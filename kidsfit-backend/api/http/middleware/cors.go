package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORSMiddleware 跨域资源共享中间件，允许所有来源访问（开发环境使用）
// 在生产环境中应限制允许的来源域名
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许所有来源访问
		c.Header("Access-Control-Allow-Origin", "*")
		// 允许的HTTP方法
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		// 允许的请求头
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		// 允许浏览器缓存预检请求结果的时间（秒）
		c.Header("Access-Control-Max-Age", "86400")
		// 允许客户端访问的响应头
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Disposition")

		// 处理预检请求（OPTIONS），直接返回204
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
