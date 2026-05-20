package reward

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// BadgeCategory 徽章类别枚举
type BadgeCategory string

const (
	// BadgeCategoryMilestone 里程碑徽章
	BadgeCategoryMilestone BadgeCategory = "milestone"
	// BadgeCategorySkill 技能徽章
	BadgeCategorySkill BadgeCategory = "skill"
	// BadgeCategoryStreak 连续打卡徽章
	BadgeCategoryStreak BadgeCategory = "streak"
	// BadgeCategoryChallenge 挑战徽章
	BadgeCategoryChallenge BadgeCategory = "challenge"
	// BadgeCategoryFamily 家庭徽章
	BadgeCategoryFamily BadgeCategory = "family"
	// BadgeCategoryVision 护眼徽章
	BadgeCategoryVision BadgeCategory = "vision"
	// BadgeCategorySpecial 特殊徽章
	BadgeCategorySpecial BadgeCategory = "special"
)

// PointType 积分类型枚举
type PointType string

const (
	// PointTypeExercise 运动获得积分
	PointTypeExercise PointType = "exercise"
	// PointTypeRecordBreak 破纪录积分
	PointTypeRecordBreak PointType = "record_break"
	// PointTypeFamilyActivity 家庭活动积分
	PointTypeFamilyActivity PointType = "family_activity"
	// PointTypeStreak 连续打卡积分
	PointTypeStreak PointType = "streak"
	// PointTypeVisionTask 护眼任务积分
	PointTypeVisionTask PointType = "vision_task"
	// PointTypeRedeem 积分兑换
	PointTypeRedeem PointType = "redeem"
)

// ChallengeType 挑战类型枚举
type ChallengeType string

const (
	// ChallengeTypeSync 同步挑战
	ChallengeTypeSync ChallengeType = "sync"
	// ChallengeTypeAsync 异步挑战
	ChallengeTypeAsync ChallengeType = "async"
	// ChallengeTypeTimed 限时挑战
	ChallengeTypeTimed ChallengeType = "timed"
)

// ChallengeStatus 挑战状态枚举
type ChallengeStatus string

const (
	// ChallengeStatusPending 待接受
	ChallengeStatusPending ChallengeStatus = "pending"
	// ChallengeStatusAccepted 已接受
	ChallengeStatusAccepted ChallengeStatus = "accepted"
	// ChallengeStatusCompleted 已完成
	ChallengeStatusCompleted ChallengeStatus = "completed"
	// ChallengeStatusExpired 已过期
	ChallengeStatusExpired ChallengeStatus = "expired"
)

// Badge 徽章领域模型
type Badge struct {
	ID          string          `json:"id" db:"id"`
	Code        string          `json:"code" db:"code"`
	Name        string          `json:"name" db:"name"`
	Description string          `json:"description" db:"description"`
	Category    BadgeCategory   `json:"category" db:"category"`
	Icon        string          `json:"icon" db:"icon"`
	Condition   json.RawMessage `json:"condition" db:"condition"`
	Points      int             `json:"points" db:"points"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
}

// NewBadge 创建徽章实例
func NewBadge(code, name, description string, category BadgeCategory, icon string, condition json.RawMessage, points int) *Badge {
	return &Badge{
		ID:          uuid.New().String(),
		Code:        code,
		Name:        name,
		Description: description,
		Category:    category,
		Icon:        icon,
		Condition:   condition,
		Points:      points,
		CreatedAt:   time.Now(),
	}
}

// UserBadge 用户徽章领域模型
type UserBadge struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	BadgeID   string    `json:"badge_id" db:"badge_id"`
	EarnedAt  time.Time `json:"earned_at" db:"earned_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// NewUserBadge 创建用户徽章实例
func NewUserBadge(userID, badgeID string) *UserBadge {
	now := time.Now()
	return &UserBadge{
		ID:        uuid.New().String(),
		UserID:    userID,
		BadgeID:   badgeID,
		EarnedAt:  now,
		CreatedAt: now,
	}
}

// PointRecord 积分记录领域模型
type PointRecord struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"user_id" db:"user_id"`
	Points      int       `json:"points" db:"points"`
	Type        PointType `json:"type" db:"type"`
	SourceID    string    `json:"source_id" db:"source_id"`
	SourceType  string    `json:"source_type" db:"source_type"`
	Description string    `json:"description" db:"description"`
	Balance     int       `json:"balance" db:"balance"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// NewPointRecord 创建积分记录实例
func NewPointRecord(userID string, points int, pointType PointType, sourceID, sourceType, description string) *PointRecord {
	return &PointRecord{
		ID:          uuid.New().String(),
		UserID:      userID,
		Points:      points,
		Type:        pointType,
		SourceID:    sourceID,
		SourceType:  sourceType,
		Description: description,
		CreatedAt:   time.Now(),
	}
}

// Challenge 挑战领域模型
type Challenge struct {
	ID             string          `json:"id" db:"id"`
	Type           ChallengeType   `json:"type" db:"type"`
	InitiatorID    string          `json:"initiator_id" db:"initiator_id"`
	AcceptorID     *string         `json:"acceptor_id,omitempty" db:"acceptor_id"`
	ExerciseType   string          `json:"exercise_type" db:"exercise_type"`
	TargetValue    int             `json:"target_value" db:"target_value"`
	InitiatorScore *int            `json:"initiator_score,omitempty" db:"initiator_score"`
	AcceptorScore  *int            `json:"acceptor_score,omitempty" db:"acceptor_score"`
	WinnerID       *string         `json:"winner_id,omitempty" db:"winner_id"`
	Status         ChallengeStatus `json:"status" db:"status"`
	ExpiresAt      *time.Time      `json:"expires_at,omitempty" db:"expires_at"`
	CompletedAt    *time.Time      `json:"completed_at,omitempty" db:"completed_at"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
}

// NewChallenge 创建挑战实例
func NewChallenge(challengeType ChallengeType, initiatorID, exerciseType string, targetValue int) *Challenge {
	return &Challenge{
		ID:           uuid.New().String(),
		Type:         challengeType,
		InitiatorID:  initiatorID,
		ExerciseType: exerciseType,
		TargetValue:  targetValue,
		Status:       ChallengeStatusPending,
		CreatedAt:    time.Now(),
	}
}

// IsExpired 判断挑战是否已过期
func (c *Challenge) IsExpired() bool {
	if c.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*c.ExpiresAt)
}

// CanAccept 判断挑战是否可以被接受
func (c *Challenge) CanAccept() bool {
	return c.Status == ChallengeStatusPending && !c.IsExpired()
}

// CanSubmit 判断挑战是否可以提交成绩
func (c *Challenge) CanSubmit() bool {
	return c.Status == ChallengeStatusAccepted && !c.IsExpired()
}
