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
	"github.com/kidsfit/api/internal/application/training"
	"github.com/kidsfit/api/internal/infrastructure/persistence/postgresql"
	"github.com/kidsfit/api/internal/infrastructure/persistence/redis"
	"github.com/kidsfit/api/internal/pkg/config"
)

// main 训练服务入口函数，负责初始化训练相关依赖并启动HTTP服务
func main() {
	// 1. 加载配置文件
	cfg, err := config.Load("configs/dev.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 覆盖数据库名称为训练服务专用库
	cfg.Database.DBName = "kidsfit_training"

	// 2. 初始化PostgreSQL数据库连接
	db, err := postgresql.NewDB(&cfg.Database)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer db.Close()

	// 3. 初始化Redis缓存连接
	redisCache, err := redis.NewRedisCache(&cfg.Redis)
	if err != nil {
		log.Fatalf("初始化Redis失败: %v", err)
	}
	defer redisCache.Close()

	// 4. 初始化训练相关仓储层
	exerciseRepo := postgresql.NewPostgresExerciseRecordRepo(db)
	planRepo := postgresql.NewPostgresTrainingPlanRepo(db)
	assessmentRepo := postgresql.NewPostgresFitnessAssessmentRepo(db)

	// 5. 初始化训练应用服务层
	trainingAppService := training.NewTrainingAppService(exerciseRepo, planRepo, assessmentRepo)

	// 6. 初始化训练HTTP处理器
	trainingHandler := handler.NewTrainingHandler(trainingAppService)

	// 7. 初始化路由，训练服务仅注册训练相关路由
	router := httpAPI.SetupTrainingRouter(trainingHandler, cfg.JWT.Secret, redisCache)

	// 8. 启动HTTP服务，监听8002端口
	addr := fmt.Sprintf(":%d", 8002)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// 在goroutine中启动服务
	go func() {
		log.Printf("训练服务启动，监听地址: %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP服务启动失败: %v", err)
		}
	}()

	// 9. 优雅关闭：等待系统信号后安全退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭训练服务...")

	// 设置5秒超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅关闭HTTP服务
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("训练服务关闭异常: %v", err)
	}

	log.Println("训练服务已关闭")
}
