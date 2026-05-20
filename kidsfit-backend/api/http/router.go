package http

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/kidsfit/api/api/http/handler"
	"github.com/kidsfit/api/api/http/middleware"
	"github.com/kidsfit/api/internal/infrastructure/persistence/redis"
)

// SetupRouter 初始化并配置Gin路由引擎
// 注册所有中间件和路由，按API路径分组组织
// userHandler: 用户处理器，trainingHandler: 训练处理器
// visionHandler: 视力处理器，rewardHandler: 激励处理器
// jwtSecret: JWT签名密钥，redisCache: Redis缓存客户端
func SetupRouter(
	userHandler *handler.UserHandler,
	trainingHandler *handler.TrainingHandler,
	visionHandler *handler.VisionHandler,
	rewardHandler *handler.RewardHandler,
	jwtSecret string,
	redisCache *redis.RedisCache,
) *gin.Engine {
	r := gin.New()

	// 注册全局中间件
	r.Use(middleware.LoggerMiddleware())   // 请求日志
	r.Use(middleware.CORSMiddleware())     // 跨域支持
	r.Use(gin.Recovery())                 // 异常恢复

	// 健康检查接口
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, map[string]string{"status": "ok"})
	})

	// API v1 路由组
	v1 := r.Group("/api/v1")

	// 认证相关路由（无需认证）
	auth := v1.Group("/auth")
	{
		auth.POST("/register", userHandler.Register)      // 注册家长账号
		auth.POST("/login", userHandler.Login)            // 用户登录
		auth.POST("/refresh", userHandler.RefreshToken)   // 刷新令牌
	}

	// 需要认证的路由组
	authenticated := v1.Group("")
	authenticated.Use(middleware.AuthMiddleware(jwtSecret))
	{
		// 认证相关（需认证）
		authenticated.POST("/auth/logout", userHandler.Logout) // 用户登出

		// 用户相关路由
		users := authenticated.Group("/users")
		{
			users.GET("/me", userHandler.GetCurrentUser)   // 获取当前用户信息
			users.PUT("/me", userHandler.UpdateUser)       // 更新当前用户信息
			users.POST("/avatar", userHandler.UploadAvatar) // 上传用户头像
		}

		// 儿童相关路由（仅家长）
		children := authenticated.Group("/children")
		children.Use(middleware.RequireParent())
		{
			children.POST("", userHandler.AddChild)    // 添加儿童账号
			children.GET("", userHandler.GetChildren)  // 获取儿童列表
		}

		// 家长设置路由（仅家长）
		settings := authenticated.Group("/settings")
		settings.Use(middleware.RequireParent())
		{
			settings.GET("", userHandler.GetParentSettings)      // 获取家长设置
			settings.PUT("", userHandler.UpdateParentSettings)   // 更新家长设置
		}

		// 运动记录路由（限流：每分钟60次）
		exercises := authenticated.Group("/exercises")
		exercises.Use(middleware.RateLimitMiddleware(redisCache, 60, time.Minute))
		{
			exercises.POST("", trainingHandler.CreateExerciseRecord)       // 创建运动记录
			exercises.GET("", trainingHandler.GetExerciseRecords)          // 查询运动记录
			exercises.GET("/personal-best", trainingHandler.GetPersonalBest) // 获取个人最佳
		}

		// 训练计划路由
		plans := authenticated.Group("/plans")
		{
			plans.GET("/today", trainingHandler.GetTodayPlan)              // 获取今日计划
			plans.POST("/:id/complete", trainingHandler.CompletePlan)      // 完成计划
		}

		// 体能评估路由
		assessments := authenticated.Group("/assessments")
		{
			assessments.POST("", trainingHandler.CreateFitnessAssessment)  // 创建体能评估
			assessments.GET("/latest", trainingHandler.GetLatestAssessment) // 获取最新评估
		}

		// 运动统计路由
		stats := authenticated.Group("/stats")
		{
			stats.GET("/weekly", trainingHandler.GetWeeklyStats)   // 获取周统计
			stats.GET("/monthly", trainingHandler.GetMonthlyStats) // 获取月统计
		}

		// 视力记录路由
		visionRecords := authenticated.Group("/vision-records")
		{
			visionRecords.POST("", visionHandler.CreateVisionRecord)      // 创建视力记录
			visionRecords.POST("/ocr", visionHandler.OCRVisionRecord)     // OCR识别验光单
			visionRecords.GET("", visionHandler.GetVisionRecords)         // 查询视力记录
		}

		// 视力趋势路由
		visionGroup := authenticated.Group("/vision")
		{
			visionGroup.GET("/trend", visionHandler.GetVisionTrend) // 获取视力趋势
		}

		// 户外活动路由
		outdoor := authenticated.Group("/outdoor")
		{
			outdoor.GET("/today", visionHandler.GetTodayOutdoor) // 获取今日户外活动
			outdoor.POST("/sync", visionHandler.SyncOutdoorData) // 同步户外活动数据
		}

		// 护眼提醒路由
		reminders := authenticated.Group("/reminders")
		{
			reminders.GET("", visionHandler.GetReminders)           // 查询护眼提醒
			reminders.POST("/:id/ack", visionHandler.AckReminder)   // 确认提醒
		}

		// 徽章路由
		badges := authenticated.Group("/badges")
		{
			badges.GET("", rewardHandler.GetBadges)       // 查询徽章列表
			badges.GET("/my", rewardHandler.GetMyBadges)  // 获取已获得徽章
		}

		// 积分路由
		points := authenticated.Group("/points")
		{
			points.GET("", rewardHandler.GetPoints)           // 查询积分记录
			points.GET("/balance", rewardHandler.GetPointsBalance) // 获取积分余额
		}

		// 挑战路由
		challenges := authenticated.Group("/challenges")
		{
			challenges.POST("", rewardHandler.CreateChallenge)              // 创建挑战
			challenges.GET("", rewardHandler.GetChallenges)                // 获取挑战列表
			challenges.POST("/:id/accept", rewardHandler.AcceptChallenge)  // 接受挑战
			challenges.POST("/:id/submit", rewardHandler.SubmitChallengeScore) // 提交挑战成绩
		}

		// 排行榜路由
		leaderboard := authenticated.Group("/leaderboard")
		{
			leaderboard.GET("/family", rewardHandler.GetFamilyLeaderboard) // 家庭排行榜
			leaderboard.GET("/global", rewardHandler.GetGlobalLeaderboard) // 全局排行榜
		}
	}

	return r
}

// SetupTrainingRouter 初始化训练服务的Gin路由引擎
// 仅注册训练相关的路由，供训练微服务独立部署使用
// trainingHandler: 训练处理器，jwtSecret: JWT签名密钥，redisCache: Redis缓存客户端
func SetupTrainingRouter(
	trainingHandler *handler.TrainingHandler,
	jwtSecret string,
	redisCache *redis.RedisCache,
) *gin.Engine {
	r := gin.New()

	// 注册全局中间件
	r.Use(middleware.LoggerMiddleware())   // 请求日志
	r.Use(middleware.CORSMiddleware())     // 跨域支持
	r.Use(gin.Recovery())                 // 异常恢复

	// 健康检查接口
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, map[string]string{"status": "ok"})
	})

	// API v1 路由组
	v1 := r.Group("/api/v1")

	// 需要认证的路由组
	authenticated := v1.Group("")
	authenticated.Use(middleware.AuthMiddleware(jwtSecret))
	{
		// 运动记录路由（限流：每分钟60次）
		exercises := authenticated.Group("/exercises")
		exercises.Use(middleware.RateLimitMiddleware(redisCache, 60, time.Minute))
		{
			exercises.POST("", trainingHandler.CreateExerciseRecord)       // 创建运动记录
			exercises.GET("", trainingHandler.GetExerciseRecords)          // 查询运动记录
			exercises.GET("/personal-best", trainingHandler.GetPersonalBest) // 获取个人最佳
		}

		// 训练计划路由
		plans := authenticated.Group("/plans")
		{
			plans.GET("/today", trainingHandler.GetTodayPlan)              // 获取今日计划
			plans.POST("/:id/complete", trainingHandler.CompletePlan)      // 完成计划
		}

		// 体能评估路由
		assessments := authenticated.Group("/assessments")
		{
			assessments.POST("", trainingHandler.CreateFitnessAssessment)  // 创建体能评估
			assessments.GET("/latest", trainingHandler.GetLatestAssessment) // 获取最新评估
		}

		// 运动统计路由
		stats := authenticated.Group("/stats")
		{
			stats.GET("/weekly", trainingHandler.GetWeeklyStats)   // 获取周统计
			stats.GET("/monthly", trainingHandler.GetMonthlyStats) // 获取月统计
		}
	}

	return r
}

// SetupVisionRouter 初始化视力服务的Gin路由引擎
// 仅注册视力相关的路由，供视力微服务独立部署使用
// visionHandler: 视力处理器，jwtSecret: JWT签名密钥，redisCache: Redis缓存客户端
func SetupVisionRouter(
	visionHandler *handler.VisionHandler,
	jwtSecret string,
	redisCache *redis.RedisCache,
) *gin.Engine {
	r := gin.New()

	// 注册全局中间件
	r.Use(middleware.LoggerMiddleware())   // 请求日志
	r.Use(middleware.CORSMiddleware())     // 跨域支持
	r.Use(gin.Recovery())                 // 异常恢复

	// 健康检查接口
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, map[string]string{"status": "ok"})
	})

	// API v1 路由组
	v1 := r.Group("/api/v1")

	// 需要认证的路由组
	authenticated := v1.Group("")
	authenticated.Use(middleware.AuthMiddleware(jwtSecret))
	{
		// 视力记录路由
		visionRecords := authenticated.Group("/vision-records")
		{
			visionRecords.POST("", visionHandler.CreateVisionRecord)      // 创建视力记录
			visionRecords.POST("/ocr", visionHandler.OCRVisionRecord)     // OCR识别验光单
			visionRecords.GET("", visionHandler.GetVisionRecords)         // 查询视力记录
		}

		// 视力趋势路由
		visionGroup := authenticated.Group("/vision")
		{
			visionGroup.GET("/trend", visionHandler.GetVisionTrend) // 获取视力趋势
		}

		// 户外活动路由
		outdoor := authenticated.Group("/outdoor")
		{
			outdoor.GET("/today", visionHandler.GetTodayOutdoor) // 获取今日户外活动
			outdoor.POST("/sync", visionHandler.SyncOutdoorData) // 同步户外活动数据
		}

		// 护眼提醒路由
		reminders := authenticated.Group("/reminders")
		{
			reminders.GET("", visionHandler.GetReminders)           // 查询护眼提醒
			reminders.POST("/:id/ack", visionHandler.AckReminder)   // 确认提醒
		}
	}

	return r
}

// SetupRewardRouter 初始化激励服务的Gin路由引擎
// 仅注册激励相关的路由，供激励微服务独立部署使用
// rewardHandler: 激励处理器，jwtSecret: JWT签名密钥，redisCache: Redis缓存客户端
func SetupRewardRouter(
	rewardHandler *handler.RewardHandler,
	jwtSecret string,
	redisCache *redis.RedisCache,
) *gin.Engine {
	r := gin.New()

	// 注册全局中间件
	r.Use(middleware.LoggerMiddleware())   // 请求日志
	r.Use(middleware.CORSMiddleware())     // 跨域支持
	r.Use(gin.Recovery())                 // 异常恢复

	// 健康检查接口
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, map[string]string{"status": "ok"})
	})

	// API v1 路由组
	v1 := r.Group("/api/v1")

	// 需要认证的路由组
	authenticated := v1.Group("")
	authenticated.Use(middleware.AuthMiddleware(jwtSecret))
	{
		// 徽章路由
		badges := authenticated.Group("/badges")
		{
			badges.GET("", rewardHandler.GetBadges)       // 查询徽章列表
			badges.GET("/my", rewardHandler.GetMyBadges)  // 获取已获得徽章
		}

		// 积分路由
		points := authenticated.Group("/points")
		{
			points.GET("", rewardHandler.GetPoints)           // 查询积分记录
			points.GET("/balance", rewardHandler.GetPointsBalance) // 获取积分余额
		}

		// 挑战路由
		challenges := authenticated.Group("/challenges")
		{
			challenges.POST("", rewardHandler.CreateChallenge)              // 创建挑战
			challenges.GET("", rewardHandler.GetChallenges)                // 获取挑战列表
			challenges.POST("/:id/accept", rewardHandler.AcceptChallenge)  // 接受挑战
			challenges.POST("/:id/submit", rewardHandler.SubmitChallengeScore) // 提交挑战成绩
		}

		// 排行榜路由
		leaderboard := authenticated.Group("/leaderboard")
		{
			leaderboard.GET("/family", rewardHandler.GetFamilyLeaderboard) // 家庭排行榜
			leaderboard.GET("/global", rewardHandler.GetGlobalLeaderboard) // 全局排行榜
		}
	}

	return r
}
