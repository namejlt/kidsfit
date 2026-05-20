package user

import "context"

// UserFilter 用户列表查询过滤条件
type UserFilter struct {
	Type   UserType  `json:"type,omitempty"`
	Status UserStatus `json:"status,omitempty"`
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

// UserRepository 用户仓储接口，负责用户聚合根的持久化操作
type UserRepository interface {
	// Create 创建新用户
	Create(ctx context.Context, user *User) error

	// GetByID 根据用户ID获取用户
	GetByID(ctx context.Context, id string) (*User, error)

	// GetByPhone 根据手机号获取用户
	GetByPhone(ctx context.Context, phone string) (*User, error)

	// Update 更新用户信息
	Update(ctx context.Context, user *User) error

	// Delete 软删除用户
	Delete(ctx context.Context, id string) error

	// List 根据过滤条件和分页参数查询用户列表
	List(ctx context.Context, filter UserFilter, pagination Pagination) (*PaginatedResult[User], error)
}

// FamilyRepository 家庭关系仓储接口
type FamilyRepository interface {
	// Create 创建家庭关系
	Create(ctx context.Context, family *Family) error

	// GetByParentID 根据家长ID获取所有家庭关系
	GetByParentID(ctx context.Context, parentID string) ([]*Family, error)

	// GetByChildID 根据儿童ID获取家庭关系
	GetByChildID(ctx context.Context, childID string) ([]*Family, error)

	// Delete 删除家庭关系
	Delete(ctx context.Context, id string) error
}

// ParentSettingsRepository 家长设置仓储接口
type ParentSettingsRepository interface {
	// Create 创建家长设置
	Create(ctx context.Context, settings *ParentSettings) error

	// GetByParentID 根据家长ID获取设置
	GetByParentID(ctx context.Context, parentID string) (*ParentSettings, error)

	// Update 更新家长设置
	Update(ctx context.Context, settings *ParentSettings) error
}
