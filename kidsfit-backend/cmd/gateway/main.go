package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/kidsfit/api/internal/pkg/config"
)

// routeRule 路由转发规则，定义URL前缀到后端服务的映射
type routeRule struct {
	prefix    string // URL路径前缀
	targetURL string // 后端服务地址
}

// routeRules 路由转发规则列表，按优先级顺序匹配
var routeRules = []routeRule{
	// 用户服务路由
	{prefix: "/api/v1/auth", targetURL: "http://localhost:8001"},
	{prefix: "/api/v1/users", targetURL: "http://localhost:8001"},
	{prefix: "/api/v1/children", targetURL: "http://localhost:8001"},
	{prefix: "/api/v1/settings", targetURL: "http://localhost:8001"},
	// 训练服务路由
	{prefix: "/api/v1/exercises", targetURL: "http://localhost:8002"},
	{prefix: "/api/v1/plans", targetURL: "http://localhost:8002"},
	{prefix: "/api/v1/assessments", targetURL: "http://localhost:8002"},
	{prefix: "/api/v1/stats", targetURL: "http://localhost:8002"},
	// 视力服务路由
	{prefix: "/api/v1/vision-records", targetURL: "http://localhost:8003"},
	{prefix: "/api/v1/vision", targetURL: "http://localhost:8003"},
	{prefix: "/api/v1/outdoor", targetURL: "http://localhost:8003"},
	{prefix: "/api/v1/reminders", targetURL: "http://localhost:8003"},
	// 激励服务路由
	{prefix: "/api/v1/badges", targetURL: "http://localhost:8004"},
	{prefix: "/api/v1/points", targetURL: "http://localhost:8004"},
	{prefix: "/api/v1/challenges", targetURL: "http://localhost:8004"},
	{prefix: "/api/v1/leaderboard", targetURL: "http://localhost:8004"},
}

// rateLimitRecord 限流记录，记录每个客户端的请求时间戳
type rateLimitRecord struct {
	timestamps []int64 // 请求时间戳列表
}

// gateway 网关实例，包含JWT密钥和限流记录
type gateway struct {
	jwtSecret    string                       // JWT签名密钥
	rateLimits   map[string]*rateLimitRecord  // 客户端限流记录
	rateLimitMax int                          // 每分钟最大请求数
}

// newGateway 创建网关实例
// jwtSecret: JWT签名密钥，rateLimitMax: 每分钟最大请求数
func newGateway(jwtSecret string, rateLimitMax int) *gateway {
	return &gateway{
		jwtSecret:    jwtSecret,
		rateLimits:   make(map[string]*rateLimitRecord),
		rateLimitMax: rateLimitMax,
	}
}

// authMiddleware JWT认证中间件，验证请求中的Bearer Token
// 对于/auth/register和/auth/login路径跳过认证
func (g *gateway) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// 认证相关公开接口跳过JWT验证
		if strings.HasSuffix(path, "/auth/register") ||
			strings.HasSuffix(path, "/auth/login") ||
			strings.HasSuffix(path, "/auth/refresh") {
			c.Next()
			return
		}

		// 从Header中提取Authorization字段
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未提供认证令牌"})
			c.Abort()
			return
		}

		// 解析Bearer Token格式
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "认证令牌格式错误"})
			c.Abort()
			return
		}

		// 解析并验证JWT令牌
		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("不支持的签名方法: %v", token.Header["alg"])
			}
			return []byte(g.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "认证令牌无效或已过期"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// rateLimitMiddleware 统一限流中间件，基于滑动窗口算法限制每个客户端的请求频率
// 每个客户端（按IP或用户ID标识）在60秒内最多允许rateLimitMax次请求
func (g *gateway) rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取客户端标识，优先使用认证用户ID，否则使用客户端IP
		clientID := c.ClientIP()
		if authHeader := c.GetHeader("Authorization"); authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				// 尝试从token中提取用户ID
				if token, _, err := new(jwt.Parser).ParseUnverified(parts[1], &jwt.MapClaims{}); err == nil {
					if claims, ok := token.Claims.(*jwt.MapClaims); ok {
						if uid, ok := (*claims)["user_id"].(string); ok && uid != "" {
							clientID = uid
						}
					}
				}
			}
		}

		now := time.Now().UnixNano()
		windowStart := now - int64(time.Minute)

		// 获取或创建限流记录
		record, exists := g.rateLimits[clientID]
		if !exists {
			record = &rateLimitRecord{timestamps: []int64{}}
			g.rateLimits[clientID] = record
		}

		// 移除窗口外的旧记录
		validIdx := 0
		for i, ts := range record.timestamps {
			if ts >= windowStart {
				validIdx = i
				break
			}
			validIdx = i + 1
		}
		record.timestamps = record.timestamps[validIdx:]

		// 检查是否超过限流阈值
		if len(record.timestamps) >= g.rateLimitMax {
			c.JSON(http.StatusTooManyRequests, gin.H{"code": 429, "message": "请求过于频繁，请稍后再试"})
			c.Abort()
			return
		}

		// 记录当前请求时间戳
		record.timestamps = append(record.timestamps, now)

		c.Next()
	}
}

// matchRoute 根据请求路径匹配路由规则，返回对应的后端服务地址
// path: 请求路径，返回匹配的目标URL，未匹配则返回空字符串
func matchRoute(path string) string {
	for _, rule := range routeRules {
		if strings.HasPrefix(path, rule.prefix) {
			return rule.targetURL
		}
	}
	return ""
}

// createReverseProxy 创建反向代理处理器
// targetURL: 后端服务地址字符串
func createReverseProxy(targetURL string) (*httputil.ReverseProxy, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("解析目标地址失败: %w", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	// 自定义请求修改器，设置X-Forwarded-For头部
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Header.Set("X-Forwarded-Proto", "http")
	}

	return proxy, nil
}

// proxyHandler 反向代理处理器，将请求转发到匹配的后端服务
// 未匹配路由时返回404
func (g *gateway) proxyHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// 匹配路由规则
		targetURL := matchRoute(path)
		if targetURL == "" {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "路由未找到"})
			c.Abort()
			return
		}

		// 创建反向代理
		proxy, err := createReverseProxy(targetURL)
		if err != nil {
			log.Printf("创建反向代理失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "网关内部错误"})
			c.Abort()
			return
		}

		// 执行反向代理
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// setupGatewayRouter 初始化网关路由引擎，注册中间件和代理处理器
// gw: 网关实例
func setupGatewayRouter(gw *gateway) *gin.Engine {
	r := gin.New()

	// 注册全局中间件
	r.Use(gin.Logger())                // 请求日志
	r.Use(gin.Recovery())              // 异常恢复
	r.Use(gw.rateLimitMiddleware())    // 统一限流
	r.Use(gw.authMiddleware())         // 统一JWT认证

	// 健康检查接口
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "gateway"})
	})

	// 代理所有API请求到后端服务
	r.Any("/api/v1/*path", gw.proxyHandler())

	return r
}

// main API网关服务入口函数，负责初始化网关依赖并启动HTTP服务
func main() {
	// 1. 加载配置文件
	cfg, err := config.Load("configs/dev.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 2. 创建网关实例，每分钟限流120次
	gw := newGateway(cfg.JWT.Secret, 120)

	// 3. 初始化路由
	router := setupGatewayRouter(gw)

	// 4. 启动HTTP服务，监听8080端口
	addr := fmt.Sprintf(":%d", 8080)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// 在goroutine中启动服务
	go func() {
		log.Printf("API网关启动，监听地址: %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP服务启动失败: %v", err)
		}
	}()

	// 5. 优雅关闭：等待系统信号后安全退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭API网关...")

	// 设置5秒超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅关闭HTTP服务
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("API网关关闭异常: %v", err)
	}

	log.Println("API网关已关闭")
}
