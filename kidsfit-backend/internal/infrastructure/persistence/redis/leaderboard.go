package redis

import (
	"context"
	"fmt"

	goredis "github.com/redis/go-redis/v9"
)

// MemberScore 排行榜成员及其分数
type MemberScore struct {
	// Member 成员标识
	Member string `json:"member"`
	// Score 成员分数
	Score float64 `json:"score"`
}

// Leaderboard 排行榜，基于Redis有序集合（ZSet）实现
type Leaderboard struct {
	cache *RedisCache
}

// NewLeaderboard 创建排行榜实例
// cache: Redis缓存客户端
func NewLeaderboard(cache *RedisCache) *Leaderboard {
	return &Leaderboard{cache: cache}
}

// AddScore 向排行榜添加或更新成员分数
// 使用ZADD命令，如果成员已存在则更新分数
// ctx: 上下文，key: 排行榜键名，member: 成员标识，score: 成员分数
func (lb *Leaderboard) AddScore(ctx context.Context, key string, member string, score float64) error {
	if err := lb.cache.Client().ZAdd(ctx, key, goredis.Z{
		Score:  score,
		Member: member,
	}).Err(); err != nil {
		return fmt.Errorf("添加排行榜分数失败: %w", err)
	}
	return nil
}

// GetTopN 获取排行榜前N名成员，按分数从高到低排序
// 使用ZREVRANGE命令获取降序排列的成员
// ctx: 上下文，key: 排行榜键名，n: 获取的成员数量
func (lb *Leaderboard) GetTopN(ctx context.Context, key string, n int64) ([]MemberScore, error) {
	// 使用ZREVRANGE按分数从高到低获取前N名
	results, err := lb.cache.Client().ZRevRangeWithScores(ctx, key, 0, n-1).Result()
	if err != nil {
		return nil, fmt.Errorf("获取排行榜前N名失败: %w", err)
	}

	members := make([]MemberScore, 0, len(results))
	for _, z := range results {
		member, ok := z.Member.(string)
		if !ok {
			continue
		}
		members = append(members, MemberScore{
			Member: member,
			Score:  z.Score,
		})
	}
	return members, nil
}

// GetRank 获取成员在排行榜中的排名，排名从0开始（0为第一名）
// 使用ZREVRANK命令获取降序排名
// ctx: 上下文，key: 排行榜键名，member: 成员标识
func (lb *Leaderboard) GetRank(ctx context.Context, key string, member string) (int64, error) {
	rank, err := lb.cache.Client().ZRevRank(ctx, key, member).Result()
	if err == goredis.Nil {
		// 成员不在排行榜中
		return -1, nil
	}
	if err != nil {
		return -1, fmt.Errorf("获取排行榜排名失败: %w", err)
	}
	return rank, nil
}

// IncrScore 增加成员分数，使用ZINCRBY命令原子性增加
// ctx: 上下文，key: 排行榜键名，member: 成员标识，increment: 增量值
func (lb *Leaderboard) IncrScore(ctx context.Context, key string, member string, increment float64) (float64, error) {
	newScore, err := lb.cache.Client().ZIncrBy(ctx, key, increment, member).Result()
	if err != nil {
		return 0, fmt.Errorf("增加排行榜分数失败: %w", err)
	}
	return newScore, nil
}

// RemoveMember 从排行榜中移除成员
// ctx: 上下文，key: 排行榜键名，member: 待移除的成员标识
func (lb *Leaderboard) RemoveMember(ctx context.Context, key string, member string) error {
	if err := lb.cache.Client().ZRem(ctx, key, member).Err(); err != nil {
		return fmt.Errorf("移除排行榜成员失败: %w", err)
	}
	return nil
}

// GetScore 获取成员的当前分数
// ctx: 上下文，key: 排行榜键名，member: 成员标识
func (lb *Leaderboard) GetScore(ctx context.Context, key string, member string) (float64, error) {
	score, err := lb.cache.Client().ZScore(ctx, key, member).Result()
	if err == goredis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("获取排行榜分数失败: %w", err)
	}
	return score, nil
}
