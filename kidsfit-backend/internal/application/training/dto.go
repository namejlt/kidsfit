package training

import "time"

// CreateExerciseRequest 创建运动记录请求DTO
type CreateExerciseRequest struct {
	// Type 运动类型
	Type string `json:"type" binding:"required"`
	// DurationSeconds 运动时长（秒）
	DurationSeconds int `json:"duration_seconds" binding:"required,min=1"`
	// Count 运动次数
	Count int `json:"count" binding:"required,min=1"`
	// Score 综合评分
	Score int `json:"score" binding:"min=0,max=100"`
	// RhythmScore 节奏评分
	RhythmScore int `json:"rhythm_score"`
	// AmplitudeScore 幅度评分
	AmplitudeScore int `json:"amplitude_score"`
	// SymmetryScore 对称性评分
	SymmetryScore int `json:"symmetry_score"`
	// ContinuityScore 连续性评分
	ContinuityScore int `json:"continuity_score"`
	// Corrections 纠正建议列表
	Corrections []string `json:"corrections"`
	// IsOffline 是否离线记录
	IsOffline bool `json:"is_offline"`
	// StartedAt 开始时间
	StartedAt time.Time `json:"started_at" binding:"required"`
	// CompletedAt 完成时间
	CompletedAt time.Time `json:"completed_at" binding:"required"`
}

// ExerciseDTO 运动记录响应DTO
type ExerciseDTO struct {
	// ID 记录ID
	ID string `json:"id"`
	// UserID 用户ID
	UserID string `json:"user_id"`
	// Type 运动类型
	Type string `json:"type"`
	// DurationSeconds 运动时长（秒）
	DurationSeconds int `json:"duration_seconds"`
	// Count 运动次数
	Count int `json:"count"`
	// Score 综合评分
	Score int `json:"score"`
	// RhythmScore 节奏评分
	RhythmScore int `json:"rhythm_score"`
	// AmplitudeScore 幅度评分
	AmplitudeScore int `json:"amplitude_score"`
	// SymmetryScore 对称性评分
	SymmetryScore int `json:"symmetry_score"`
	// ContinuityScore 连续性评分
	ContinuityScore int `json:"continuity_score"`
	// Corrections 纠正建议列表
	Corrections []string `json:"corrections"`
	// IsOffline 是否离线记录
	IsOffline bool `json:"is_offline"`
	// StartedAt 开始时间
	StartedAt time.Time `json:"started_at"`
	// CompletedAt 完成时间
	CompletedAt time.Time `json:"completed_at"`
	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`
	// PointsEarned 获得积分
	PointsEarned int `json:"points_earned"`
	// BadgesEarned 获得徽章列表
	BadgesEarned []string `json:"badges_earned"`
	// IsRecordBroken 是否打破个人纪录
	IsRecordBroken bool `json:"is_record_broken"`
}

// TrainingPlanDTO 训练计划响应DTO
type TrainingPlanDTO struct {
	// ID 计划ID
	ID string `json:"id"`
	// Date 计划日期
	Date time.Time `json:"date"`
	// Status 计划状态
	Status string `json:"status"`
	// TotalDuration 总时长（秒）
	TotalDuration int `json:"total_duration"`
	// ActualDuration 实际时长（秒）
	ActualDuration int `json:"actual_duration"`
	// Items 训练项目列表
	Items []ExerciseItemDTO `json:"items"`
}

// ExerciseItemDTO 运动项目响应DTO
type ExerciseItemDTO struct {
	// ID 项目ID
	ID string `json:"id"`
	// Type 运动类型
	Type string `json:"type"`
	// Name 运动名称
	Name string `json:"name"`
	// DurationSec 时长（秒）
	DurationSec int `json:"duration_sec"`
	// TargetCount 目标次数
	TargetCount int `json:"target_count"`
	// Difficulty 难度等级（1-5）
	Difficulty int `json:"difficulty"`
	// Tips 运动提示
	Tips string `json:"tips"`
	// Order 排序序号
	Order int `json:"order"`
	// Phase 运动阶段（warmup/main/cooldown）
	Phase string `json:"phase"`
}

// FitnessAssessmentDTO 体能评估响应DTO
type FitnessAssessmentDTO struct {
	// ID 评估ID
	ID string `json:"id"`
	// UserID 用户ID
	UserID string `json:"user_id"`
	// Endurance 耐力评分（1-10）
	Endurance int `json:"endurance"`
	// Agility 敏捷评分（1-10）
	Agility int `json:"agility"`
	// Strength 力量评分（1-10）
	Strength int `json:"strength"`
	// Speed 速度评分（1-10）
	Speed int `json:"speed"`
	// Coordination 协调评分（1-10）
	Coordination int `json:"coordination"`
	// Balance 平衡评分（1-10）
	Balance int `json:"balance"`
	// Flexibility 柔韧评分（1-10）
	Flexibility int `json:"flexibility"`
	// AssessedAt 评估时间
	AssessedAt time.Time `json:"assessed_at"`
	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`
}

// WeeklyStatsDTO 周统计响应DTO
type WeeklyStatsDTO struct {
	// TotalExercises 运动总次数
	TotalExercises int `json:"total_exercises"`
	// TotalDuration 总时长（秒）
	TotalDuration int `json:"total_duration"`
	// TotalCount 运动总次数（计数类）
	TotalCount int `json:"total_count"`
	// AverageScore 平均评分
	AverageScore float64 `json:"average_score"`
	// ActiveDays 活跃天数
	ActiveDays int `json:"active_days"`
	// DailyBreakdown 每日明细
	DailyBreakdown []DailyStats `json:"daily_breakdown"`
}

// DailyStats 每日统计数据
type DailyStats struct {
	// Date 日期
	Date string `json:"date"`
	// ExerciseCount 运动次数
	ExerciseCount int `json:"exercise_count"`
	// Duration 运动时长（秒）
	Duration int `json:"duration"`
	// Count 运动计数
	Count int `json:"count"`
	// Score 平均评分
	Score float64 `json:"score"`
}

// MonthlyStatsDTO 月统计响应DTO
type MonthlyStatsDTO struct {
	// TotalExercises 运动总次数
	TotalExercises int `json:"total_exercises"`
	// TotalDuration 总时长（秒）
	TotalDuration int `json:"total_duration"`
	// TotalCount 运动总次数（计数类）
	TotalCount int `json:"total_count"`
	// AverageScore 平均评分
	AverageScore float64 `json:"average_score"`
	// ActiveDays 活跃天数
	ActiveDays int `json:"active_days"`
	// WeeklyBreakdown 每周明细
	WeeklyBreakdown []WeeklyBreakdown `json:"weekly_breakdown"`
}

// WeeklyBreakdown 每周统计数据
type WeeklyBreakdown struct {
	// WeekStart 周开始日期
	WeekStart string `json:"week_start"`
	// ExerciseCount 运动次数
	ExerciseCount int `json:"exercise_count"`
	// Duration 运动时长（秒）
	Duration int `json:"duration"`
	// Count 运动计数
	Count int `json:"count"`
	// Score 平均评分
	Score float64 `json:"score"`
}

// PersonalBestDTO 个人最佳记录响应DTO
type PersonalBestDTO struct {
	// Type 运动类型
	Type string `json:"type"`
	// BestScore 最佳评分
	BestScore int `json:"best_score"`
	// BestCount 最佳次数
	BestCount int `json:"best_count"`
	// AchievedAt 达成时间
	AchievedAt time.Time `json:"achieved_at"`
}
