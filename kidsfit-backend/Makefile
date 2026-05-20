# KidsFit 后端构建脚本

# 项目模块名
MODULE := github.com/kidsfit/api

# 构建输出目录
BUILD_DIR := bin

# 服务列表（包含所有微服务和网关）
SERVICES := user-svc training-svc vision-svc reward-svc gateway

# Go 编译参数
LDFLAGS := -s -w

# 编译所有服务
.PHONY: build
build:
	@echo "编译所有服务..."
	@mkdir -p $(BUILD_DIR)
	@for svc in $(SERVICES); do \
		echo "编译 $${svc}..."; \
		go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$${svc} $(MODULE)/cmd/$${svc}; \
	done
	@echo "编译完成"

# 运行测试
.PHONY: test
test:
	@echo "运行测试..."
	go test -v -race ./...

# 运行代码检查
.PHONY: lint
lint:
	@echo "运行golangci-lint..."
	golangci-lint run ./...

# 执行用户服务数据库迁移
.PHONY: migrate-up
migrate-up:
	@echo "执行用户服务数据库迁移..."
	@migrate -path scripts/migrate -database "postgres://kidsfit:kidsfit_dev@localhost:5432/kidsfit_users?sslmode=disable" up

# 执行训练服务数据库迁移
.PHONY: migrate-up-training
migrate-up-training:
	@echo "执行训练服务数据库迁移..."
	@migrate -path scripts/migrate -database "postgres://kidsfit:kidsfit_dev@localhost:5432/kidsfit_training?sslmode=disable" up

# 执行视力服务数据库迁移
.PHONY: migrate-up-vision
migrate-up-vision:
	@echo "执行视力服务数据库迁移..."
	@migrate -path scripts/migrate -database "postgres://kidsfit:kidsfit_dev@localhost:5432/kidsfit_vision?sslmode=disable" up

# 执行激励服务数据库迁移
.PHONY: migrate-up-reward
migrate-up-reward:
	@echo "执行激励服务数据库迁移..."
	@migrate -path scripts/migrate -database "postgres://kidsfit:kidsfit_dev@localhost:5432/kidsfit_rewards?sslmode=disable" up

# 整理Go模块依赖
.PHONY: tidy
tidy:
	@echo "整理Go模块依赖..."
	go mod tidy

# 清理构建产物
.PHONY: clean
clean:
	@echo "清理构建产物..."
	rm -rf $(BUILD_DIR)

# 运行用户服务（开发模式）
.PHONY: run-user-svc
run-user-svc:
	go run cmd/user-svc/main.go

# 运行训练服务（开发模式）
.PHONY: run-training-svc
run-training-svc:
	go run cmd/training-svc/main.go

# 运行视力服务（开发模式）
.PHONY: run-vision-svc
run-vision-svc:
	go run cmd/vision-svc/main.go

# 运行激励服务（开发模式）
.PHONY: run-reward-svc
run-reward-svc:
	go run cmd/reward-svc/main.go

# 运行API网关（开发模式）
.PHONY: run-gateway
run-gateway:
	go run cmd/gateway/main.go
