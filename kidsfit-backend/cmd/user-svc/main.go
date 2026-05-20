package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpAPI "github.com/kidsfit/api/api/http"
	"github.com/kidsfit/api/api/http/handler"
	"github.com/kidsfit/api/internal/application/reward"
	"github.com/kidsfit/api/internal/application/training"
	"github.com/kidsfit/api/internal/application/user"
	"github.com/kidsfit/api/internal/application/vision"
	"github.com/kidsfit/api/internal/infrastructure/persistence/postgresql"
	"github.com/kidsfit/api/internal/infrastructure/persistence/redis"
	"github.com/kidsfit/api/internal/pkg/config"
)

// main 服务入口函数，负责初始化所有依赖并启动HTTP服务
func main() {
	// 1. 加载配置
	cfg, err := config.Load("configs/dev.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 2. 初始化数据库连接
	db, err := postgresql.NewDB(&cfg.Database)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer db.Close()

	// 3. 初始化Redis缓存
	redisCache, err := redis.NewRedisCache(&cfg.Redis)
	if err != nil {
		log.Fatalf("初始化Redis失败: %v", err)
	}
	defer redisCache.Close()

	// 4. 初始化仓储层
	userRepo := postgresql.NewPostgresUserRepo(db)
	familyRepo := postgresql.NewPostgresFamilyRepo(db)
	settingsRepo := postgresql.NewPostgresParentSettingsRepo(db)
	exerciseRepo := postgresql.NewPostgresExerciseRecordRepo(db)
	planRepo := postgresql.NewPostgresTrainingPlanRepo(db)
	assessmentRepo := postgresql.NewPostgresFitnessAssessmentRepo(db)
	visionRecordRepo := postgresql.NewPostgresVisionRecordRepo(db)
	outdoorRepo := postgresql.NewPostgresOutdoorActivityRepo(db)
	reminderRepo := postgresql.NewPostgresEyeReminderRepo(db)
	badgeRepo := postgresql.NewPostgresBadgeRepo(db)
	userBadgeRepo := postgresql.NewPostgresUserBadgeRepo(db)
	pointRepo := postgresql.NewPostgresPointRecordRepo(db)
	challengeRepo := postgresql.NewPostgresChallengeRepo(db)
	leaderboard := redis.NewLeaderboard(redisCache)

	// 5. 初始化应用服务层
	userAppService := user.NewUserAppService(userRepo, familyRepo, settingsRepo, redisCache, &cfg.JWT)
	trainingAppService := training.NewTrainingAppService(exerciseRepo, planRepo, assessmentRepo)
	visionAppService := vision.NewVisionAppService(visionRecordRepo, outdoorRepo, reminderRepo)
	rewardAppService := reward.NewRewardAppService(badgeRepo, userBadgeRepo, pointRepo, challengeRepo, userRepo, familyRepo, leaderboard)

	// 6. 初始化HTTP处理器
	userHandler := handler.NewUserHandler(userAppService)
	trainingHandler := handler.NewTrainingHandler(trainingAppService)
	visionHandler := handler.NewVisionHandler(visionAppService)
	rewardHandler := handler.NewRewardHandler(rewardAppService)

	// 7. 初始化路由
	router := httpAPI.SetupRouter(userHandler, trainingHandler, visionHandler, rewardHandler, cfg.JWT.Secret, redisCache)

	// 8. 启动HTTP服务
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// 在goroutine中启动服务
	go func() {
		log.Printf("用户服务启动，监听地址: %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP服务启动失败: %v", err)
		}
	}()

	// 9. 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务...")

	// 设置5秒超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅关闭HTTP服务
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("服务关闭异常: %v", err)
	}

	log.Println("服务已关闭")
}
