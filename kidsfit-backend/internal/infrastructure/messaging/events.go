package events

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// EventType 事件类型别名，用于标识不同业务事件
type EventType string

// 事件类型常量定义
const (
	// EventUserRegistered 用户注册事件
	EventUserRegistered EventType = "user.registered"
	// EventUserLoggedIn 用户登录事件
	EventUserLoggedIn EventType = "user.logged_in"
	// EventExerciseCompleted 运动完成事件
	EventExerciseCompleted EventType = "exercise.completed"
	// EventExerciseRecordBroken 运动记录打破事件
	EventExerciseRecordBroken EventType = "exercise.record_broken"
	// EventVisionOutdoorTargetMet 户外活动目标达成事件
	EventVisionOutdoorTargetMet EventType = "vision.outdoor_target_met"
	// EventVisionAlert 视力预警事件
	EventVisionAlert EventType = "vision.alert"
	// EventRewardBadgeEarned 勋章获得事件
	EventRewardBadgeEarned EventType = "reward.badge_earned"
	// EventRewardPointsEarned 积分获得事件
	EventRewardPointsEarned EventType = "reward.points_earned"
	// EventRewardChallengeCompleted 挑战完成事件
	EventRewardChallengeCompleted EventType = "reward.challenge_completed"
)

// Event 事件结构体，用于在系统间传递业务事件
type Event struct {
	// ID 事件唯一标识
	ID string `json:"id"`
	// Type 事件类型
	Type EventType `json:"type"`
	// Payload 事件负载数据，使用json.RawMessage延迟解析
	Payload json.RawMessage `json:"payload"`
	// Timestamp 事件发生时间
	Timestamp time.Time `json:"timestamp"`
	// Source 事件来源服务标识
	Source string `json:"source"`
}

// NewEvent 创建新的事件实例，自动生成UUID和时间戳
// eventType: 事件类型，payload: 事件负载对象（将被序列化为JSON）
func NewEvent(eventType EventType, payload interface{}) (*Event, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return &Event{
		ID:        uuid.New().String(),
		Type:      eventType,
		Payload:   data,
		Timestamp: time.Now(),
	}, nil
}

// EventPayload 事件负载辅助结构体，提供通用的事件数据封装
type EventPayload struct {
	// UserID 关联的用户ID
	UserID string `json:"user_id"`
	// Data 事件附加数据
	Data interface{} `json:"data,omitempty"`
}

// NewEventPayload 创建事件负载实例
// userID: 关联用户ID，data: 附加数据
func NewEventPayload(userID string, data interface{}) *EventPayload {
	return &EventPayload{
		UserID: userID,
		Data:   data,
	}
}
