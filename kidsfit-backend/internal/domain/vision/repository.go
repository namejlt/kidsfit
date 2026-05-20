package vision

import (
	"context"
	"time"
)

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

// VisionRecordRepository 视力记录仓储接口
type VisionRecordRepository interface {
	// Create 创建视力记录
	Create(ctx context.Context, record *VisionRecord) error

	// GetByID 根据ID获取视力记录
	GetByID(ctx context.Context, id string) (*VisionRecord, error)

	// GetByChildID 根据儿童ID分页查询视力记录
	GetByChildID(ctx context.Context, childID string, pagination Pagination) (*PaginatedResult[VisionRecord], error)

	// GetByChildIDAndDateRange 根据儿童ID和日期范围查询视力记录
	GetByChildIDAndDateRange(ctx context.Context, childID string, startDate, endDate time.Time) ([]*VisionRecord, error)
}

// OutdoorActivityRepository 户外活动仓储接口
type OutdoorActivityRepository interface {
	// Create 创建户外活动记录
	Create(ctx context.Context, activity *OutdoorActivity) error

	// GetByUserIDAndDate 根据用户ID和日期获取户外活动
	GetByUserIDAndDate(ctx context.Context, userID string, date time.Time) (*OutdoorActivity, error)

	// GetByUserIDAndDateRange 根据用户ID和日期范围查询户外活动
	GetByUserIDAndDateRange(ctx context.Context, userID string, startDate, endDate time.Time) ([]*OutdoorActivity, error)

	// Update 更新户外活动记录
	Update(ctx context.Context, activity *OutdoorActivity) error
}

// EyeReminderRepository 护眼提醒仓储接口
type EyeReminderRepository interface {
	// Create 创建护眼提醒
	Create(ctx context.Context, reminder *EyeReminder) error

	// GetByUserID 根据用户ID分页查询护眼提醒
	GetByUserID(ctx context.Context, userID string, pagination Pagination) (*PaginatedResult[EyeReminder], error)

	// UpdateAcknowledged 更新提醒确认状态
	UpdateAcknowledged(ctx context.Context, id string, acknowledged bool) error
}
