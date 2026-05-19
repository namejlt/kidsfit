# KidsFit 小勇士 — 技术设计文档 (Technical Design Document)

> 版本：v1.0
> 日期：2026-05-19
> 状态：设计定稿
> 对应PRD：docs/PRD.md

---

## 一、设计目标与原则

### 1.1 设计目标

- 支撑3-12岁儿童及家长的AI运动健康服务，峰值QPS 10,000
- 保障儿童隐私数据安全，符合国内儿童个人信息保护法规
- 支持AI骨骼识别实时反馈（端到端延迟 < 200ms）
- 系统可用性达到99.9%，支持水平扩展

### 1.2 设计原则

| 原则 | 说明 |
|-----|------|
| **隐私优先** | 敏感数据本地加密，云端脱敏，摄像头数据不落地 |
| **高并发** | Go协程处理连接，连接池管理数据库，缓存抗读热点 |
| **可扩展** | 微服务拆分，独立部署，服务间通过gRPC/HTTP通信 |
| **容错性** | 核心功能支持离线，服务降级策略，熔断限流 |
| **可观测** | 全链路日志、指标监控、分布式追踪 |

---

## 二、整体技术架构

### 2.1 技术栈总览

| 层级 | 技术选型 | 版本/说明 |
|-----|---------|----------|
| 跨平台客户端 | Flutter | 3.x，Dart |
| 客户端状态管理 | Riverpod | 类型安全 |
| 客户端本地存储 | Hive + SQLite | Hive配置/缓存，SQLite结构化数据 |
| 客户端AI/ML | Google ML Kit + TFLite | Pose Detection + 自定义模型 |
| 客户端AR | ARCore | 体感交互 |
| **后端服务** | **Go + Gin** | **v1.21+，高性能HTTP框架** |
| 后端RPC | gRPC + Protocol Buffers | 服务间通信 |
| 数据库 | PostgreSQL 15 | 主业务数据 |
| 缓存 | Redis 7 | 会话、缓存、限流、排行榜 |
| 消息队列 | RabbitMQ / Kafka | 异步任务、事件驱动 |
| 对象存储 | 阿里云OSS / 腾讯云COS | 图片、视频、验光单 |
| 搜索引擎 | Elasticsearch 8 | 内容搜索、日志检索 |
| 推送服务 | Firebase Cloud Messaging | 跨平台推送 |
| 监控告警 | Prometheus + Grafana + Jaeger | 指标、日志、追踪 |
| 容器编排 | Kubernetes | 容器化部署 |
| API网关 | Kong / Nginx Ingress | 路由、限流、认证 |

### 2.2 系统拓扑图

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              接入层 (Access Layer)                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌─────────────┐    ┌─────────────┐    ┌─────────────┐                    │
│   │   CDN加速    │    │  API网关     │    │  WAF防火墙   │                    │
│   │  (静态资源)  │    │  (Kong)     │    │  (防护)     │                    │
│   └──────┬──────┘    └──────┬──────┘    └─────────────┘                    │
│          │                  │                                               │
│          └──────────────────┘                                               │
│                     │                                                       │
└─────────────────────┼───────────────────────────────────────────────────────┘
                      │ HTTPS/TLS 1.3
                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                            网关层 (Gateway Layer)                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │                         API Gateway (Kong)                           │  │
│   │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐  │  │
│   │  │ 路由转发  │ │ JWT认证  │ │ 速率限制  │ │ 请求鉴权  │ │ 日志记录  │  │  │
│   │  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘  │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
                      │
        ┌─────────────┼─────────────┬─────────────┐
        ▼             ▼             ▼             ▼
┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐
│  用户服务    │ │  训练服务    │ │  视力服务    │ │  激励服务    │
│  (user-svc) │ │(training-svc│ │(vision-svc) │ │(reward-svc) │
│   Go/Gin    │ │   Go/Gin    │ │   Go/Gin    │ │   Go/Gin    │
└──────┬──────┘ └──────┬──────┘ └──────┬──────┘ └──────┬──────┘
       │               │               │               │
       └───────────────┴───────────────┴───────────────┘
                           │
              ┌────────────┼────────────┐
              ▼            ▼            ▼
       ┌──────────┐ ┌──────────┐ ┌──────────┐
       │PostgreSQL│ │  Redis   │ │  OSS/COS │
       │  (主库)   │ │ (缓存)   │ │ (文件)   │
       └──────────┘ └──────────┘ └──────────┘
              │            │            │
              └────────────┴────────────┘
                           │
              ┌────────────┼────────────┐
              ▼            ▼            ▼
       ┌──────────┐ ┌──────────┐ ┌──────────┐
       │RabbitMQ/ │ │   ES     │ │Prometheus│
       │  Kafka   │ │ (搜索)   │ │ (监控)   │
       └──────────┘ └──────────┘ └──────────┘
```

### 2.3 服务拆分策略

采用**领域驱动设计(DDD)**进行服务拆分，按业务边界划分为独立微服务：

| 服务名 | 职责 | 端口 | 数据库 |
|-------|------|------|--------|
| **user-svc** | 用户管理、家长-儿童关系、认证授权 | 8001 | users |
| **training-svc** | 训练计划、运动记录、AI评分、课程管理 | 8002 | training |
| **vision-svc** | 视力档案、户外追踪、用眼提醒、OCR | 8003 | vision |
| **reward-svc** | 积分、勋章、成就、排行榜、亲子挑战 | 8004 | rewards |
| **content-svc** | 技能课程、故事线、教学视频、内容配置 | 8005 | content |
| **notification-svc** | 推送通知、短信、站内消息 | 8006 | - |
| **gateway** | API网关、统一认证、限流、路由 | 8080 | - |

---

## 三、Go后端服务架构设计

### 3.1 项目目录结构

```
kidsfit-backend/
├── api/                          # API定义
│   ├── proto/                    # Protocol Buffers定义
│   │   ├── user.proto
│   │   ├── training.proto
│   │   ├── vision.proto
│   │   └── reward.proto
│   └── http/                     # HTTP API路由定义
│       ├── router.go
│       ├── middleware/
│       │   ├── auth.go           # JWT认证中间件
│       │   ├── rate_limit.go     # 限流中间件
│       │   ├── cors.go           # 跨域中间件
│       │   └── logger.go         # 请求日志中间件
│       └── handler/
│           ├── user_handler.go
│           ├── training_handler.go
│           ├── vision_handler.go
│           └── reward_handler.go
│
├── cmd/                          # 服务入口
│   ├── user-svc/                 # 用户服务
│   │   └── main.go
│   ├── training-svc/             # 训练服务
│   │   └── main.go
│   ├── vision-svc/               # 视力服务
│   │   └── main.go
│   ├── reward-svc/               # 激励服务
│   │   └── main.go
│   └── gateway/                  # API网关
│       └── main.go
│
├── internal/                     # 内部实现（不可外部导入）
│   ├── domain/                   # 领域模型
│   │   ├── user/
│   │   │   ├── model.go          # 领域实体
│   │   │   ├── repository.go     # 仓储接口
│   │   │   └── service.go        # 领域服务
│   │   ├── training/
│   │   ├── vision/
│   │   └── reward/
│   │
│   ├── application/              # 应用层
│   │   ├── user/
│   │   │   ├── dto.go            # 数据传输对象
│   │   │   ├── app_service.go    # 应用服务
│   │   │   └── assembler.go      # DTO-领域模型转换
│   │   ├── training/
│   │   ├── vision/
│   │   └── reward/
│   │
│   ├── infrastructure/           # 基础设施层
│   │   ├── persistence/          # 数据持久化
│   │   │   ├── postgresql/       # PostgreSQL实现
│   │   │   │   ├── user_repo.go
│   │   │   │   ├── training_repo.go
│   │   │   │   └── db.go         # 数据库连接管理
│   │   │   └── redis/            # Redis实现
│   │   │       ├── cache.go
│   │   │       ├── session.go
│   │   │       └── lock.go
│   │   ├── messaging/            # 消息队列
│   │   │   ├── rabbitmq/
│   │   │   └── kafka/
│   │   ├── storage/              # 对象存储
│   │   │   └── oss.go
│   │   ├── ai/                   # AI服务客户端
│   │   │   └── pose_client.go
│   │   └── ocr/                  # OCR服务客户端
│   │       └── ocr_client.go
│   │
│   └── pkg/                      # 公共包
│       ├── config/               # 配置管理
│       │   └── config.go
│       ├── logger/               # 日志
│       │   └── zap_logger.go
│       ├── errors/               # 错误码
│       │   └── errors.go
│       ├── validator/            # 参数校验
│       │   └── validator.go
│       ├── crypto/               # 加密
│       │   ├── aes.go
│       │   └── hash.go
│       ├── jwt/                  # JWT工具
│       │   └── jwt.go
│       ├── pagination/           # 分页
│       │   └── pagination.go
│       └── response/             # 统一响应
│           └── response.go
│
├── configs/                      # 配置文件
│   ├── dev.yaml
│   ├── prod.yaml
│   └── test.yaml
│
├── deployments/                  # 部署配置
│   ├── docker/
│   │   ├── Dockerfile
│   │   └── docker-compose.yml
│   └── k8s/
│       ├── namespace.yaml
│       ├── configmap.yaml
│       ├── secret.yaml
│       ├── deployment.yaml
│       ├── service.yaml
│       └── ingress.yaml
│
├── scripts/                      # 脚本
│   ├── migrate/                  # 数据库迁移
│   │   └── migrate.go
│   └── seed/                     # 数据初始化
│       └── seed.go
│
├── tests/                        # 测试
│   ├── integration/              # 集成测试
│   ├── e2e/                      # 端到端测试
│   └── benchmark/                # 性能测试
│
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### 3.2 分层架构设计

采用**整洁架构(Clean Architecture)**，依赖关系由外向内：

```
┌─────────────────────────────────────────┐
│           接口层 (Interface)              │  HTTP Handler / gRPC Server
│         api/http/handler/*.go             │
├─────────────────────────────────────────┤
│           应用层 (Application)            │  DTO、应用服务、事务控制
│       internal/application/*/             │
├─────────────────────────────────────────┤
│           领域层 (Domain)                 │  实体、值对象、仓储接口、领域服务
│          internal/domain/*/               │
├─────────────────────────────────────────┤
│         基础设施层 (Infrastructure)        │  数据库、缓存、消息队列、外部服务
│       internal/infrastructure/*/          │
└─────────────────────────────────────────┘
```

**依赖规则**：
- 内层不依赖外层
- 外层通过接口依赖内层
- 领域层完全独立，不依赖任何框架

### 3.3 核心领域模型

#### 3.3.1 用户领域 (User Domain)

```go
// internal/domain/user/model.go

package user

import (
	"time"
	"github.com/google/uuid"
)

// UserType 用户类型
type UserType string

const (
	UserTypeParent UserType = "parent"
	UserTypeChild  UserType = "child"
)

// User 用户实体
type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Type      UserType  `json:"type" gorm:"type:varchar(20);not null;index"`
	ParentID  *uuid.UUID `json:"parent_id,omitempty" gorm:"type:uuid;index"`
	Age       int       `json:"age" gorm:"type:int;check:age >= 3 AND age <= 12"`
	Nickname  string    `json:"nickname" gorm:"type:varchar(50);not null"`
	Avatar    string    `json:"avatar" gorm:"type:varchar(255)"`
	Phone     string    `json:"phone,omitempty" gorm:"type:varchar(20);index"` // 仅家长
	Password  string    `json:"-" gorm:"type:varchar(255)"` // 加密存储
	Status    UserStatus `json:"status" gorm:"type:varchar(20);default:active"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// UserStatus 用户状态
type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusDeleted  UserStatus = "deleted"
)

// Family 家庭关系
type Family struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	ParentID  uuid.UUID `json:"parent_id" gorm:"type:uuid;not null;index"`
	ChildID   uuid.UUID `json:"child_id" gorm:"type:uuid;not null;index"`
	Relation  string    `json:"relation" gorm:"type:varchar(20)"` // father, mother, etc.
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// ParentSettings 家长设置
type ParentSettings struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	ParentID        uuid.UUID `json:"parent_id" gorm:"type:uuid;not null;index"`
	DailyLimitMin   int       `json:"daily_limit_min" gorm:"type:int;default:30"`
	AvailableFrom   string    `json:"available_from" gorm:"type:varchar(5);default:08:00"`
	AvailableTo     string    `json:"available_to" gorm:"type:varchar(5);default:21:00"`
	CameraAllowed   bool      `json:"camera_allowed" gorm:"default:true"`
	LocationAllowed bool      `json:"location_allowed" gorm:"default:true"`
	DataUploadCloud bool      `json:"data_upload_cloud" gorm:"default:false"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
```

#### 3.3.2 训练领域 (Training Domain)

```go
// internal/domain/training/model.go

package training

import (
	"time"
	"github.com/google/uuid"
)

// ExerciseType 运动类型
type ExerciseType string

const (
	ExerciseJumpRope     ExerciseType = "jump_rope"
	ExerciseJumpingJack  ExerciseType = "jumping_jack"
	ExerciseSquat        ExerciseType = "squat"
	ExerciseSitUp        ExerciseType = "sit_up"
	ExerciseHighKnee     ExerciseType = "high_knee"
	ExercisePushUp       ExerciseType = "push_up"
)

// ExerciseRecord 运动记录
type ExerciseRecord struct {
	ID               uuid.UUID    `json:"id" gorm:"type:uuid;primary_key"`
	UserID           uuid.UUID    `json:"user_id" gorm:"type:uuid;not null;index"`
	Type             ExerciseType `json:"type" gorm:"type:varchar(30);not null;index"`
	DurationSeconds  int          `json:"duration_seconds" gorm:"type:int;not null"`
	Count            int          `json:"count" gorm:"type:int;default:0"`
	Score            int          `json:"score" gorm:"type:int;check:score >= 0 AND score <= 100"`
	RhythmScore      int          `json:"rhythm_score" gorm:"type:int"`
	AmplitudeScore   int          `json:"amplitude_score" gorm:"type:int"`
	SymmetryScore    int          `json:"symmetry_score" gorm:"type:int"`
	ContinuityScore  int          `json:"continuity_score" gorm:"type:int"`
	Corrections      []string     `json:"corrections" gorm:"type:text[]"`
	IsOffline        bool         `json:"is_offline" gorm:"default:false"`
	StartedAt        time.Time    `json:"started_at" gorm:"not null;index"`
	CompletedAt      time.Time    `json:"completed_at"`
	CreatedAt        time.Time    `json:"created_at" gorm:"autoCreateTime"`
}

// TrainingPlan 训练计划
type TrainingPlan struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primary_key"`
	UserID          uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	Date            time.Time      `json:"date" gorm:"type:date;not null;index"`
	Status          PlanStatus     `json:"status" gorm:"type:varchar(20);default:pending"`
	TotalDuration   int            `json:"total_duration" gorm:"type:int;not null"`
	ActualDuration  *int           `json:"actual_duration,omitempty" gorm:"type:int"`
	WarmupItems     []ExerciseItem `json:"warmup_items" gorm:"foreignKey:PlanID"`
	MainItems       []ExerciseItem `json:"main_items" gorm:"foreignKey:PlanID"`
	CooldownItems   []ExerciseItem `json:"cooldown_items" gorm:"foreignKey:PlanID"`
	CompletedAt     *time.Time     `json:"completed_at,omitempty"`
	CreatedAt       time.Time      `json:"created_at" gorm:"autoCreateTime"`
}

// PlanStatus 计划状态
type PlanStatus string

const (
	PlanStatusPending    PlanStatus = "pending"
	PlanStatusCompleted  PlanStatus = "completed"
	PlanStatusSkipped    PlanStatus = "skipped"
)

// ExerciseItem 训练项目
type ExerciseItem struct {
	ID          uuid.UUID    `json:"id" gorm:"type:uuid;primary_key"`
	PlanID      uuid.UUID    `json:"plan_id" gorm:"type:uuid;not null;index"`
	Type        ExerciseType `json:"type" gorm:"type:varchar(30);not null"`
	Name        string       `json:"name" gorm:"type:varchar(100);not null"`
	DurationSec int          `json:"duration_sec" gorm:"type:int"`
	TargetCount int          `json:"target_count" gorm:"type:int"`
	Difficulty  int          `json:"difficulty" gorm:"type:int;check:difficulty >= 1 AND difficulty <= 5"`
	Tips        string       `json:"tips" gorm:"type:varchar(500)"`
	Order       int          `json:"order" gorm:"type:int;not null"`
	Phase       ExercisePhase `json:"phase" gorm:"type:varchar(20);not null"`
}

// ExercisePhase 训练阶段
type ExercisePhase string

const (
	PhaseWarmup   ExercisePhase = "warmup"
	PhaseMain     ExercisePhase = "main"
	PhaseCooldown ExercisePhase = "cooldown"
)

// FitnessAssessment 体能评估
type FitnessAssessment struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	Endurance   int       `json:"endurance" gorm:"type:int;check:score >= 1 AND score <= 10"`
	Agility     int       `json:"agility" gorm:"type:int"`
	Strength    int       `json:"strength" gorm:"type:int"`
	Speed       int       `json:"speed" gorm:"type:int"`
	Coordination int      `json:"coordination" gorm:"type:int"`
	Balance     int       `json:"balance" gorm:"type:int"`
	Flexibility int       `json:"flexibility" gorm:"type:int"`
	AssessedAt  time.Time `json:"assessed_at" gorm:"not null;index"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
}
```

#### 3.3.3 视力领域 (Vision Domain)

```go
// internal/domain/vision/model.go

package vision

import (
	"time"
	"github.com/google/uuid"
)

// VisionRecord 视力档案记录
type VisionRecord struct {
	ID                uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	UserID            uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	ChildID           uuid.UUID `json:"child_id" gorm:"type:uuid;not null;index"`
	Date              time.Time `json:"date" gorm:"type:date;not null;index"`
	RightEye          EyeData   `json:"right_eye" gorm:"embedded;embeddedPrefix:right_"`
	LeftEye           EyeData   `json:"left_eye" gorm:"embedded;embeddedPrefix:left_"`
	AxialLengthRight  *float64  `json:"axial_length_right,omitempty" gorm:"type:decimal(5,2)"`
	AxialLengthLeft   *float64  `json:"axial_length_left,omitempty" gorm:"type:decimal(5,2)"`
	HyperopiaReserve  *float64  `json:"hyperopia_reserve,omitempty" gorm:"type:decimal(4,2)"`
	Source            DataSource `json:"source" gorm:"type:varchar(20);not null"`
	ImageURL          *string   `json:"image_url,omitempty" gorm:"type:varchar(255)"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// EyeData 单眼数据
type EyeData struct {
	SPH float64 `json:"sph" gorm:"type:decimal(5,2)"` // 球镜
	CYL float64 `json:"cyl" gorm:"type:decimal(5,2)"` // 柱镜
	AXIS int    `json:"axis" gorm:"type:int"`         // 轴位
	VA  float64 `json:"va" gorm:"type:decimal(3,2)"`  // 矫正视力
}

// DataSource 数据来源
type DataSource string

const (
	SourceOCR    DataSource = "ocr"
	SourceManual DataSource = "manual"
)

// OutdoorActivity 户外活动记录
type OutdoorActivity struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	Date        time.Time `json:"date" gorm:"type:date;not null;index"`
	DurationMin int       `json:"duration_min" gorm:"type:int;not null;default:0"`
	Segments    []OutdoorSegment `json:"segments" gorm:"foreignKey:ActivityID"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// OutdoorSegment 户外时段
type OutdoorSegment struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	ActivityID uuid.UUID `json:"activity_id" gorm:"type:uuid;not null;index"`
	StartTime  time.Time `json:"start_time" gorm:"not null"`
	EndTime    time.Time `json:"end_time" gorm:"not null"`
	DurationMin int      `json:"duration_min" gorm:"type:int;not null"`
	Location   *string   `json:"location,omitempty" gorm:"type:varchar(100)"`
}

// EyeReminder 用眼提醒记录
type EyeReminder struct {
	ID        uuid.UUID     `json:"id" gorm:"type:uuid;primary_key"`
	UserID    uuid.UUID     `json:"user_id" gorm:"type:uuid;not null;index"`
	Type      ReminderType  `json:"type" gorm:"type:varchar(30);not null"`
	TriggeredAt time.Time   `json:"triggered_at" gorm:"not null;index"`
	Acknowledged bool       `json:"acknowledged" gorm:"default:false"`
	CreatedAt time.Time     `json:"created_at" gorm:"autoCreateTime"`
}

// ReminderType 提醒类型
type ReminderType string

const (
	Reminder202020  ReminderType = "20_20_20"
	ReminderOutdoor ReminderType = "outdoor"
	ReminderBreak   ReminderType = "break"
)
```

#### 3.3.4 激励领域 (Reward Domain)

```go
// internal/domain/reward/model.go

package reward

import (
	"time"
	"github.com/google/uuid"
)

// Badge 勋章
type Badge struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	Code        string     `json:"code" gorm:"type:varchar(50);not null;uniqueIndex"`
	Name        string     `json:"name" gorm:"type:varchar(100);not null"`
	Description string     `json:"description" gorm:"type:varchar(500)"`
	Category    BadgeCategory `json:"category" gorm:"type:varchar(30);not null"`
	Icon        string     `json:"icon" gorm:"type:varchar(255)"`
	Condition   string     `json:"condition" gorm:"type:varchar(500)"` // 获得条件JSON
	Points      int        `json:"points" gorm:"type:int;default:0"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
}

// BadgeCategory 勋章分类
type BadgeCategory string

const (
	BadgeCategoryMilestone  BadgeCategory = "milestone"
	BadgeCategorySkill      BadgeCategory = "skill"
	BadgeCategoryStreak     BadgeCategory = "streak"
	BadgeCategoryChallenge  BadgeCategory = "challenge"
	BadgeCategoryFamily     BadgeCategory = "family"
	BadgeCategoryVision     BadgeCategory = "vision"
	BadgeCategorySpecial    BadgeCategory = "special"
)

// UserBadge 用户勋章
type UserBadge struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	BadgeID   uuid.UUID `json:"badge_id" gorm:"type:uuid;not null;index"`
	EarnedAt  time.Time `json:"earned_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// PointRecord 积分记录
type PointRecord struct {
	ID          uuid.UUID   `json:"id" gorm:"type:uuid;primary_key"`
	UserID      uuid.UUID   `json:"user_id" gorm:"type:uuid;not null;index"`
	Points      int         `json:"points" gorm:"type:int;not null"`
	Type        PointType   `json:"type" gorm:"type:varchar(30);not null"`
	SourceID    *uuid.UUID  `json:"source_id,omitempty" gorm:"type:uuid"`
	SourceType  *string     `json:"source_type,omitempty" gorm:"type:varchar(50)"`
	Description string      `json:"description" gorm:"type:varchar(255)"`
	Balance     int         `json:"balance" gorm:"type:int;not null"` // 变动后余额
	CreatedAt   time.Time   `json:"created_at" gorm:"autoCreateTime"`
}

// PointType 积分类型
type PointType string

const (
	PointTypeExercise   PointType = "exercise"
	PointTypeRecord     PointType = "record_break"
	PointTypeFamily     PointType = "family_activity"
	PointTypeStreak     PointType = "streak"
	PointTypeVision     PointType = "vision_task"
	PointTypeRedeem     PointType = "redeem"
)

// Challenge 挑战
type Challenge struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key"`
	Type        ChallengeType  `json:"type" gorm:"type:varchar(30);not null"`
	InitiatorID uuid.UUID      `json:"initiator_id" gorm:"type:uuid;not null;index"`
	AcceptorID  uuid.UUID      `json:"acceptor_id" gorm:"type:uuid;not null;index"`
	ExerciseType string        `json:"exercise_type" gorm:"type:varchar(30);not null"`
	TargetValue int           `json:"target_value" gorm:"type:int;not null"`
	InitiatorScore *int       `json:"initiator_score,omitempty" gorm:"type:int"`
	AcceptorScore  *int       `json:"acceptor_score,omitempty" gorm:"type:int"`
	WinnerID    *uuid.UUID    `json:"winner_id,omitempty" gorm:"type:uuid"`
	Status      ChallengeStatus `json:"status" gorm:"type:varchar(20);default:pending"`
	ExpiresAt   time.Time     `json:"expires_at" gorm:"not null"`
	CompletedAt *time.Time    `json:"completed_at,omitempty"`
	CreatedAt   time.Time     `json:"created_at" gorm:"autoCreateTime"`
}

// ChallengeType 挑战类型
type ChallengeType string

const (
	ChallengeTypeSync  ChallengeType = "sync"
	ChallengeTypeAsync ChallengeType = "async"
	ChallengeTypeTimed ChallengeType = "timed"
)

// ChallengeStatus 挑战状态
type ChallengeStatus string

const (
	ChallengeStatusPending   ChallengeStatus = "pending"
	ChallengeStatusAccepted  ChallengeStatus = "accepted"
	ChallengeStatusCompleted ChallengeStatus = "completed"
	ChallengeStatusExpired   ChallengeStatus = "expired"
)
```

### 3.4 仓储接口设计

```go
// internal/domain/user/repository.go

package user

import (
	"context"
	"github.com/google/uuid"
)

// Repository 用户仓储接口
type Repository interface {
	// Create 创建用户
	Create(ctx context.Context, user *User) error
	// GetByID 根据ID获取用户
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	// GetByPhone 根据手机号获取用户
	GetByPhone(ctx context.Context, phone string) (*User, error)
	// GetChildrenByParentID 获取家长下的所有儿童
	GetChildrenByParentID(ctx context.Context, parentID uuid.UUID) ([]*User, error)
	// Update 更新用户信息
	Update(ctx context.Context, user *User) error
	// Delete 软删除用户
	Delete(ctx context.Context, id uuid.UUID) error
	// List 分页查询用户
	List(ctx context.Context, filter UserFilter, page, pageSize int) ([]*User, int64, error)
}

// UserFilter 用户查询过滤条件
type UserFilter struct {
	Type     *UserType
	ParentID *uuid.UUID
	Status   *UserStatus
	Keyword  string
}
```

```go
// internal/domain/training/repository.go

package training

import (
	"context"
	"time"
	"github.com/google/uuid"
)

// Repository 训练仓储接口
type Repository interface {
	// ExerciseRecord
	CreateExerciseRecord(ctx context.Context, record *ExerciseRecord) error
	GetExerciseRecordByID(ctx context.Context, id uuid.UUID) (*ExerciseRecord, error)
	GetExerciseRecordsByUserID(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*ExerciseRecord, error)
	GetExerciseRecordsByUserIDAndType(ctx context.Context, userID uuid.UUID, exerciseType ExerciseType, limit int) ([]*ExerciseRecord, error)
	GetPersonalBest(ctx context.Context, userID uuid.UUID, exerciseType ExerciseType) (*ExerciseRecord, error)

	// TrainingPlan
	CreateTrainingPlan(ctx context.Context, plan *TrainingPlan) error
	GetTrainingPlanByID(ctx context.Context, id uuid.UUID) (*TrainingPlan, error)
	GetTrainingPlanByUserIDAndDate(ctx context.Context, userID uuid.UUID, date time.Time) (*TrainingPlan, error)
	UpdateTrainingPlan(ctx context.Context, plan *TrainingPlan) error

	// FitnessAssessment
	CreateFitnessAssessment(ctx context.Context, assessment *FitnessAssessment) error
	GetLatestFitnessAssessment(ctx context.Context, userID uuid.UUID) (*FitnessAssessment, error)
}
```

### 3.5 服务间通信

#### 3.5.1 gRPC服务定义

```protobuf
// api/proto/user.proto

syntax = "proto3";

package user;

option go_package = "github.com/kidsfit/api/proto/user";

service UserService {
  rpc GetUserByID(GetUserByIDRequest) returns (GetUserByIDResponse);
  rpc GetChildrenByParentID(GetChildrenRequest) returns (GetChildrenResponse);
  rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
}

message GetUserByIDRequest {
  string id = 1;
}

message GetUserByIDResponse {
  User user = 1;
}

message User {
  string id = 1;
  string type = 2;
  string nickname = 3;
  int32 age = 4;
  string avatar = 5;
}

message GetChildrenRequest {
  string parent_id = 1;
}

message GetChildrenResponse {
  repeated User children = 1;
}

message ValidateTokenRequest {
  string token = 1;
}

message ValidateTokenResponse {
  bool valid = 1;
  string user_id = 2;
  string user_type = 3;
}
```

#### 3.5.2 事件驱动（异步消息）

```go
// internal/infrastructure/messaging/events.go

package messaging

// EventType 事件类型
type EventType string

const (
	// 用户事件
	EventUserRegistered    EventType = "user.registered"
	EventUserLoggedIn      EventType = "user.logged_in"

	// 运动事件
	EventExerciseCompleted EventType = "exercise.completed"
	EventRecordBroken      EventType = "exercise.record_broken"

	// 视力事件
	EventOutdoorTargetMet  EventType = "vision.outdoor_target_met"
	EventVisionAlert       EventType = "vision.alert"

	// 激励事件
	EventBadgeEarned       EventType = "reward.badge_earned"
	EventPointsEarned      EventType = "reward.points_earned"
	EventChallengeCompleted EventType = "reward.challenge_completed"
)

// Event 事件结构
type Event struct {
	ID        string    `json:"id"`
	Type      EventType `json:"type"`
	Payload   []byte    `json:"payload"`
	Timestamp int64     `json:"timestamp"`
	Source    string    `json:"source"`
}
```

**事件流设计**：

```
运动完成事件流：

training-svc (EventExerciseCompleted)
    │
    ├──► reward-svc (计算积分、检查勋章)
    │
    ├──► notification-svc (推送完成通知)
    │
    └──► user-svc (更新用户统计)

记录突破事件流：

training-svc (EventRecordBroken)
    │
    ├──► reward-svc (+50积分、颁发突破勋章)
    │
    ├──► notification-svc (推送庆祝通知)
    │
    └──► content-svc (解锁新故事章节)
```

---

## 四、数据库设计

### 4.1 数据库选型与分库策略

| 数据库 | 用途 | 说明 |
|-------|------|------|
| PostgreSQL 15 | 主业务数据 | 用户、训练、视力、激励等结构化数据 |
| Redis 7 | 缓存/会话/排行榜 | 热点数据、分布式锁、实时排行榜 |
| Elasticsearch 8 | 搜索/日志 | 内容搜索、审计日志检索 |
| OSS/COS | 文件存储 | 头像、验光单照片、教学视频 |

**分库设计**：按服务拆分数据库，避免单库瓶颈

```
kidsfit_users      → user-svc
kidsfit_training   → training-svc
kidsfit_vision     → vision-svc
kidsfit_rewards    → reward-svc
kidsfit_content    → content-svc
```

### 4.2 核心表结构

#### 4.2.1 用户库 (kidsfit_users)

```sql
-- 用户表
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type VARCHAR(20) NOT NULL CHECK (type IN ('parent', 'child')),
    parent_id UUID REFERENCES users(id) ON DELETE CASCADE,
    age INT CHECK (age >= 3 AND age <= 12),
    nickname VARCHAR(50) NOT NULL,
    avatar VARCHAR(255),
    phone VARCHAR(20) UNIQUE,
    password_hash VARCHAR(255),
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'deleted')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_users_type ON users(type);
CREATE INDEX idx_users_parent_id ON users(parent_id);
CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_deleted_at ON users(deleted_at) WHERE deleted_at IS NULL;

-- 家庭关系表
CREATE TABLE families (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    child_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    relation VARCHAR(20) CHECK (relation IN ('father', 'mother', 'grandfather', 'grandmother', 'other')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(parent_id, child_id)
);

CREATE INDEX idx_families_parent_id ON families(parent_id);
CREATE INDEX idx_families_child_id ON families(child_id);

-- 家长设置表
CREATE TABLE parent_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    daily_limit_min INT DEFAULT 30 CHECK (daily_limit_min >= 5 AND daily_limit_min <= 120),
    available_from TIME DEFAULT '08:00',
    available_to TIME DEFAULT '21:00',
    camera_allowed BOOLEAN DEFAULT true,
    location_allowed BOOLEAN DEFAULT true,
    data_upload_cloud BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

#### 4.2.2 训练库 (kidsfit_training)

```sql
-- 运动记录表
CREATE TABLE exercise_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    type VARCHAR(30) NOT NULL CHECK (type IN ('jump_rope', 'jumping_jack', 'squat', 'sit_up', 'high_knee', 'push_up')),
    duration_seconds INT NOT NULL CHECK (duration_seconds > 0),
    count INT DEFAULT 0 CHECK (count >= 0),
    score INT CHECK (score >= 0 AND score <= 100),
    rhythm_score INT CHECK (rhythm_score >= 0 AND rhythm_score <= 100),
    amplitude_score INT CHECK (amplitude_score >= 0 AND amplitude_score <= 100),
    symmetry_score INT CHECK (symmetry_score >= 0 AND symmetry_score <= 100),
    continuity_score INT CHECK (continuity_score >= 0 AND continuity_score <= 100),
    corrections TEXT[],
    is_offline BOOLEAN DEFAULT false,
    started_at TIMESTAMPTZ NOT NULL,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_exercise_records_user_id ON exercise_records(user_id);
CREATE INDEX idx_exercise_records_user_type ON exercise_records(user_id, type);
CREATE INDEX idx_exercise_records_started_at ON exercise_records(started_at);
CREATE INDEX idx_exercise_records_user_date ON exercise_records(user_id, started_at);

-- 训练计划表
CREATE TABLE training_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    date DATE NOT NULL,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'skipped')),
    total_duration INT NOT NULL CHECK (total_duration > 0),
    actual_duration INT,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, date)
);

CREATE INDEX idx_training_plans_user_id ON training_plans(user_id);
CREATE INDEX idx_training_plans_user_date ON training_plans(user_id, date);

-- 训练项目表
CREATE TABLE exercise_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id UUID NOT NULL REFERENCES training_plans(id) ON DELETE CASCADE,
    type VARCHAR(30) NOT NULL,
    name VARCHAR(100) NOT NULL,
    duration_sec INT,
    target_count INT,
    difficulty INT CHECK (difficulty >= 1 AND difficulty <= 5),
    tips VARCHAR(500),
    item_order INT NOT NULL,
    phase VARCHAR(20) NOT NULL CHECK (phase IN ('warmup', 'main', 'cooldown'))
);

CREATE INDEX idx_exercise_items_plan_id ON exercise_items(plan_id);

-- 体能评估表
CREATE TABLE fitness_assessments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    endurance INT CHECK (endurance >= 1 AND endurance <= 10),
    agility INT CHECK (agility >= 1 AND agility <= 10),
    strength INT CHECK (strength >= 1 AND strength <= 10),
    speed INT CHECK (speed >= 1 AND speed <= 10),
    coordination INT CHECK (coordination >= 1 AND coordination <= 10),
    balance INT CHECK (balance >= 1 AND balance <= 10),
    flexibility INT CHECK (flexibility >= 1 AND flexibility <= 10),
    assessed_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_fitness_assessments_user_id ON fitness_assessments(user_id);
CREATE INDEX idx_fitness_assessments_user_assessed ON fitness_assessments(user_id, assessed_at);
```

#### 4.2.3 视力库 (kidsfit_vision)

```sql
-- 视力档案表
CREATE TABLE vision_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    child_id UUID NOT NULL,
    date DATE NOT NULL,
    right_sph DECIMAL(5,2),
    right_cyl DECIMAL(5,2),
    right_axis INT,
    right_va DECIMAL(3,2),
    left_sph DECIMAL(5,2),
    left_cyl DECIMAL(5,2),
    left_axis INT,
    left_va DECIMAL(3,2),
    axial_length_right DECIMAL(5,2),
    axial_length_left DECIMAL(5,2),
    hyperopia_reserve DECIMAL(4,2),
    source VARCHAR(20) NOT NULL CHECK (source IN ('ocr', 'manual')),
    image_url VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(child_id, date)
);

CREATE INDEX idx_vision_records_user_id ON vision_records(user_id);
CREATE INDEX idx_vision_records_child_id ON vision_records(child_id);
CREATE INDEX idx_vision_records_child_date ON vision_records(child_id, date);

-- 户外活动记录表
CREATE TABLE outdoor_activities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    date DATE NOT NULL,
    duration_min INT NOT NULL DEFAULT 0 CHECK (duration_min >= 0),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, date)
);

CREATE INDEX idx_outdoor_activities_user_id ON outdoor_activities(user_id);
CREATE INDEX idx_outdoor_activities_user_date ON outdoor_activities(user_id, date);

-- 户外时段表
CREATE TABLE outdoor_segments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    activity_id UUID NOT NULL REFERENCES outdoor_activities(id) ON DELETE CASCADE,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    duration_min INT NOT NULL CHECK (duration_min > 0),
    location VARCHAR(100)
);

CREATE INDEX idx_outdoor_segments_activity_id ON outdoor_segments(activity_id);

-- 用眼提醒记录表
CREATE TABLE eye_reminders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    type VARCHAR(30) NOT NULL CHECK (type IN ('20_20_20', 'outdoor', 'break')),
    triggered_at TIMESTAMPTZ NOT NULL,
    acknowledged BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_eye_reminders_user_id ON eye_reminders(user_id);
CREATE INDEX idx_eye_reminders_triggered_at ON eye_reminders(triggered_at);
```

#### 4.2.4 激励库 (kidsfit_rewards)

```sql
-- 勋章定义表
CREATE TABLE badges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500),
    category VARCHAR(30) NOT NULL CHECK (category IN ('milestone', 'skill', 'streak', 'challenge', 'family', 'vision', 'special')),
    icon VARCHAR(255),
    condition JSONB NOT NULL,
    points INT DEFAULT 0 CHECK (points >= 0),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_badges_category ON badges(category);

-- 用户勋章表
CREATE TABLE user_badges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    badge_id UUID NOT NULL REFERENCES badges(id),
    earned_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, badge_id)
);

CREATE INDEX idx_user_badges_user_id ON user_badges(user_id);
CREATE INDEX idx_user_badges_badge_id ON user_badges(badge_id);

-- 积分记录表
CREATE TABLE point_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    points INT NOT NULL,
    type VARCHAR(30) NOT NULL CHECK (type IN ('exercise', 'record_break', 'family_activity', 'streak', 'vision_task', 'redeem')),
    source_id UUID,
    source_type VARCHAR(50),
    description VARCHAR(255),
    balance INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_point_records_user_id ON point_records(user_id);
CREATE INDEX idx_point_records_user_created ON point_records(user_id, created_at);

-- 挑战表
CREATE TABLE challenges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type VARCHAR(30) NOT NULL CHECK (type IN ('sync', 'async', 'timed')),
    initiator_id UUID NOT NULL,
    acceptor_id UUID NOT NULL,
    exercise_type VARCHAR(30) NOT NULL,
    target_value INT NOT NULL CHECK (target_value > 0),
    initiator_score INT,
    acceptor_score INT,
    winner_id UUID,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'completed', 'expired')),
    expires_at TIMESTAMPTZ NOT NULL,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_challenges_initiator ON challenges(initiator_id);
CREATE INDEX idx_challenges_acceptor ON challenges(acceptor_id);
CREATE INDEX idx_challenges_status ON challenges(status);
```

### 4.3 数据库迁移策略

使用 **golang-migrate** 进行版本化迁移：

```
scripts/migrate/
├── 000001_create_users_table.up.sql
├── 000001_create_users_table.down.sql
├── 000002_create_training_tables.up.sql
├── 000002_create_training_tables.down.sql
├── 000003_create_vision_tables.up.sql
├── 000003_create_vision_tables.down.sql
├── 000004_create_reward_tables.up.sql
└── 000004_create_reward_tables.down.sql
```

---

## 五、API接口设计

### 5.1 接口规范

- **协议**：HTTPS/TLS 1.3
- **数据格式**：JSON
- **字符编码**：UTF-8
- **时间格式**：RFC3339 (2024-01-01T00:00:00Z)
- **版本控制**：URL路径版本 `/api/v1/...`
- **认证方式**：JWT Bearer Token

### 5.2 统一响应格式

```json
// 成功响应
{
  "code": 0,
  "message": "success",
  "data": { ... }
}

// 错误响应
{
  "code": 10001,
  "message": "用户不存在",
  "data": null
}

// 分页响应
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [ ... ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 100,
      "total_pages": 5
    }
  }
}
```

### 5.3 错误码定义

```go
// internal/pkg/errors/errors.go

package errors

// 通用错误码 (0-999)
const (
	CodeSuccess        = 0
	CodeInternalError  = 500
	CodeBadRequest     = 400
	CodeUnauthorized   = 401
	CodeForbidden      = 403
	CodeNotFound       = 404
	CodeTooManyRequests = 429
)

// 用户服务错误码 (1000-1999)
const (
	CodeUserNotFound      = 1001
	CodeUserExists        = 1002
	CodeInvalidPassword   = 1003
	CodeInvalidToken      = 1004
	CodeTokenExpired      = 1005
	CodePhoneExists       = 1006
	CodeInvalidPhone      = 1007
	CodeChildLimitReached = 1008
)

// 训练服务错误码 (2000-2999)
const (
	CodeExerciseNotFound   = 2001
	CodePlanNotFound       = 2002
	CodeInvalidExerciseType = 2003
	CodePlanAlreadyExists  = 2004
)

// 视力服务错误码 (3000-3999)
const (
	CodeVisionRecordNotFound = 3001
	CodeInvalidVisionData    = 3002
	CodeOCRFailed           = 3003
)

// 激励服务错误码 (4000-4999)
const (
	CodeBadgeNotFound      = 4001
	CodeBadgeAlreadyEarned = 4002
	CodeChallengeNotFound  = 4003
	CodeChallengeExpired   = 4004
	CodeInsufficientPoints = 4005
)
```

### 5.4 核心接口列表

#### 5.4.1 用户服务 (user-svc)

| 方法 | 路径 | 说明 | 认证 |
|-----|------|------|------|
| POST | /api/v1/auth/register | 家长注册 | 否 |
| POST | /api/v1/auth/login | 登录 | 否 |
| POST | /api/v1/auth/refresh | 刷新Token | 是 |
| POST | /api/v1/auth/logout | 登出 | 是 |
| GET | /api/v1/users/me | 获取当前用户 | 是 |
| PUT | /api/v1/users/me | 更新用户信息 | 是 |
| POST | /api/v1/children | 添加儿童 | 是 |
| GET | /api/v1/children | 获取儿童列表 | 是 |
| GET | /api/v1/children/:id | 获取儿童详情 | 是 |
| PUT | /api/v1/children/:id | 更新儿童信息 | 是 |
| DELETE | /api/v1/children/:id | 删除儿童 | 是 |
| GET | /api/v1/settings | 获取家长设置 | 是 |
| PUT | /api/v1/settings | 更新家长设置 | 是 |

#### 5.4.2 训练服务 (training-svc)

| 方法 | 路径 | 说明 | 认证 |
|-----|------|------|------|
| GET | /api/v1/plans/today | 获取今日计划 | 是 |
| GET | /api/v1/plans | 获取计划列表 | 是 |
| POST | /api/v1/plans/:id/complete | 完成计划 | 是 |
| POST | /api/v1/exercises | 提交运动记录 | 是 |
| GET | /api/v1/exercises | 获取运动记录 | 是 |
| GET | /api/v1/exercises/personal-best | 获取个人最佳 | 是 |
| POST | /api/v1/assessments | 提交体能评估 | 是 |
| GET | /api/v1/assessments/latest | 获取最新评估 | 是 |
| GET | /api/v1/stats/weekly | 获取周统计 | 是 |
| GET | /api/v1/stats/monthly | 获取月统计 | 是 |

#### 5.4.3 视力服务 (vision-svc)

| 方法 | 路径 | 说明 | 认证 |
|-----|------|------|------|
| POST | /api/v1/vision-records | 添加视力记录 | 是 |
| GET | /api/v1/vision-records | 获取视力记录 | 是 |
| GET | /api/v1/vision-records/:id | 获取视力详情 | 是 |
| POST | /api/v1/vision-records/ocr | OCR识别验光单 | 是 |
| GET | /api/v1/vision/trend | 获取视力趋势 | 是 |
| GET | /api/v1/outdoor/today | 获取今日户外时间 | 是 |
| GET | /api/v1/outdoor/history | 获取户外历史 | 是 |
| POST | /api/v1/outdoor/sync | 同步户外数据 | 是 |
| GET | /api/v1/reminders | 获取提醒记录 | 是 |
| POST | /api/v1/reminders/:id/ack | 确认提醒 | 是 |

#### 5.4.4 激励服务 (reward-svc)

| 方法 | 路径 | 说明 | 认证 |
|-----|------|------|------|
| GET | /api/v1/badges | 获取勋章列表 | 是 |
| GET | /api/v1/badges/my | 获取我的勋章 | 是 |
| GET | /api/v1/points | 获取积分记录 | 是 |
| GET | /api/v1/points/balance | 获取积分余额 | 是 |
| POST | /api/v1/challenges | 发起挑战 | 是 |
| GET | /api/v1/challenges | 获取挑战列表 | 是 |
| POST | /api/v1/challenges/:id/accept | 接受挑战 | 是 |
| POST | /api/v1/challenges/:id/submit | 提交挑战成绩 | 是 |
| GET | /api/v1/leaderboard/family | 家庭排行榜 | 是 |
| GET | /api/v1/leaderboard/global | 全局排行榜 | 是 |

### 5.5 接口详细定义示例

#### 提交运动记录

```
POST /api/v1/exercises
Authorization: Bearer <jwt_token>
Content-Type: application/json

Request:
{
  "type": "jump_rope",
  "duration_seconds": 120,
  "count": 150,
  "score": 85,
  "rhythm_score": 90,
  "amplitude_score": 80,
  "symmetry_score": 85,
  "continuity_score": 88,
  "corrections": ["手臂再收拢一点"],
  "is_offline": false,
  "started_at": "2024-01-15T08:30:00Z",
  "completed_at": "2024-01-15T08:32:00Z"
}

Response (200 OK):
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "type": "jump_rope",
    "duration_seconds": 120,
    "count": 150,
    "score": 85,
    "points_earned": 10,
    "badges_earned": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440001",
        "name": "跳绳达人",
        "icon": "https://cdn.kidsfit.com/badges/jump_rope_master.png"
      }
    ],
    "is_record_broken": false,
    "created_at": "2024-01-15T08:32:05Z"
  }
}
```

---

## 六、AI模块集成方案

### 6.1 架构设计

AI识别主要在客户端完成（ML Kit Pose Detection），服务端负责：
- 接收运动结果数据（评分、计数、纠正建议）
- 训练计划智能推荐算法
- 动作质量趋势分析
- 模型版本管理（OTA更新）

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   客户端 (Flutter) │     │   API Gateway   │     │  training-svc   │
│                 │     │                 │     │   (Go/Gin)      │
│  ML Kit Pose    │────►│  /api/v1/       │────►│  接收运动数据    │
│  Detection      │     │  exercises      │     │  存储记录       │
│  (本地实时)      │     │                 │     │  触发积分计算    │
└─────────────────┘     └─────────────────┘     └─────────────────┘
                                                        │
                                                        ▼
                                               ┌─────────────────┐
                                               │   reward-svc    │
                                               │  计算积分/勋章   │
                                               └─────────────────┘
```

### 6.2 训练计划推荐算法（服务端实现）

```go
// internal/application/training/recommendation_service.go

package training

import (
	"context"
	"math"
	"sort"
	"time"

	"github.com/google/uuid"
)

// RecommendationService 训练推荐服务
type RecommendationService struct {
	trainingRepo   training.Repository
	visionRepo     vision.Repository
	exerciseRepo   exercise.Repository
}

// GenerateDailyPlan 生成每日训练计划
func (s *RecommendationService) GenerateDailyPlan(ctx context.Context, userID uuid.UUID, date time.Time) (*training.TrainingPlan, error) {
	// 1. 获取用户画像
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 2. 获取最新体能评估
	assessment, err := s.trainingRepo.GetLatestFitnessAssessment(ctx, userID)
	if err != nil {
		// 使用默认评估
		assessment = &training.FitnessAssessment{
			Endurance:    5,
			Agility:      5,
			Strength:     5,
			Speed:        5,
			Coordination: 5,
			Balance:      5,
			Flexibility:  5,
		}
	}

	// 3. 获取历史运动数据
	weekAgo := date.AddDate(0, 0, -7)
	records, err := s.trainingRepo.GetExerciseRecordsByUserID(ctx, userID, weekAgo, date)
	if err != nil {
		records = []*training.ExerciseRecord{}
	}

	// 4. 获取视力数据
	outdoorActivity, err := s.visionRepo.GetOutdoorActivityByUserIDAndDate(ctx, userID, date)
	if err != nil {
		outdoorActivity = &vision.OutdoorActivity{DurationMin: 0}
	}

	// 5. 构建用户画像
	profile := UserProfile{
		Age:              user.Age,
		Fitness:          assessment,
		WeeklyRecords:    records,
		DailyOutdoorMin:  outdoorActivity.DurationMin,
		TargetOutdoorMin: 120,
	}

	// 6. 计算训练总量
	baseDuration := ageBasedDuration(user.Age)
	adjustment := calculateDurationAdjustment(profile)
	totalDuration := clamp(baseDuration+adjustment, 15, 30)

	// 7. 构建训练计划
	plan := &training.TrainingPlan{
		UserID:        userID,
		Date:          date,
		Status:        training.PlanStatusPending,
		TotalDuration: totalDuration,
	}

	// 8. 生成热身
	warmupDuration := int(float64(totalDuration) * 0.2)
	plan.WarmupItems = generateWarmupExercises(user.Age, warmupDuration)

	// 9. 生成主训练
	mainDuration := int(float64(totalDuration) * 0.6)
	plan.MainItems = generateMainExercises(profile, mainDuration)

	// 10. 生成拉伸
	cooldownDuration := int(float64(totalDuration) * 0.2)
	plan.CooldownItems = generateCooldownExercises(user.Age, cooldownDuration)

	return plan, nil
}

// ageBasedDuration 根据年龄确定基础训练时长
func ageBasedDuration(age int) int {
	switch {
	case age >= 3 && age <= 6:
		return 15
	case age >= 7 && age <= 9:
		return 20
	case age >= 10 && age <= 12:
		return 25
	default:
		return 20
	}
}

// calculateDurationAdjustment 计算时长调整
func calculateDurationAdjustment(profile UserProfile) int {
	adjustment := 0

	// 根据完成率调整
	completionRate := calculateCompletionRate(profile.WeeklyRecords)
	if completionRate > 0.8 {
		adjustment += 5
	} else if completionRate < 0.5 {
		adjustment -= 5
	}

	// 户外时间不足时增加时长（鼓励户外）
	if profile.DailyOutdoorMin < profile.TargetOutdoorMin {
		adjustment += 5
	}

	return adjustment
}

// generateMainExercises 生成主训练项目
func generateMainExercises(profile UserProfile, duration int) []training.ExerciseItem {
	// 找出弱项
	weaknesses := findWeaknesses(profile.Fitness)

	// 根据弱项和偏好选择运动
	exercises := selectExercises(profile.Age, weaknesses, duration)

	return exercises
}

// findWeaknesses 找出体能弱项
func findWeaknesses(assessment *training.FitnessAssessment) []string {
	dimensions := map[string]int{
		"endurance":    assessment.Endurance,
		"agility":      assessment.Agility,
		"strength":     assessment.Strength,
		"speed":        assessment.Speed,
		"coordination": assessment.Coordination,
		"balance":      assessment.Balance,
		"flexibility":  assessment.Flexibility,
	}

	// 按分数排序，取最低的2项
	type pair struct {
		name  string
		score int
	}
	var pairs []pair
	for name, score := range dimensions {
		pairs = append(pairs, pair{name, score})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].score < pairs[j].score
	})

	var weaknesses []string
	for i := 0; i < min(2, len(pairs)); i++ {
		weaknesses = append(weaknesses, pairs[i].name)
	}
	return weaknesses
}
```

### 6.3 模型管理

```go
// internal/infrastructure/ai/model_manager.go

package ai

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// ModelManager AI模型管理器
type ModelManager struct {
	modelDir   string
	versionURL string
	httpClient *http.Client
}

// ModelInfo 模型信息
type ModelInfo struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	URL        string `json:"url"`
	Checksum   string `json:"checksum"`
	MinAge     int    `json:"min_age"`
	MaxAge     int    `json:"max_age"`
	Required   bool   `json:"required"`
}

// CheckUpdate 检查模型更新
func (m *ModelManager) CheckUpdate(ctx context.Context) ([]ModelInfo, error) {
	resp, err := m.httpClient.Get(m.versionURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var models []ModelInfo
	if err := json.NewDecoder(resp.Body).Decode(&models); err != nil {
		return nil, err
	}

	// 检查本地版本
	var updates []ModelInfo
	for _, model := range models {
		localVersion := m.getLocalVersion(model.Name)
		if localVersion != model.Version {
			updates = append(updates, model)
		}
	}

	return updates, nil
}

// DownloadModel 下载模型
func (m *ModelManager) DownloadModel(ctx context.Context, model ModelInfo) error {
	resp, err := m.httpClient.Get(model.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 创建临时文件
	tmpFile := filepath.Join(m.modelDir, fmt.Sprintf("%s_%s.tmp", model.Name, model.Version))
	f, err := os.Create(tmpFile)
	if err != nil {
		return err
	}
	defer f.Close()

	// 写入文件
	if _, err := io.Copy(f, resp.Body); err != nil {
		return err
	}

	// 校验checksum
	// ...

	// 替换旧版本
	modelPath := filepath.Join(m.modelDir, fmt.Sprintf("%s.tflite", model.Name))
	return os.Rename(tmpFile, modelPath)
}
```

---

## 七、安全与隐私设计

### 7.1 认证与授权

#### 7.1.1 JWT认证

```go
// internal/pkg/jwt/jwt.go

package jwt

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT声明
type Claims struct {
	UserID   string `json:"user_id"`
	UserType string `json:"user_type"`
	ParentID string `json:"parent_id,omitempty"`
	jwt.RegisteredClaims
}

// TokenPair Token对
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// Generator Token生成器
type Generator struct {
	secret        []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// GenerateTokenPair 生成Token对
func (g *Generator) GenerateTokenPair(userID, userType, parentID string) (*TokenPair, error) {
	now := time.Now()

	// Access Token
	accessClaims := Claims{
		UserID:   userID,
		UserType: userType,
		ParentID: parentID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(g.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "kidsfit",
			Subject:   userID,
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(g.secret)
	if err != nil {
		return nil, err
	}

	// Refresh Token
	refreshClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(g.refreshExpiry)),
		IssuedAt:  jwt.NewNumericDate(now),
		Subject:   userID,
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(g.secret)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(g.accessExpiry.Seconds()),
	}, nil
}

// ValidateToken 验证Token
func (g *Generator) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return g.secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}
```

#### 7.1.2 权限控制

```go
// internal/api/http/middleware/auth.go

package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware(jwtGen *jwt.Generator) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, response.Error(errors.CodeUnauthorized, "缺少认证信息"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, response.Error(errors.CodeUnauthorized, "认证格式错误"))
			c.Abort()
			return
		}

		claims, err := jwtGen.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.Error(errors.CodeInvalidToken, "Token无效或已过期"))
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("user_type", claims.UserType)
		c.Set("parent_id", claims.ParentID)

		c.Next()
	}
}

// RequireParent 要求家长权限
func RequireParent() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists || userType != "parent" {
			c.JSON(http.StatusForbidden, response.Error(errors.CodeForbidden, "需要家长权限"))
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireChildOrParent 要求儿童本人或关联家长
func RequireChildOrParent() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, _ := c.Get("user_type")
		userID, _ := c.Get("user_id")
		parentID, _ := c.Get("parent_id")

		// 获取请求中的child_id
		childID := c.Param("child_id")

		// 如果是儿童本人
		if userType == "child" && userID == childID {
			c.Next()
			return
		}

		// 如果是关联家长
		if userType == "parent" {
			// 验证家长是否关联该儿童
			// ...
			c.Next()
			return
		}

		c.JSON(http.StatusForbidden, response.Error(errors.CodeForbidden, "无权访问"))
		c.Abort()
	}
}
```

### 7.2 数据加密

```go
// internal/pkg/crypto/aes.go

package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

// AESGCM AES-GCM加密
type AESGCM struct {
	key []byte
}

// NewAESGCM 创建加密器
func NewAESGCM(key []byte) (*AESGCM, error) {
	if len(key) != 32 {
		return nil, errors.New("key must be 32 bytes")
	}
	return &AESGCM{key: key}, nil
}

// Encrypt 加密
func (a *AESGCM) Encrypt(plaintext []byte) (string, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 解密
func (a *AESGCM) Decrypt(ciphertext string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertextBytes, nil)
}
```

### 7.3 隐私保护措施

| 措施 | 实现方式 |
|-----|---------|
| 数据最小化 | 仅收集年龄（用于分龄内容）、运动数据、视力数据，不收集姓名、学校等 |
| 本地优先 | 运动视频、照片仅本地存储，不上传云端 |
| 实时处理 | 摄像头数据实时分析，服务端不接收原始视频帧 |
| 位置模糊 | GPS坐标模糊到100米精度，仅存储户外时间统计 |
| 数据加密 | 视力数据、体征数据AES-256加密存储 |
| 传输加密 | 全站HTTPS/TLS 1.3 |
| 家长控制 | 所有儿童数据收集需家长知情同意，家长可随时删除全部数据 |
| 审计日志 | 所有数据访问操作记录审计日志 |

---

## 八、性能与扩展性设计

### 8.1 缓存策略

```go
// internal/infrastructure/redis/cache.go

package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache Redis缓存
type Cache struct {
	client *redis.Client
	prefix string
}

// KeyPattern 缓存键模式
type KeyPattern string

const (
	KeyUserProfile     KeyPattern = "user:profile:%s"
	KeyUserChildren    KeyPattern = "user:children:%s"
	KeyTodayPlan       KeyPattern = "plan:today:%s"
	KeyPersonalBest    KeyPattern = "pb:%s:%s"       // pb:user_id:exercise_type
	KeyLeaderboard     KeyPattern = "lb:%s:%s"       // lb:type:period
	KeyVisionTrend     KeyPattern = "vision:trend:%s"
	KeyOutdoorToday    KeyPattern = "outdoor:today:%s"
	KeyRateLimit       KeyPattern = "ratelimit:%s:%s" // ratelimit:user_id:action
)

// GetUserProfile 获取用户资料（带缓存）
func (c *Cache) GetUserProfile(ctx context.Context, userID string) (*user.User, error) {
	key := fmt.Sprintf(string(KeyUserProfile), userID)

	// 尝试从缓存获取
	data, err := c.client.Get(ctx, key).Result()
	if err == nil {
		var u user.User
		if err := json.Unmarshal([]byte(data), &u); err == nil {
			return &u, nil
		}
	}

	// 缓存未命中，从数据库获取
	// ...

	return nil, err
}

// SetUserProfile 缓存用户资料
func (c *Cache) SetUserProfile(ctx context.Context, userID string, u *user.User, ttl time.Duration) error {
	key := fmt.Sprintf(string(KeyUserProfile), userID)
	data, err := json.Marshal(u)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}

// UpdateLeaderboard 更新排行榜
func (c *Cache) UpdateLeaderboard(ctx context.Context, boardType string, period string, userID string, score int) error {
	key := fmt.Sprintf(string(KeyLeaderboard), boardType, period)
	return c.client.ZAdd(ctx, key, redis.Z{
		Score:  float64(score),
		Member: userID,
	}).Err()
}

// GetLeaderboard 获取排行榜
func (c *Cache) GetLeaderboard(ctx context.Context, boardType string, period string, topN int64) ([]redis.Z, error) {
	key := fmt.Sprintf(string(KeyLeaderboard), boardType, period)
	return c.client.ZRevRangeWithScores(ctx, key, 0, topN-1).Result()
}
```

### 8.2 限流策略

```go
// internal/api/http/middleware/rate_limit.go

package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimiter 限流器
type RateLimiter struct {
	redisClient *redis.Client
	window      time.Duration
	maxRequests int
}

// NewRateLimiter 创建限流器
func NewRateLimiter(redisClient *redis.Client, window time.Duration, maxRequests int) *RateLimiter {
	return &RateLimiter{
		redisClient: redisClient,
		window:      window,
		maxRequests: maxRequests,
	}
}

// RateLimit 限流中间件
func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.Next()
			return
		}

		key := fmt.Sprintf("ratelimit:%s:%s", userID, c.Request.URL.Path)
		now := time.Now().Unix()
		windowStart := now - int64(rl.window.Seconds())

		// 使用Redis有序集合实现滑动窗口限流
		pipe := rl.redisClient.Pipeline()
		pipe.ZRemRangeByScore(c.Request.Context(), key, "0", fmt.Sprintf("%d", windowStart))
		pipe.ZCard(c.Request.Context(), key)
		pipe.ZAdd(c.Request.Context(), key, redis.Z{Score: float64(now), Member: now})
		pipe.Expire(c.Request.Context(), key, rl.window)

		results, err := pipe.Exec(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(errors.CodeInternalError, "限流服务异常"))
			c.Abort()
			return
		}

		count := results[1].(*redis.IntCmd).Val()
		if int(count) >= rl.maxRequests {
			c.JSON(http.StatusTooManyRequests, response.Error(errors.CodeTooManyRequests, "请求过于频繁，请稍后再试"))
			c.Abort()
			return
		}

		c.Next()
	}
}
```

### 8.3 数据库优化

| 优化措施 | 实现 |
|---------|------|
| 读写分离 | 主库写、从库读，使用pgpool-II或应用层实现 |
| 连接池 | 最大连接数100，最大空闲连接10，连接超时30s |
| 索引优化 | 所有外键、查询条件字段建立索引，定期分析执行计划 |
| 分表策略 | 运动记录按用户ID哈希分表（未来数据量大时） |
| 归档策略 | 超过1年的历史数据归档到冷存储 |

---

## 九、部署与运维架构

### 9.1 容器化部署

```dockerfile
# deployments/docker/Dockerfile

# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装依赖
RUN apk add --no-cache git

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源码
COPY . .

# 构建二进制
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/user-svc ./cmd/user-svc
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/training-svc ./cmd/training-svc
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/vision-svc ./cmd/vision-svc
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/reward-svc ./cmd/reward-svc
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/gateway ./cmd/gateway

# 运行阶段
FROM alpine:3.18

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# 从构建阶段复制二进制
COPY --from=builder /app/bin/* /app/
COPY configs/ /app/configs/

# 非root用户运行
RUN adduser -D -u 1000 appuser
USER appuser

EXPOSE 8080

ENTRYPOINT ["/app/gateway"]
```

### 9.2 Kubernetes部署

```yaml
# deployments/k8s/deployment.yaml

apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-svc
  namespace: kidsfit
  labels:
    app: user-svc
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  selector:
    matchLabels:
      app: user-svc
  template:
    metadata:
      labels:
        app: user-svc
    spec:
      containers:
      - name: user-svc
        image: registry.cn-hangzhou.aliyuncs.com/kidsfit/user-svc:v1.0.0
        ports:
        - containerPort: 8001
          name: http
        env:
        - name: ENV
          value: "production"
        - name: DB_HOST
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: host
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8001
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8001
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: user-svc
  namespace: kidsfit
spec:
  selector:
    app: user-svc
  ports:
  - port: 8001
    targetPort: 8001
    name: http
  type: ClusterIP
```

### 9.3 监控告警

```yaml
# Prometheus监控配置

# 服务指标
- name: http_requests_total
  type: counter
  labels: [method, path, status]
  help: HTTP请求总数

- name: http_request_duration_seconds
  type: histogram
  labels: [method, path]
  buckets: [0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10]
  help: HTTP请求耗时

- name: active_users
  type: gauge
  labels: [user_type]
  help: 当前活跃用户

- name: exercise_records_total
  type: counter
  labels: [exercise_type, age_group]
  help: 运动记录总数

# 告警规则
alerting:
  rules:
  - alert: HighErrorRate
    expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: "错误率过高"
      description: "服务 {{ $labels.service }} 错误率超过5%"

  - alert: HighLatency
    expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "响应延迟过高"
      description: "P95延迟超过1秒"

  - alert: DatabaseConnectionHigh
    expr: pg_stat_activity_count > 80
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "数据库连接数过高"
```

### 9.4 CI/CD流水线

```yaml
# .github/workflows/ci-cd.yml

name: CI/CD

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    - run: go mod download
    - run: go test -v -race -coverprofile=coverage.out ./...
    - run: go tool cover -html=coverage.out -o coverage.html

  lint:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - uses: golangci/golangci-lint-action@v3
      with:
        version: latest

  build:
    needs: [test, lint]
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - name: Build images
      run: |
        docker build -t user-svc:${{ github.sha }} -f deployments/docker/Dockerfile --target user-svc .
        docker build -t training-svc:${{ github.sha }} -f deployments/docker/Dockerfile --target training-svc .
    - name: Push to registry
      run: |
        docker tag user-svc:${{ github.sha }} registry.cn-hangzhou.aliyuncs.com/kidsfit/user-svc:${{ github.sha }}
        docker push registry.cn-hangzhou.aliyuncs.com/kidsfit/user-svc:${{ github.sha }}

  deploy:
    needs: build
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - name: Deploy to K8s
      run: |
        kubectl set image deployment/user-svc user-svc=registry.cn-hangzhou.aliyuncs.com/kidsfit/user-svc:${{ github.sha }} -n kidsfit
        kubectl rollout status deployment/user-svc -n kidsfit
```

---

## 十、开发规范

### 10.1 代码规范

- 遵循 **Go Code Review Comments** 和 **Effective Go**
- 使用 **golangci-lint** 进行静态检查
- 单元测试覆盖率 ≥ 70%
- 核心业务流程覆盖率 ≥ 90%

### 10.2 日志规范

```go
// 使用结构化日志
logger.Info("用户登录",
    zap.String("user_id", userID),
    zap.String("ip", clientIP),
    zap.Duration("latency