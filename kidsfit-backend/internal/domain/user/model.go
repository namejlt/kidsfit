package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// UserType 用户类型枚举
type UserType string

const (
	// UserTypeParent 家长类型
	UserTypeParent UserType = "parent"
	// UserTypeChild 儿童类型
	UserTypeChild UserType = "child"
)

// UserStatus 用户状态枚举
type UserStatus string

const (
	// UserStatusActive 活跃状态
	UserStatusActive UserStatus = "active"
	// UserStatusInactive 未激活状态
	UserStatusInactive UserStatus = "inactive"
	// UserStatusDeleted 已删除状态
	UserStatusDeleted UserStatus = "deleted"
)

// Relation 家庭关系类型枚举
type Relation string

const (
	// RelationFather 父亲
	RelationFather Relation = "father"
	// RelationMother 母亲
	RelationMother Relation = "mother"
	// RelationGrandfather 祖父/外祖父
	RelationGrandfather Relation = "grandfather"
	// RelationGrandmother 祖母/外祖母
	RelationGrandmother Relation = "grandmother"
	// RelationOther 其他关系
	RelationOther Relation = "other"
)

// User 用户领域模型
type User struct {
	ID           string      `json:"id" db:"id"`
	Type         UserType    `json:"type" db:"type"`
	ParentID     *string     `json:"parent_id,omitempty" db:"parent_id"`
	Age          *int        `json:"age,omitempty" db:"age"`
	Nickname     string      `json:"nickname" db:"nickname"`
	Avatar       string      `json:"avatar" db:"avatar"`
	Phone        *string     `json:"phone,omitempty" db:"phone"`
	PasswordHash *string     `json:"-" db:"password_hash"`
	Status       UserStatus  `json:"status" db:"status"`
	CreatedAt    time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time  `json:"deleted_at,omitempty" db:"deleted_at"`
}

// NewUser 创建新用户实例，自动生成UUID并设置默认状态
func NewUser(userType UserType, nickname string) *User {
	now := time.Now()
	return &User{
		ID:        uuid.New().String(),
		Type:      userType,
		Nickname:  nickname,
		Status:    UserStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// IsParent 判断用户是否为家长类型
func (u *User) IsParent() bool {
	return u.Type == UserTypeParent
}

// IsChild 判断用户是否为儿童类型
func (u *User) IsChild() bool {
	return u.Type == UserTypeChild
}

// IsActive 判断用户是否处于活跃状态
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// ValidateAge 校验儿童年龄是否在3-12岁范围内
func (u *User) ValidateAge() error {
	if u.Age == nil {
		return nil
	}
	if *u.Age < 3 || *u.Age > 12 {
		return errors.New("年龄必须在3-12岁之间")
	}
	return nil
}

// GetAgeGroup 根据年龄返回年龄分组：3-6/7-9/10-12
func (u *User) GetAgeGroup() string {
	if u.Age == nil {
		return ""
	}
	age := *u.Age
	switch {
	case age >= 3 && age <= 6:
		return "3-6"
	case age >= 7 && age <= 9:
		return "7-9"
	case age >= 10 && age <= 12:
		return "10-12"
	default:
		return ""
	}
}

// Family 家庭关系领域模型
type Family struct {
	ID        string   `json:"id" db:"id"`
	ParentID  string   `json:"parent_id" db:"parent_id"`
	ChildID   string   `json:"child_id" db:"child_id"`
	Relation  Relation `json:"relation" db:"relation"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// NewFamily 创建新的家庭关系实例
func NewFamily(parentID, childID string, relation Relation) *Family {
	return &Family{
		ID:        uuid.New().String(),
		ParentID:  parentID,
		ChildID:   childID,
		Relation:  relation,
		CreatedAt: time.Now(),
	}
}

// ParentSettings 家长设置领域模型
type ParentSettings struct {
	ID              string    `json:"id" db:"id"`
	ParentID        string    `json:"parent_id" db:"parent_id"`
	DailyLimitMin   int       `json:"daily_limit_min" db:"daily_limit_min"`
	AvailableFrom   string    `json:"available_from" db:"available_from"`
	AvailableTo     string    `json:"available_to" db:"available_to"`
	CameraAllowed   bool      `json:"camera_allowed" db:"camera_allowed"`
	LocationAllowed bool      `json:"location_allowed" db:"location_allowed"`
	DataUploadCloud bool      `json:"data_upload_cloud" db:"data_upload_cloud"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// NewParentSettings 创建家长设置实例，DailyLimitMin默认30分钟
func NewParentSettings(parentID string) *ParentSettings {
	now := time.Now()
	return &ParentSettings{
		ID:            uuid.New().String(),
		ParentID:      parentID,
		DailyLimitMin: 30,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// ValidateDailyLimit 校验每日时长限制是否在5-120分钟范围内
func (ps *ParentSettings) ValidateDailyLimit() error {
	if ps.DailyLimitMin < 5 || ps.DailyLimitMin > 120 {
		return errors.New("每日时长限制必须在5-120分钟之间")
	}
	return nil
}
