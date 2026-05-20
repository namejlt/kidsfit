package redis

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

// DistributedLock 分布式锁，基于Redis SET NX EX命令实现
type DistributedLock struct {
	cache *RedisCache
}

// NewDistributedLock 创建分布式锁实例
// cache: Redis缓存客户端
func NewDistributedLock(cache *RedisCache) *DistributedLock {
	return &DistributedLock{cache: cache}
}

// Acquire 尝试获取分布式锁
// 使用SET key value NX EX命令实现互斥，key为锁标识，ttl为锁的过期时间
// ctx: 上下文，key: 锁的键名，ttl: 锁的过期时间
// 返回true表示获取锁成功，false表示锁已被占用
func (l *DistributedLock) Acquire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	// 使用SET NX EX命令，仅当key不存在时设置值
	ok, err := l.cache.Client().SetNX(ctx, key, "locked", ttl).Result()
	if err != nil {
		return false, fmt.Errorf("获取分布式锁失败: %w", err)
	}
	return ok, nil
}

// Release 释放分布式锁
// 直接删除key来释放锁，适用于同一进程内的锁释放场景
// ctx: 上下文，key: 锁的键名
func (l *DistributedLock) Release(ctx context.Context, key string) error {
	if err := l.cache.Delete(ctx, key); err != nil {
		return fmt.Errorf("释放分布式锁失败: %w", err)
	}
	return nil
}

// AcquireWithValue 尝试获取分布式锁并设置自定义值
// 可用于标识锁的持有者，便于后续安全释放
// ctx: 上下文，key: 锁的键名，value: 锁的值（如UUID），ttl: 锁的过期时间
func (l *DistributedLock) AcquireWithValue(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
	ok, err := l.cache.Client().SetNX(ctx, key, value, ttl).Result()
	if err != nil {
		return false, fmt.Errorf("获取分布式锁失败: %w", err)
	}
	return ok, nil
}

// ReleaseWithValue 安全释放分布式锁，仅当锁的值匹配时才释放
// 使用Lua脚本保证原子性，防止误释放其他持有者的锁
// ctx: 上下文，key: 锁的键名，value: 锁的预期值
func (l *DistributedLock) ReleaseWithValue(ctx context.Context, key string, value string) error {
	// Lua脚本：仅当key的值等于预期值时才删除
	script := goredis.NewScript(`
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`)
	_, err := script.Run(ctx, l.cache.Client(), []string{key}, value).Result()
	if err != nil {
		return fmt.Errorf("安全释放分布式锁失败: %w", err)
	}
	return nil
}
