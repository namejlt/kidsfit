package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware 日志中间件，记录每个请求的方法、路径、状态码和耗时
// 使用标准库log输出请求日志，便于调试和监控
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 计算请求耗时
		duration := time.Since(startTime)

		// 输出请求日志：方法、路径、状态码、耗时
		log.Printf("[HTTP] %s %s %d %v",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			duration,
		)
	}
}
