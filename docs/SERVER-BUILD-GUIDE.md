# KidsFit 小勇士 - 服务端编译安装文档

## 1. 环境要求

| 工具 | 最低版本 | 说明 |
|------|---------|------|
| Go | 1.21+ | 推荐 1.22+ |
| PostgreSQL | 14+ | 推荐 16+ |
| Redis | 7.0+ | 推荐 7.2+ |
| golang-migrate | 4.x | 数据库迁移工具 |
| Make | GNU Make | 构建工具 |
| Git | 2.x | 版本管理 |

## 2. 安装基础依赖

### 2.1 macOS（Homebrew）

```bash
# 安装 Go
brew install go

# 安装 PostgreSQL
brew install postgresql@16
brew services start postgresql@16

# 安装 Redis
brew install redis
brew services start redis

# 安装 golang-migrate
brew install golang-migrate

# 验证
go version
psql --version
redis-cli ping
migrate -version
```

### 2.2 Ubuntu/Debian

```bash
# 安装 Go
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 安装 PostgreSQL
sudo apt install postgresql-16 postgresql-client-16
sudo systemctl start postgresql

# 安装 Redis
sudo apt install redis-server
sudo systemctl start redis-server

# 安装 golang-migrate
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/
```

### 2.3 配置 Go 国内代理（加速模块下载）

```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

## 3. 数据库配置

### 3.1 创建数据库和用户

```bash
# 登录 PostgreSQL
sudo -u postgres psql

# 执行以下 SQL
CREATE USER kidsfit WITH PASSWORD 'kidsfit_dev';
CREATE DATABASE kidsfit_users OWNER kidsfit;
CREATE DATABASE kidsfit_training OWNER kidsfit;
CREATE DATABASE kidsfit_vision OWNER kidsfit;
CREATE DATABASE kidsfit_rewards OWNER kidsfit;

# 启用 UUID 扩展
\c kidsfit_users
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
\c kidsfit_training
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
\c kidsfit_vision
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
\c kidsfit_rewards
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

\q
```

### 3.2 执行数据库迁移

```bash
cd kidsfit-backend

# 用户库迁移
migrate -path scripts/migrate \
  -database "postgres://kidsfit:kidsfit_dev@localhost:5432/kidsfit_users?sslmode=disable" \
  up

# 训练库迁移（需指定对应迁移文件）
migrate -path scripts/migrate \
  -database "postgres://kidsfit:kidsfit_dev@localhost:5432/kidsfit_training?sslmode=disable" \
  up

# 视力库迁移
migrate -path scripts/migrate \
  -database "postgres://kidsfit:kidsfit_dev@localhost:5432/kidsfit_vision?sslmode=disable" \
  up

# 激励库迁移
migrate -path scripts/migrate \
  -database "postgres://kidsfit:kidsfit_dev@localhost:5432/kidsfit_rewards?sslmode=disable" \
  up
```

或使用 Makefile：

```bash
make migrate-up
```

### 3.3 验证迁移

```bash
psql -U kidsfit -d kidsfit_users -c "\dt"
# 应显示：users, families, parent_settings

psql -U kidsfit -d kidsfit_training -c "\dt"
# 应显示：exercise_records, training_plans, exercise_items, fitness_assessments

psql -U kidsfit -d kidsfit_vision -c "\dt"
# 应显示：vision_records, outdoor_activities, outdoor_segments, eye_reminders

psql -U kidsfit -d kidsfit_rewards -c "\dt"
# 应显示：badges, user_badges, point_records, challenges
```

## 4. 项目配置

### 4.1 进入项目目录

```bash
cd kidsfit-backend
```

### 4.2 获取依赖

```bash
go mod tidy
```

### 4.3 配置文件

配置文件位于 `configs/` 目录，按环境区分：

| 文件 | 用途 |
|------|------|
| `configs/dev.yaml` | 开发环境 |
| `configs/prod.yaml` | 生产环境（需自行创建） |
| `configs/test.yaml` | 测试环境（需自行创建） |

#### 开发环境配置 (`configs/dev.yaml`)

```yaml
server:
  port: 8001
  mode: debug
database:
  host: localhost
  port: 5432
  user: kidsfit
  password: kidsfit_dev
  dbname: kidsfit_users
  sslmode: disable
  max_open_conns: 100
  max_idle_conns: 10
  conn_timeout: 30
redis:
  addr: localhost:6379
  password: ""
  db: 0
jwt:
  secret: kidsfit_dev_secret_key_2024
  access_ttl: 7200
  refresh_ttl: 604800
  issuer: kidsfit
log:
  level: debug
  format: console
```

#### 生产环境配置要点

```yaml
server:
  mode: release
database:
  sslmode: require
  password: <强密码>
redis:
  password: <强密码>
jwt:
  secret: <随机64位密钥>
  access_ttl: 1800      # 30分钟
  refresh_ttl: 604800   # 7天
log:
  level: info
  format: json
```

### 4.4 环境变量覆盖

支持 `KIDSFIT_` 前缀的环境变量覆盖配置：

```bash
export KIDSFIT_SERVER_PORT=8001
export KIDSFIT_DATABASE_HOST=db.internal
export KIDSFIT_DATABASE_PASSWORD=prod_password
export KIDSFIT_REDIS_ADDR=redis.internal:6379
export KITSFIT_JWT_SECRET=production_secret_key
```

## 5. 编译构建

### 5.1 编译所有服务

```bash
make build
```

输出位置：`bin/user-svc`

### 5.2 手动编译单个服务

```bash
# 编译 user-svc
go build -ldflags "-s -w" -o bin/user-svc ./cmd/user-svc/

# 编译时注入版本信息
VERSION=$(git describe --tags --always)
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
go build -ldflags "-s -w -X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME" \
  -o bin/user-svc ./cmd/user-svc/
```

### 5.3 交叉编译

```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o bin/user-svc-linux ./cmd/user-svc/

# ARM64 (Apple Silicon / AWS Graviton)
GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o bin/user-svc-arm64 ./cmd/user-svc/
```

## 6. 运行服务

### 6.1 开发模式（直接运行）

```bash
make run-user-svc
# 或
go run cmd/user-svc/main.go
```

### 6.2 生产模式（编译后运行）

```bash
./bin/user-svc
```

### 6.3 指定配置文件

```bash
./bin/user-svc -config configs/prod.yaml
```

### 6.4 验证服务启动

```bash
# 健康检查
curl http://localhost:8001/api/v1/auth/login

# 应返回：
# {"code":401,"message":"请求参数错误","data":null}
```

### 6.5 服务端口分配

| 服务 | 端口 | 说明 |
|------|------|------|
| user-svc | 8001 | 用户认证与管理 |
| training-svc | 8002 | 运动训练（待开发） |
| vision-svc | 8003 | 视力健康（待开发） |
| reward-svc | 8004 | 激励系统（待开发） |
| gateway | 8080 | API 网关（待开发） |

## 7. 运行测试

### 7.1 运行全部测试

```bash
make test
# 或
go test -v -race ./...
```

### 7.2 运行指定包测试

```bash
# 领域层测试
go test -v ./internal/domain/...

# 应用层测试
go test -v ./internal/application/...

# 公共包测试
go test -v ./internal/pkg/...
```

### 7.3 运行基准测试

```bash
go test -bench=. -benchmem ./internal/pkg/crypto/...
go test -bench=. -benchmem ./internal/pkg/jwt/...
```

### 7.4 测试覆盖率

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## 8. 代码检查

### 8.1 运行 lint

```bash
make lint
# 或
golangci-lint run ./...
```

### 8.2 运行 vet

```bash
go vet ./...
```

## 9. 项目结构

```
kidsfit-backend/
├── api/
│   ├── http/
│   │   ├── handler/              # HTTP 处理器（4个模块）
│   │   ├── middleware/            # 中间件（认证/CORS/限流/日志）
│   │   └── router.go             # 路由注册
│   └── proto/                    # gRPC Proto 定义（待开发）
├── cmd/
│   └── user-svc/main.go          # 服务入口
├── configs/
│   └── dev.yaml                  # 开发环境配置
├── internal/
│   ├── domain/                   # 领域层（实体 + 仓储接口）
│   │   ├── user/                 # 用户域
│   │   ├── training/             # 训练域
│   │   ├── vision/               # 视力域
│   │   └── reward/               # 激励域
│   ├── application/              # 应用层（DTO + 应用服务）
│   │   ├── user/
│   │   ├── training/
│   │   ├── vision/
│   │   └── reward/
│   ├── infrastructure/           # 基础设施层
│   │   └── persistence/
│   │       ├── postgresql/       # PostgreSQL 仓储实现
│   │       └── redis/            # Redis 缓存/锁/排行榜
│   └── pkg/                      # 公共包
│       ├── config/               # 配置管理
│       ├── crypto/               # AES-256-GCM + bcrypt
│       ├── errors/               # 统一错误码
│       ├── jwt/                  # JWT 工具
│       ├── logger/               # 日志（待完善）
│       ├── pagination/           # 分页工具
│       ├── response/             # 统一响应
│       └── validator/            # 参数校验（待完善）
├── scripts/
│   └── migrate/                  # 数据库迁移脚本（4个）
├── tests/                        # 集成/E2E/性能测试
├── go.mod / go.sum
└── Makefile
```

## 10. API 接口清单

### 10.1 用户服务 (user-svc:8001)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | /api/v1/auth/register | 家长注册 | 否 |
| POST | /api/v1/auth/login | 登录 | 否 |
| POST | /api/v1/auth/refresh | 刷新Token | 是 |
| POST | /api/v1/auth/logout | 登出 | 是 |
| GET | /api/v1/users/me | 获取当前用户 | 是 |
| PUT | /api/v1/users/me | 更新用户信息 | 是 |
| POST | /api/v1/users/avatar | 上传头像 | 是 |
| POST | /api/v1/children | 添加儿童 | 是 |
| GET | /api/v1/children | 获取儿童列表 | 是 |
| GET | /api/v1/children/:id | 获取儿童详情 | 是 |
| PUT | /api/v1/children/:id | 更新儿童信息 | 是 |
| DELETE | /api/v1/children/:id | 删除儿童 | 是 |
| GET | /api/v1/settings | 获取家长设置 | 是 |
| PUT | /api/v1/settings | 更新家长设置 | 是 |

### 10.2 训练服务 (training-svc:8002)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | /api/v1/exercises | 提交运动记录 | 是 |
| GET | /api/v1/exercises | 获取运动记录 | 是 |
| GET | /api/v1/exercises/personal-best | 个人最佳 | 是 |
| GET | /api/v1/plans/today | 今日训练计划 | 是 |
| POST | /api/v1/plans/:id/complete | 完成计划 | 是 |
| POST | /api/v1/assessments | 提交体能评估 | 是 |
| GET | /api/v1/assessments/latest | 最新评估 | 是 |
| GET | /api/v1/stats/weekly | 周统计 | 是 |
| GET | /api/v1/stats/monthly | 月统计 | 是 |

### 10.3 视力服务 (vision-svc:8003)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | /api/v1/vision-records | 添加视力记录 | 是 |
| GET | /api/v1/vision-records | 获取视力记录 | 是 |
| POST | /api/v1/vision-records/ocr | OCR识别验光单 | 是 |
| GET | /api/v1/vision/trend | 视力趋势 | 是 |
| GET | /api/v1/outdoor/today | 今日户外时间 | 是 |
| POST | /api/v1/outdoor/sync | 同步户外数据 | 是 |
| GET | /api/v1/reminders | 获取提醒 | 是 |
| POST | /api/v1/reminders/:id/ack | 确认提醒 | 是 |

### 10.4 激励服务 (reward-svc:8004)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | /api/v1/badges | 勋章列表 | 是 |
| GET | /api/v1/badges/my | 我的勋章 | 是 |
| GET | /api/v1/points | 积分记录 | 是 |
| GET | /api/v1/points/balance | 积分余额 | 是 |
| POST | /api/v1/challenges | 发起挑战 | 是 |
| GET | /api/v1/challenges | 挑战列表 | 是 |
| POST | /api/v1/challenges/:id/accept | 接受挑战 | 是 |
| POST | /api/v1/challenges/:id/submit | 提交成绩 | 是 |
| GET | /api/v1/leaderboard/family | 家庭排行 | 是 |
| GET | /api/v1/leaderboard/global | 全局排行 | 是 |

## 11. Docker 部署（推荐生产环境）

### 11.1 创建 Dockerfile

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /user-svc ./cmd/user-svc/

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /user-svc /usr/local/bin/
COPY configs/ /app/configs/
EXPOSE 8001
CMD ["user-svc"]
```

### 11.2 构建镜像

```bash
docker build -t kidsfit/user-svc:latest .
```

### 11.3 运行容器

```bash
docker run -d \
  --name kidsfit-user-svc \
  -p 8001:8001 \
  -e KIDSFIT_DATABASE_HOST=host.docker.internal \
  -e KIDSFIT_REDIS_ADDR=host.docker.internal:6379 \
  kidsfit/user-svc:latest
```

## 12. 常见问题

### 12.1 数据库连接失败

```bash
# 检查 PostgreSQL 状态
pg_isready -h localhost -p 5432

# 检查用户权限
psql -U kidsfit -d kidsfit_users -c "SELECT 1;"

# 检查 pg_hba.conf 允许本地连接
```

### 12.2 Redis 连接失败

```bash
redis-cli ping
# 应返回 PONG

# 检查 Redis 配置
redis-cli info server
```

### 12.3 Go 模块下载慢

```bash
go env -w GOPROXY=https://goproxy.cn,direct
go clean -modcache
go mod tidy
```

### 12.4 迁移失败

```bash
# 查看迁移版本
migrate -path scripts/migrate \
  -database "postgres://kidsfit:kidsfit_dev@localhost:5432/kidsfit_users?sslmode=disable" \
  version

# 回滚一步
migrate -path scripts/migrate \
  -database "postgres://kidsfit:kidsfit_dev@localhost:5432/kidsfit_users?sslmode=disable" \
  down 1

# 强制设置版本（谨慎使用）
migrate -path scripts/migrate \
  -database "postgres://kidsfit:kidsfit_dev@localhost:5432/kidsfit_users?sslmode=disable" \
  force 4
```

### 12.5 端口被占用

```bash
lsof -i :8001
kill -9 <PID>
```

## 13. 安全注意事项

| 项目 | 要求 |
|------|------|
| JWT Secret | 生产环境必须使用随机64位密钥 |
| 数据库密码 | 生产环境必须使用强密码 |
| Redis 密码 | 生产环境必须设置密码 |
| HTTPS | 生产环境必须启用 TLS 1.3 |
| 敏感数据 | 手机号/密码使用 AES-256-GCM + bcrypt 加密 |
| 摄像头数据 | 不落地存储，仅内存处理 |
| CORS | 生产环境限制允许的域名 |
| 限流 | 已实现 Redis 滑动窗口限流中间件 |
