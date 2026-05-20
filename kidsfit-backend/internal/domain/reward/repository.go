package reward

import "context"

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

// BadgeRepository 徽章仓储接口
type BadgeRepository interface {
	// GetByID 根据ID获取徽章
	GetByID(ctx context.Context, id string) (*Badge, error)

	// GetByCode 根据唯一编码获取徽章
	GetByCode(ctx context.Context, code string) (*Badge, error)

	// List 按类别过滤查询徽章列表
	List(ctx context.Context, category *BadgeCategory) ([]*Badge, error)
}

// UserBadgeRepository 用户徽章仓储接口
type UserBadgeRepository interface {
	// Create 创建用户徽章记录
	Create(ctx context.Context, userBadge *UserBadge) error

	// GetByUserID 根据用户ID查询已获得的徽章列表
	GetByUserID(ctx context.Context, userID string) ([]*UserBadge, error)

	// GetByUserIDAndBadgeID 根据用户ID和徽章ID查询用户徽章记录
	GetByUserIDAndBadgeID(ctx context.Context, userID, badgeID string) (*UserBadge, error)

	// CountByUserID 统计用户已获得的徽章数量
	CountByUserID(ctx context.Context, userID string) (int64, error)
}

// PointRecordRepository 积分记录仓储接口
type PointRecordRepository interface {
	// Create 创建积分记录
	Create(ctx context.Context, record *PointRecord) error

	// GetByUserID 根据用户ID分页查询积分记录
	GetByUserID(ctx context.Context, userID string, pagination Pagination) (*PaginatedResult[PointRecord], error)

	// GetBalanceByUserID 获取用户当前积分余额
	GetBalanceByUserID(ctx context.Context, userID string) (int, error)

	// UpdateBalance 更新用户积分余额
	UpdateBalance(ctx context.Context, userID string, balance int) error
}

// ChallengeRepository 挑战仓储接口
type ChallengeRepository interface {
	// Create 创建挑战
	Create(ctx context.Context, challenge *Challenge) error

	// GetByID 根据ID获取挑战
	GetByID(ctx context.Context, id string) (*Challenge, error)

	// GetByInitiatorID 根据发起者ID查询挑战列表
	GetByInitiatorID(ctx context.Context, initiatorID string) ([]*Challenge, error)

	// GetByAcceptorID 根据接受者ID查询挑战列表
	GetByAcceptorID(ctx context.Context, acceptorID string) ([]*Challenge, error)

	// Update 更新挑战信息
	Update(ctx context.Context, challenge *Challenge) error
}
