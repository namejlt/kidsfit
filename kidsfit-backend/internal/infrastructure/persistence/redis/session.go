package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// SessionManager Redis会话管理器，基于Redis实现用户会话的创建、查询和删除
type SessionManager struct {
	// cache Redis缓存客户端
	cache *RedisCache
}

// NewSessionManager 创建会话管理器实例
// cache: Redis缓存客户端
func NewSessionManager(cache *RedisCache) *SessionManager {
	return &SessionManager{cache: cache}
}

// CreateSession 创建用户会话，生成唯一的sessionID并存入Redis
// ctx: 上下文，userID: 用户ID，ttl: 会话过期时间
// 返回新创建的sessionID
func (sm *SessionManager) CreateSession(ctx context.Context, userID string, ttl time.Duration) (string, error) {
	sessionID := uuid.New().String()
	key := fmt.Sprintf("session:%s", sessionID)

	if err := sm.cache.Set(ctx, key, userID, ttl); err != nil {
		return "", fmt.Errorf("创建会话失败: %w", err)
	}

	return sessionID, nil
}

// GetSession 根据sessionID获取关联的用户ID
// ctx: 上下文，sessionID: 会话ID
// 返回关联的用户ID，如果会话不存在或已过期返回空字符串
func (sm *SessionManager) GetSession(ctx context.Context, sessionID string) (string, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	userID, err := sm.cache.Get(ctx, key)
	if err != nil {
		return "", fmt.Errorf("获取会话失败: %w", err)
	}
	return userID, nil
}

// DeleteSession 删除指定会话，用于用户登出场景
// ctx: 上下文，sessionID: 待删除的会话ID
func (sm *SessionManager) DeleteSession(ctx context.Context, sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	if err := sm.cache.Delete(ctx, key); err != nil {
		return fmt.Errorf("删除会话失败: %w", err)
	}
	return nil
}
