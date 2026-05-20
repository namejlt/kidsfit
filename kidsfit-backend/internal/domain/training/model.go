package training

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ExerciseType 运动类型枚举
type ExerciseType string

const (
	// ExerciseTypeJumpRope 跳绳
	ExerciseTypeJumpRope ExerciseType = "jump_rope"
	// ExerciseTypeJumpingJack 开合跳
	ExerciseTypeJumpingJack ExerciseType = "jumping_jack"
	// ExerciseTypeSquat 深蹲
	ExerciseTypeSquat ExerciseType = "squat"
	// ExerciseTypeSitUp 仰卧起坐
	ExerciseTypeSitUp ExerciseType = "sit_up"
	// ExerciseTypeHighKnee 高抬腿
	ExerciseTypeHighKnee ExerciseType = "high_knee"
	// ExerciseTypePushUp 俯卧撑
	ExerciseTypePushUp ExerciseType = "push_up"
)

// PlanStatus 训练计划状态枚举
type PlanStatus string

const (
	// PlanStatusPending 待完成
	PlanStatusPending PlanStatus = "pending"
	// PlanStatusCompleted 已完成
	PlanStatusCompleted PlanStatus = "completed"
	// PlanStatusSkipped 已跳过
	PlanStatusSkipped PlanStatus = "skipped"
)

// ExercisePhase 运动阶段枚举
type ExercisePhase string

const (
	// ExercisePhaseWarmup 热身阶段
	ExercisePhaseWarmup ExercisePhase = "warmup"
	// ExercisePhaseMain 主体阶段
	ExercisePhaseMain ExercisePhase = "main"
	// ExercisePhaseCooldown 放松阶段
	ExercisePhaseCooldown ExercisePhase = "cooldown"
)

// ExerciseRecord 运动记录领域模型
type ExerciseRecord struct {
	ID               string      `json:"id" db:"id"`
	UserID           string      `json:"user_id" db:"user_id"`
	Type             ExerciseType `json:"type" db:"type"`
	DurationSeconds  int         `json:"duration_seconds" db:"duration_seconds"`
	Count            int         `json:"count" db:"count"`
	Score            int         `json:"score" db:"score"`
	RhythmScore      int         `json:"rhythm_score" db:"rhythm_score"`
	AmplitudeScore   int         `json:"amplitude_score" db:"amplitude_score"`
	SymmetryScore    int         `json:"symmetry_score" db:"symmetry_score"`
	ContinuityScore  int         `json:"continuity_score" db:"continuity_score"`
	Corrections      []string    `json:"corrections" db:"corrections"`
	IsOffline        bool        `json:"is_offline" db:"is_offline"`
	StartedAt        time.Time   `json:"started_at" db:"started_at"`
	CompletedAt      time.Time   `json:"completed_at" db:"completed_at"`
	CreatedAt        time.Time   `json:"created_at" db:"created_at"`
}

// NewExerciseRecord 创建运动记录实例
func NewExerciseRecord(userID string, exerciseType ExerciseType) *ExerciseRecord {
	return &ExerciseRecord{
		ID:          uuid.New().String(),
		UserID:      userID,
		Type:        exerciseType,
		Corrections: []string{},
		CreatedAt:   time.Now(),
	}
}

// ScoreGrade 根据分数返回等级：S/A/B/C/D
func (r *ExerciseRecord) ScoreGrade() string {
	switch {
	case r.Score >= 90:
		return "S"
	case r.Score >= 80:
		return "A"
	case r.Score >= 70:
		return "B"
	case r.Score >= 60:
		return "C"
	default:
		return "D"
	}
}

// ValidateScore 校验分数是否在0-100范围内
func (r *ExerciseRecord) ValidateScore() error {
	if r.Score < 0 || r.Score > 100 {
		return errors.New("分数必须在0-100之间")
	}
	return nil
}

// TrainingPlan 训练计划领域模型
type TrainingPlan struct {
	ID             string         `json:"id" db:"id"`
	UserID         string         `json:"user_id" db:"user_id"`
	Date           time.Time      `json:"date" db:"date"`
	Status         PlanStatus     `json:"status" db:"status"`
	TotalDuration  int            `json:"total_duration" db:"total_duration"`
	ActualDuration int            `json:"actual_duration" db:"actual_duration"`
	WarmupItems    []ExerciseItem `json:"warmup_items" db:"warmup_items"`
	MainItems      []ExerciseItem `json:"main_items" db:"main_items"`
	CooldownItems  []ExerciseItem `json:"cooldown_items" db:"cooldown_items"`
	CompletedAt    *time.Time     `json:"completed_at,omitempty" db:"completed_at"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
}

// NewTrainingPlan 创建训练计划实例
func NewTrainingPlan(userID string, date time.Time) *TrainingPlan {
	return &TrainingPlan{
		ID:        uuid.New().String(),
		UserID:    userID,
		Date:      date,
		Status:    PlanStatusPending,
		CreatedAt: time.Now(),
	}
}

// ExerciseItem 运动项目领域模型
type ExerciseItem struct {
	ID           string        `json:"id" db:"id"`
	PlanID       string        `json:"plan_id" db:"plan_id"`
	Type         ExerciseType  `json:"type" db:"type"`
	Name         string        `json:"name" db:"name"`
	DurationSec  int           `json:"duration_sec" db:"duration_sec"`
	TargetCount  int           `json:"target_count" db:"target_count"`
	Difficulty   int           `json:"difficulty" db:"difficulty"`
	Tips         string        `json:"tips" db:"tips"`
	Order        int           `json:"order" db:"order"`
	Phase        ExercisePhase `json:"phase" db:"phase"`
}

// NewExerciseItem 创建运动项目实例
func NewExerciseItem(planID string, exerciseType ExerciseType, phase ExercisePhase, order int) *ExerciseItem {
	return &ExerciseItem{
		ID:     uuid.New().String(),
		PlanID: planID,
		Type:   exerciseType,
		Phase:  phase,
		Order:  order,
	}
}

// ValidateDifficulty 校验难度是否在1-5范围内
func (ei *ExerciseItem) ValidateDifficulty() error {
	if ei.Difficulty < 1 || ei.Difficulty > 5 {
		return errors.New("难度必须在1-5之间")
	}
	return nil
}

// FitnessAssessment 体能评估领域模型
type FitnessAssessment struct {
	ID           string    `json:"id" db:"id"`
	UserID       string    `json:"user_id" db:"user_id"`
	Endurance    int       `json:"endurance" db:"endurance"`
	Agility      int       `json:"agility" db:"agility"`
	Strength     int       `json:"strength" db:"strength"`
	Speed        int       `json:"speed" db:"speed"`
	Coordination int       `json:"coordination" db:"coordination"`
	Balance      int       `json:"balance" db:"balance"`
	Flexibility  int       `json:"flexibility" db:"flexibility"`
	AssessedAt   time.Time `json:"assessed_at" db:"assessed_at"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// NewFitnessAssessment 创建体能评估实例
func NewFitnessAssessment(userID string) *FitnessAssessment {
	return &FitnessAssessment{
		ID:        uuid.New().String(),
		UserID:    userID,
		CreatedAt: time.Now(),
	}
}

// ValidateScores 校验所有评估分数是否在1-10范围内
func (fa *FitnessAssessment) ValidateScores() error {
	scores := []int{fa.Endurance, fa.Agility, fa.Strength, fa.Speed, fa.Coordination, fa.Balance, fa.Flexibility}
	for _, s := range scores {
		if s < 1 || s > 10 {
			return errors.New("评估分数必须在1-10之间")
		}
	}
	return nil
}
