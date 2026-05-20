package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"

	"github.com/kidsfit/api/internal/infrastructure/persistence/redis"
	appErrors "github.com/kidsfit/api/internal/pkg/errors"
	"github.com/kidsfit/api/internal/pkg/response"
)

// RateLimitMiddleware 基于Redis ZSet实现滑动窗口限流中间件
// 限制每个用户在指定时间窗口内最多发送limit次请求
// cache: Redis缓存客户端，limit: 时间窗口内允许的最大请求数，window: 滑动窗口时间范围
func RateLimitMiddleware(cache *redis.RedisCache, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取客户端标识，优先使用认证用户ID，否则使用客户端IP
		userID := GetUserID(c)
		if userID == "" {
			userID = c.ClientIP()
		}

		// 构建限流键名
		key := fmt.Sprintf("ratelimit:%s", userID)
		now := time.Now()
		windowStart := now.Add(-window)

		ctx := c.Request.Context()
		client := cache.Client()

		// 使用Redis管道批量执行命令，减少网络往返
		pipe := client.Pipeline()

		// 移除窗口外的旧记录
		pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart.UnixNano()))
		// 添加当前请求记录
		pipe.ZAdd(ctx, key, goredis.Z{Score: float64(now.UnixNano()), Member: userID})
		// 获取当前窗口内的请求数量
		countCmd := pipe.ZCard(ctx, key)
		// 设置键的过期时间，防止内存泄漏
		pipe.Expire(ctx, key, window)

		// 执行管道命令
		_, err := pipe.Exec(ctx)
		if err != nil {
			// Redis异常时放行请求，避免影响正常业务
			c.Next()
			return
		}

		// 检查请求次数是否超过限制
		if countCmd.Val() > int64(limit) {
			response.Error(c, appErrors.ErrTooManyRequests)
			c.Abort()
			return
		}

		c.Next()
	}
}
