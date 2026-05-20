package rabbitmq

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	events "github.com/kidsfit/api/internal/infrastructure/messaging"
)

// Subscriber RabbitMQ事件订阅器，封装消息消费和处理逻辑
type Subscriber struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewSubscriber 创建RabbitMQ订阅器实例
// 支持连接重试机制，最多重试3次，采用指数退避策略
// url: RabbitMQ连接地址（如 amqp://guest:guest@localhost:5672/）
func NewSubscriber(url string) (*Subscriber, error) {
	var conn *amqp.Connection
	var err error

	// 指数退避重试，最多3次
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			break
		}
		// 计算退避时间：2^attempt * 基础时间
		backoff := time.Duration(math.Pow(2, float64(attempt))) * time.Second
		time.Sleep(backoff)
	}
	if err != nil {
		return nil, fmt.Errorf("连接RabbitMQ失败（重试%d次）: %w", maxRetries, err)
	}

	// 创建通道
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("创建RabbitMQ通道失败: %w", err)
	}

	// 设置QoS，每次只预取一条消息，确保手动ACK机制生效
	if err := ch.Qos(1, 0, false); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("设置QoS失败: %w", err)
	}

	return &Subscriber{
		conn:    conn,
		channel: ch,
	}, nil
}

// Subscribe 订阅指定队列的消息，绑定交换机和路由键
// queue: 队列名称，exchange: 交换机名称，routingKey: 路由键
// handler: 消息处理回调函数，返回error表示处理失败
// 消息处理采用手动ACK机制，处理失败重试3次后进入死信队列
func (s *Subscriber) Subscribe(queue, exchange, routingKey string, handler func(events.Event) error) error {
	// 声明死信交换机
	dlxExchange := queue + ".dlx"
	if err := s.channel.ExchangeDeclare(
		dlxExchange, // 名称
		"direct",    // 类型
		true,        // 持久化
		false,       // 自动删除
		false,       // 内部
		false,       // 等待
		nil,         // 参数
	); err != nil {
		return fmt.Errorf("声明死信交换机失败: %w", err)
	}

	// 声明死信队列
	dlqName := queue + ".dlq"
	if _, err := s.channel.QueueDeclare(
		dlqName, // 名称
		true,    // 持久化
		false,   // 自动删除
		false,   // 独占
		false,   // 等待
		nil,     // 参数
	); err != nil {
		return fmt.Errorf("声明死信队列失败: %w", err)
	}

	// 绑定死信队列到死信交换机
	if err := s.channel.QueueBind(
		dlqName,     // 队列名
		routingKey,  // 路由键
		dlxExchange, // 交换机
		false,       // 等待
		nil,         // 参数
	); err != nil {
		return fmt.Errorf("绑定死信队列失败: %w", err)
	}

	// 声明主队列，配置死信路由
	if _, err := s.channel.QueueDeclare(
		queue,  // 名称
		true,   // 持久化
		false,  // 自动删除
		false,  // 独占
		false,  // 等待
		amqp.Table{
			"x-dead-letter-exchange":    dlxExchange, // 死信交换机
			"x-dead-letter-routing-key": routingKey,  // 死信路由键
		},
	); err != nil {
		return fmt.Errorf("声明队列失败: %w", err)
	}

	// 绑定主队列到交换机
	if err := s.channel.QueueBind(
		queue,      // 队列名
		routingKey, // 路由键
		exchange,   // 交换机
		false,      // 等待
		nil,        // 参数
	); err != nil {
		return fmt.Errorf("绑定队列失败: %w", err)
	}

	// 开始消费消息，手动ACK
	deliveries, err := s.channel.Consume(
		queue, // 队列名
		"",    // 消费者标签
		false, // 手动ACK
		false, // 独占
		false, // 不等待
		false, // 参数
		nil,
	)
	if err != nil {
		return fmt.Errorf("开始消费消息失败: %w", err)
	}

	// 启动消息处理协程
	go func() {
		for delivery := range deliveries {
			var event events.Event
			if err := json.Unmarshal(delivery.Body, &event); err != nil {
				// 反序列化失败，直接拒绝不重试
				_ = delivery.Reject(false)
				continue
			}

			// 处理消息，最多重试3次
			maxRetries := 3
			retryCount := 0
			if xDeath, ok := delivery.Headers["x-death"]; ok {
				if deaths, ok := xDeath.([]interface{}); ok && len(deaths) > 0 {
					if death, ok := deaths[0].(amqp.Table); ok {
						if count, ok := death["count"].(int64); ok {
							retryCount = int(count)
						}
					}
				}
			}

			if err := handler(event); err != nil {
				if retryCount >= maxRetries {
					// 超过最大重试次数，拒绝消息（进入死信队列）
					_ = delivery.Reject(false)
				} else {
					// 重试：拒绝并重新入队
					_ = delivery.Reject(true)
				}
				continue
			}

			// 处理成功，确认消息
			_ = delivery.Ack(false)
		}
	}()

	return nil
}

// Close 关闭订阅器，释放RabbitMQ连接和通道资源
func (s *Subscriber) Close() error {
	if s.channel != nil {
		if err := s.channel.Close(); err != nil {
			return fmt.Errorf("关闭RabbitMQ通道失败: %w", err)
		}
	}
	if s.conn != nil {
		if err := s.conn.Close(); err != nil {
			return fmt.Errorf("关闭RabbitMQ连接失败: %w", err)
		}
	}
	return nil
}
