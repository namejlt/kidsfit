package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"github.com/kidsfit/api/internal/pkg/config"
)

// RedisCache Redis缓存客户端，封装go-redis操作
type RedisCache struct {
	client *goredis.Client
}

// NewRedisCache 根据配置创建Redis缓存客户端实例
// cfg: Redis连接配置，包含地址、密码和数据库编号
func NewRedisCache(cfg *config.RedisConfig) (*RedisCache, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 验证连接是否正常
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("Redis连接失败: %w", err)
	}

	return &RedisCache{client: client}, nil
}

// Get 根据键获取字符串值
// ctx: 上下文，key: 缓存键
func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == goredis.Nil {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("获取缓存失败: %w", err)
	}
	return val, nil
}

// Set 设置字符串键值对，支持过期时间
// ctx: 上下文，key: 缓存键，value: 缓存值，ttl: 过期时间
func (r *RedisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	if err := r.client.Set(ctx, key, value, ttl).Err(); err != nil {
		return fmt.Errorf("设置缓存失败: %w", err)
	}
	return nil
}

// Delete 删除一个或多个缓存键
// ctx: 上下文，keys: 待删除的缓存键列表
func (r *RedisCache) Delete(ctx context.Context, keys ...string) error {
	if err := r.client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("删除缓存失败: %w", err)
	}
	return nil
}

// GetJSON 根据键获取JSON值并反序列化到目标对象
// ctx: 上下文，key: 缓存键，dest: 反序列化目标对象指针
func (r *RedisCache) GetJSON(ctx context.Context, key string, dest interface{}) error {
	val, err := r.Get(ctx, key)
	if err != nil {
		return err
	}
	if val == "" {
		return goredis.Nil
	}
	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return fmt.Errorf("反序列化缓存JSON失败: %w", err)
	}
	return nil
}

// SetJSON 将对象序列化为JSON并设置到缓存
// ctx: 上下文，key: 缓存键，value: 待序列化的对象，ttl: 过期时间
func (r *RedisCache) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化缓存JSON失败: %w", err)
	}
	return r.Set(ctx, key, string(data), ttl)
}

// Client 获取底层go-redis客户端，用于高级操作
func (r *RedisCache) Client() *goredis.Client {
	return r.client
}

// Close 关闭Redis连接
func (r *RedisCache) Close() error {
	return r.client.Close()
}
