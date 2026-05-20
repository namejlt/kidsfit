package user

import "time"

// RegisterRequest 注册请求DTO，用于家长账号注册
type RegisterRequest struct {
	// Phone 手机号
	Phone string `json:"phone" binding:"required"`
	// Password 密码
	Password string `json:"password" binding:"required,min=6,max=20"`
	// Nickname 昵称
	Nickname string `json:"nickname" binding:"required"`
}

// LoginRequest 登录请求DTO
type LoginRequest struct {
	// Phone 手机号
	Phone string `json:"phone" binding:"required"`
	// Password 密码
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应DTO，包含令牌和用户信息
type LoginResponse struct {
	// AccessToken 访问令牌
	AccessToken string `json:"access_token"`
	// RefreshToken 刷新令牌
	RefreshToken string `json:"refresh_token"`
	// ExpiresIn 令牌过期时间（秒）
	ExpiresIn int64 `json:"expires_in"`
	// User 用户信息
	User *UserDTO `json:"user"`
}

// UserDTO 用户信息DTO
type UserDTO struct {
	// ID 用户ID
	ID string `json:"id"`
	// Type 用户类型（parent/child）
	Type string `json:"type"`
	// Nickname 昵称
	Nickname string `json:"nickname"`
	// Avatar 头像URL
	Avatar string `json:"avatar"`
	// Phone 手机号（仅家长）
	Phone string `json:"phone,omitempty"`
	// Age 年龄（仅儿童）
	Age int `json:"age,omitempty"`
	// Status 用户状态
	Status string `json:"status"`
	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`
}

// ChildDTO 儿童信息DTO
type ChildDTO struct {
	// ID 儿童用户ID
	ID string `json:"id"`
	// Nickname 昵称
	Nickname string `json:"nickname"`
	// Avatar 头像URL
	Avatar string `json:"avatar"`
	// Age 年龄
	Age int `json:"age"`
	// AgeGroup 年龄分组（3-6/7-9/10-12）
	AgeGroup string `json:"age_group"`
	// Status 用户状态
	Status string `json:"status"`
}

// ParentSettingsDTO 家长设置DTO
type ParentSettingsDTO struct {
	// DailyLimitMin 每日运动时长限制（分钟）
	DailyLimitMin int `json:"daily_limit_min"`
	// AvailableFrom 允许运动开始时间
	AvailableFrom string `json:"available_from"`
	// AvailableTo 允许运动结束时间
	AvailableTo string `json:"available_to"`
	// CameraAllowed 是否允许使用摄像头
	CameraAllowed bool `json:"camera_allowed"`
	// LocationAllowed 是否允许获取位置
	LocationAllowed bool `json:"location_allowed"`
	// DataUploadCloud 是否允许数据上传云端
	DataUploadCloud bool `json:"data_upload_cloud"`
}

// UpdateUserRequest 更新用户信息请求DTO
type UpdateUserRequest struct {
	// Nickname 昵称
	Nickname string `json:"nickname,omitempty"`
	// Avatar 头像URL
	Avatar string `json:"avatar,omitempty"`
}

// FamilyDTO 家庭关系DTO
type FamilyDTO struct {
	// ID 关系ID
	ID string `json:"id"`
	// ParentID 家长ID
	ParentID string `json:"parent_id"`
	// ChildID 儿童ID
	ChildID string `json:"child_id"`
	// Relation 关系类型
	Relation string `json:"relation"`
	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`
}

// AddChildRequest 添加儿童请求DTO
type AddChildRequest struct {
	// Nickname 儿童昵称
	Nickname string `json:"nickname" binding:"required"`
	// Age 儿童年龄
	Age int `json:"age" binding:"required,min=3,max=12"`
	// Avatar 儿童头像URL
	Avatar string `json:"avatar,omitempty"`
}
