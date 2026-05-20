package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	events "github.com/kidsfit/api/internal/infrastructure/messaging"
)

// Publisher RabbitMQ事件发布器，封装AMQP连接和消息发布操作
type Publisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewPublisher 创建RabbitMQ发布器实例
// 支持连接重试机制，最多重试3次，采用指数退避策略
// url: RabbitMQ连接地址（如 amqp://guest:guest@localhost:5672/）
func NewPublisher(url string) (*Publisher, error) {
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

	return &Publisher{
		conn:    conn,
		channel: ch,
	}, nil
}

// Publish 发布事件到指定的RabbitMQ交换机
// event: 待发布的事件，exchange: 交换机名称，routingKey: 路由键
func (p *Publisher) Publish(event *events.Event, exchange, routingKey string) error {
	// 序列化事件为JSON
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("序列化事件失败: %w", err)
	}

	// 发布消息
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = p.channel.PublishWithContext(
		ctx,
		exchange,   // 交换机
		routingKey, // 路由键
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent, // 持久化消息
			MessageId:    event.ID,
			Timestamp:    event.Timestamp,
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("发布事件失败: %w", err)
	}

	return nil
}

// Close 关闭发布器，释放RabbitMQ连接和通道资源
func (p *Publisher) Close() error {
	if p.channel != nil {
		if err := p.channel.Close(); err != nil {
			return fmt.Errorf("关闭RabbitMQ通道失败: %w", err)
		}
	}
	if p.conn != nil {
		if err := p.conn.Close(); err != nil {
			return fmt.Errorf("关闭RabbitMQ连接失败: %w", err)
		}
	}
	return nil
}
