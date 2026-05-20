package reward

import "time"

// BadgeDTO 徽章响应DTO
type BadgeDTO struct {
	// ID 徽章ID
	ID string `json:"id"`
	// Code 徽章唯一编码
	Code string `json:"code"`
	// Name 徽章名称
	Name string `json:"name"`
	// Description 徽章描述
	Description string `json:"description"`
	// Category 徽章类别
	Category string `json:"category"`
	// Icon 徽章图标
	Icon string `json:"icon"`
	// Points 徽章积分值
	Points int `json:"points"`
	// Earned 是否已获得
	Earned bool `json:"earned"`
}

// PointRecordDTO 积分记录响应DTO
type PointRecordDTO struct {
	// ID 记录ID
	ID string `json:"id"`
	// UserID 用户ID
	UserID string `json:"user_id"`
	// Points 积分变动值（正为获得，负为消耗）
	Points int `json:"points"`
	// Type 积分类型
	Type string `json:"type"`
	// SourceID 来源ID
	SourceID *string `json:"source_id,omitempty"`
	// SourceType 来源类型
	SourceType *string `json:"source_type,omitempty"`
	// Description 描述
	Description string `json:"description"`
	// Balance 变动后余额
	Balance int `json:"balance"`
	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`
}

// ChallengeDTO 挑战响应DTO
type ChallengeDTO struct {
	// ID 挑战ID
	ID string `json:"id"`
	// Type 挑战类型
	Type string `json:"type"`
	// InitiatorID 发起者ID
	InitiatorID string `json:"initiator_id"`
	// AcceptorID 接受者ID
	AcceptorID *string `json:"acceptor_id,omitempty"`
	// ExerciseType 运动类型
	ExerciseType string `json:"exercise_type"`
	// TargetValue 目标值
	TargetValue int `json:"target_value"`
	// InitiatorScore 发起者成绩
	InitiatorScore *int `json:"initiator_score,omitempty"`
	// AcceptorScore 接受者成绩
	AcceptorScore *int `json:"acceptor_score,omitempty"`
	// WinnerID 获胜者ID
	WinnerID *string `json:"winner_id,omitempty"`
	// Status 挑战状态
	Status string `json:"status"`
	// ExpiresAt 过期时间
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	// CompletedAt 完成时间
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`
}

// CreateChallengeRequest 创建挑战请求DTO
type CreateChallengeRequest struct {
	// Type 挑战类型（sync/async/timed）
	Type string `json:"type" binding:"required"`
	// AcceptorID 接受者ID（可选，为空则公开挑战）
	AcceptorID string `json:"acceptor_id,omitempty"`
	// ExerciseType 运动类型
	ExerciseType string `json:"exercise_type" binding:"required"`
	// TargetValue 目标值
	TargetValue int `json:"target_value" binding:"required,min=1"`
}

// UserBadgeDTO 用户徽章响应DTO
type UserBadgeDTO struct {
	// ID 记录ID
	ID string `json:"id"`
	// UserID 用户ID
	UserID string `json:"user_id"`
	// BadgeID 徽章ID
	BadgeID string `json:"badge_id"`
	// BadgeCode 徽章编码
	BadgeCode string `json:"badge_code"`
	// BadgeName 徽章名称
	BadgeName string `json:"badge_name"`
	// BadgeIcon 徽章图标
	BadgeIcon string `json:"badge_icon"`
	// EarnedAt 获得时间
	EarnedAt time.Time `json:"earned_at"`
}

// LeaderboardDTO 排行榜响应DTO
type LeaderboardDTO struct {
	// Rank 排名
	Rank int64 `json:"rank"`
	// UserID 用户ID
	UserID string `json:"user_id"`
	// Nickname 昵称
	Nickname string `json:"nickname"`
	// Avatar 头像URL
	Avatar string `json:"avatar"`
	// Score 分数
	Score float64 `json:"score"`
}
