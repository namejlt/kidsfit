package training

import (
	"context"
	"time"
)

// ExerciseRecordFilter 运动记录查询过滤条件
type ExerciseRecordFilter struct {
	Type     ExerciseType `json:"type,omitempty"`
	FromDate *time.Time   `json:"from_date,omitempty"`
	ToDate   *time.Time   `json:"to_date,omitempty"`
}

// Pagination 分页参数
type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// PaginatedResult 分页查询结果
type PaginatedResult[T any] struct {
	Items      []T   `json:"items"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

// ExerciseRecordRepository 运动记录仓储接口
type ExerciseRecordRepository interface {
	// Create 创建运动记录
	Create(ctx context.Context, record *ExerciseRecord) error

	// GetByID 根据ID获取运动记录
	GetByID(ctx context.Context, id string) (*ExerciseRecord, error)

	// GetByUserID 根据用户ID分页查询运动记录
	GetByUserID(ctx context.Context, userID string, filter ExerciseRecordFilter, pagination Pagination) (*PaginatedResult[ExerciseRecord], error)

	// GetByUserIDAndType 根据用户ID和运动类型查询运动记录
	GetByUserIDAndType(ctx context.Context, userID string, exerciseType ExerciseType, pagination Pagination) (*PaginatedResult[ExerciseRecord], error)

	// GetPersonalBest 获取用户某类运动的个人最佳记录
	GetPersonalBest(ctx context.Context, userID string, exerciseType ExerciseType) (*ExerciseRecord, error)

	// CreateBatch 批量创建运动记录
	CreateBatch(ctx context.Context, records []*ExerciseRecord) error
}

// TrainingPlanRepository 训练计划仓储接口
type TrainingPlanRepository interface {
	// Create 创建训练计划
	Create(ctx context.Context, plan *TrainingPlan) error

	// GetByID 根据ID获取训练计划
	GetByID(ctx context.Context, id string) (*TrainingPlan, error)

	// GetByUserIDAndDate 根据用户ID和日期获取训练计划
	GetByUserIDAndDate(ctx context.Context, userID string, date time.Time) (*TrainingPlan, error)

	// GetByUserID 根据用户ID分页查询训练计划
	GetByUserID(ctx context.Context, userID string, pagination Pagination) (*PaginatedResult[TrainingPlan], error)

	// Update 更新训练计划
	Update(ctx context.Context, plan *TrainingPlan) error
}

// FitnessAssessmentRepository 体能评估仓储接口
type FitnessAssessmentRepository interface {
	// Create 创建体能评估
	Create(ctx context.Context, assessment *FitnessAssessment) error

	// GetLatestByUserID 获取用户最新的体能评估
	GetLatestByUserID(ctx context.Context, userID string) (*FitnessAssessment, error)
}
