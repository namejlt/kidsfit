package training

import (
	"context"
	"fmt"
	"time"

	domain "github.com/kidsfit/api/internal/domain/training"
)

// ageGroupConfig 年龄组训练强度配置
type ageGroupConfig struct {
	WarmupCount     int
	MainCount       int
	CooldownCount   int
	TotalDurationSec int
	MaxDifficulty   int
}

// 预定义年龄组配置
var ageGroupConfigs = map[string]ageGroupConfig{
	"3-6": {
		WarmupCount:     2,
		MainCount:       3,
		CooldownCount:   1,
		TotalDurationSec: 5 * 60,
		MaxDifficulty:   2,
	},
	"7-9": {
		WarmupCount:     2,
		MainCount:       4,
		CooldownCount:   1,
		TotalDurationSec: 15 * 60,
		MaxDifficulty:   3,
	},
	"10-12": {
		WarmupCount:     3,
		MainCount:       5,
		CooldownCount:   2,
		TotalDurationSec: 20 * 60,
		MaxDifficulty:   5,
	},
}

// exerciseDefinition 预定义运动项目定义
type exerciseDefinition struct {
	Type         domain.ExerciseType
	Name         string
	Difficulty   int
	DurationSec  int
	TargetCount  int
	Tips         string
	Phase        domain.ExercisePhase
	Dimensions   []string
}

// 预定义运动项目库
var exerciseLibrary = []exerciseDefinition{
	// === 热身运动 ===
	{Type: domain.ExerciseTypeJumpingJack, Name: "开合跳（轻松版）", Difficulty: 1, DurationSec: 60, TargetCount: 20, Tips: "保持节奏均匀，不要过快", Phase: domain.ExercisePhaseWarmup, Dimensions: []string{"endurance", "coordination"}},
	{Type: domain.ExerciseTypeHighKnee, Name: "原地高抬腿（慢速）", Difficulty: 1, DurationSec: 60, TargetCount: 20, Tips: "膝盖尽量抬高，保持平衡", Phase: domain.ExercisePhaseWarmup, Dimensions: []string{"endurance", "speed"}},
	{Type: domain.ExerciseTypeJumpingJack, Name: "开合跳（标准版）", Difficulty: 2, DurationSec: 90, TargetCount: 30, Tips: "手臂伸展到位，落地轻柔", Phase: domain.ExercisePhaseWarmup, Dimensions: []string{"endurance", "coordination"}},
	{Type: domain.ExerciseTypeHighKnee, Name: "原地高抬腿（标准）", Difficulty: 2, DurationSec: 90, TargetCount: 30, Tips: "加快频率，保持呼吸均匀", Phase: domain.ExercisePhaseWarmup, Dimensions: []string{"endurance", "speed"}},
	{Type: domain.ExerciseTypeJumpingJack, Name: "开合跳（加强版）", Difficulty: 3, DurationSec: 120, TargetCount: 40, Tips: "提高速度，注意手脚协调", Phase: domain.ExercisePhaseWarmup, Dimensions: []string{"endurance", "coordination"}},

	// === 主体运动 - 跳绳 ===
	{Type: domain.ExerciseTypeJumpRope, Name: "跳绳（入门）", Difficulty: 1, DurationSec: 120, TargetCount: 30, Tips: "手腕发力，脚尖着地", Phase: domain.ExercisePhaseMain, Dimensions: []string{"endurance", "coordination", "speed"}},
	{Type: domain.ExerciseTypeJumpRope, Name: "跳绳（基础）", Difficulty: 2, DurationSec: 180, TargetCount: 50, Tips: "保持节奏，不要断绳", Phase: domain.ExercisePhaseMain, Dimensions: []string{"endurance", "coordination", "speed"}},
	{Type: domain.ExerciseTypeJumpRope, Name: "跳绳（进阶）", Difficulty: 3, DurationSec: 240, TargetCount: 80, Tips: "尝试交替脚跳，提高连续性", Phase: domain.ExercisePhaseMain, Dimensions: []string{"endurance", "coordination", "speed"}},
	{Type: domain.ExerciseTypeJumpRope, Name: "跳绳（挑战）", Difficulty: 4, DurationSec: 300, TargetCount: 120, Tips: "保持高速跳绳，挑战个人纪录", Phase: domain.ExercisePhaseMain, Dimensions: []string{"endurance", "coordination", "speed"}},
	{Type: domain.ExerciseTypeJumpRope, Name: "跳绳（极限）", Difficulty: 5, DurationSec: 360, TargetCount: 160, Tips: "尝试花式跳法，突破自我", Phase: domain.ExercisePhaseMain, Dimensions: []string{"endurance", "coordination", "speed"}},

	// === 主体运动 - 深蹲 ===
	{Type: domain.ExerciseTypeSquat, Name: "深蹲（入门）", Difficulty: 1, DurationSec: 90, TargetCount: 10, Tips: "膝盖不超过脚尖，背部挺直", Phase: domain.ExercisePhaseMain, Dimensions: []string{"strength", "balance"}},
	{Type: domain.ExerciseTypeSquat, Name: "深蹲（基础）", Difficulty: 2, DurationSec: 120, TargetCount: 15, Tips: "下蹲至大腿与地面平行", Phase: domain.ExercisePhaseMain, Dimensions: []string{"strength", "balance"}},
	{Type: domain.ExerciseTypeSquat, Name: "深蹲（进阶）", Difficulty: 3, DurationSec: 150, TargetCount: 20, Tips: "增加速度，注意动作标准", Phase: domain.ExercisePhaseMain, Dimensions: []string{"strength", "balance", "endurance"}},
	{Type: domain.ExerciseTypeSquat, Name: "深蹲（挑战）", Difficulty: 4, DurationSec: 180, TargetCount: 30, Tips: "尝试跳跃深蹲，增强爆发力", Phase: domain.ExercisePhaseMain, Dimensions: []string{"strength", "balance", "endurance", "speed"}},
	{Type: domain.ExerciseTypeSquat, Name: "深蹲（极限）", Difficulty: 5, DurationSec: 240, TargetCount: 40, Tips: "挑战单腿深蹲，提升核心力量", Phase: domain.ExercisePhaseMain, Dimensions: []string{"strength", "balance", "endurance"}},

	// === 主体运动 - 仰卧起坐 ===
	{Type: domain.ExerciseTypeSitUp, Name: "仰卧起坐（入门）", Difficulty: 1, DurationSec: 90, TargetCount: 10, Tips: "双手放耳侧，不要抱头", Phase: domain.ExercisePhaseMain, Dimensions: []string{"strength", "flexibility"}},
	{Type: domain.ExerciseTypeSitUp, Name: "仰卧起坐（基础）", Difficulty: 2, DurationSec: 120, TargetCount: 15, Tips: "起身时呼气，下落时吸气", Phase: domain.ExercisePhaseMain, Dimensions: []string{"strength", "flexibility"}},
	{Type: domain.ExerciseTypeSitUp, Name: "仰卧起坐（进阶）", Difficulty: 3, DurationSec: 150, TargetCount: 25, Tips: "增加速度，保持腹部发力", Phase: domain.ExercisePhaseMain, Dimensions: []string{"strength", "endurance"}},
	{Type: domain.ExerciseTypeSitUp, Name: "仰卧起坐（挑战）", Difficulty: 4, DurationSec: 180, TargetCount: 35, Tips: "尝试交叉仰卧起坐，锻炼侧腹", Phase: domain.ExercisePhaseMain, Dimensions: []string{"strength", "endurance", "coordination"}},

	// === 主体运动 - 高抬腿 ===
	{Type: domain.ExerciseTypeHighKnee, Name: "高抬腿（入门）", Difficulty: 1, DurationSec: 60, TargetCount: 20, Tips: "膝盖抬至腰部高度", Phase: domain.ExercisePhaseMain, Dimensions: []string{"speed", "endurance"}},
	{Type: domain.ExerciseTypeHighKnee, Name: "高抬腿（基础）", Difficulty: 2, DurationSec: 90, TargetCount: 30, Tips: "加快频率，保持上身稳定", Phase: domain.ExercisePhaseMain, Dimensions: []string{"speed", "endurance", "coordination"}},
	{Type: domain.ExerciseTypeHighKnee, Name: "高抬腿（进阶）", Difficulty: 3, DurationSec: 120, TargetCount: 50, Tips: "提高速度，保持呼吸节奏", Phase: domain.ExercisePhaseMain, Dimensions: []string{"speed", "endurance", "coordination"}},
	{Type: domain.ExerciseTypeHighKnee, Name: "高抬腿（挑战）", Difficulty: 4, DurationSec: 150, TargetCount: 70, Tips: "冲刺速度，挑战极限", Phase: domain.ExercisePhaseMain, Dimensions: []string{"speed", "endurance"}},

	// === 主体运动 - 俯卧撑 ===
	{Type: domain.ExerciseTypePushUp, Name: "俯卧撑（入门）", Difficulty: 1, DurationSec: 60, TargetCount: 5, Tips: "可以跪姿俯卧撑开始", Phase: domain.ExercisePhaseMain, Dimensions: []string{"strength", "coordination"}},
	{Type: domain.ExerciseTypePushUp, Name: "俯卧撑（基础）", Difficulty: 2, DurationSec: 90, TargetCount: 10, Tips: "身体保持一条直线", Phase: domain.ExercisePhaseMain, Dimensions: []string{"strength", "coordination"}},
	{Type: domain.ExerciseTypePushUp, Name: "俯卧撑（进阶）", Difficulty: 3, DurationSec: 120, TargetCount: 15, Tips: "下降时慢速控制，上升时快速推起", Phase: domain.ExercisePhaseMain, Dimensions: []string{"strength", "endurance"}},
	{Type: domain.ExerciseTypePushUp, Name: "俯卧撑（挑战）", Difficulty: 4, DurationSec: 150, TargetCount: 25, Tips: "尝试宽距俯卧撑，增加难度", Phase: domain.ExercisePhaseMain, Dimensions: []string{"strength", "endurance", "coordination"}},
	{Type: domain.ExerciseTypePushUp, Name: "俯卧撑（极限）", Difficulty: 5, DurationSec: 180, TargetCount: 35, Tips: "尝试击掌俯卧撑，挑战爆发力", Phase: domain.ExercisePhaseMain, Dimensions: []string{"strength", "speed", "coordination"}},

	// === 主体运动 - 开合跳 ===
	{Type: domain.ExerciseTypeJumpingJack, Name: "开合跳（入门）", Difficulty: 1, DurationSec: 90, TargetCount: 20, Tips: "动作幅度适中，保持节奏", Phase: domain.ExercisePhaseMain, Dimensions: []string{"endurance", "coordination"}},
	{Type: domain.ExerciseTypeJumpingJack, Name: "开合跳（基础）", Difficulty: 2, DurationSec: 120, TargetCount: 35, Tips: "加快速度，注意手脚协调", Phase: domain.ExercisePhaseMain, Dimensions: []string{"endurance", "coordination", "agility"}},
	{Type: domain.ExerciseTypeJumpingJack, Name: "开合跳（进阶）", Difficulty: 3, DurationSec: 180, TargetCount: 50, Tips: "持续跳跃，提高心率", Phase: domain.ExercisePhaseMain, Dimensions: []string{"endurance", "coordination", "agility"}},

	// === 放松运动 ===
	{Type: domain.ExerciseTypeJumpingJack, Name: "慢速开合跳", Difficulty: 1, DurationSec: 60, TargetCount: 10, Tips: "放慢速度，调整呼吸", Phase: domain.ExercisePhaseCooldown, Dimensions: []string{"coordination"}},
	{Type: domain.ExerciseTypeHighKnee, Name: "慢速高抬腿", Difficulty: 1, DurationSec: 60, TargetCount: 10, Tips: "缓慢抬腿，放松肌肉", Phase: domain.ExercisePhaseCooldown, Dimensions: []string{"flexibility"}},
	{Type: domain.ExerciseTypeSquat, Name: "浅蹲拉伸", Difficulty: 1, DurationSec: 60, TargetCount: 8, Tips: "下蹲后保持5秒，拉伸腿部", Phase: domain.ExercisePhaseCooldown, Dimensions: []string{"flexibility", "balance"}},
	{Type: domain.ExerciseTypeSitUp, Name: "腹部拉伸", Difficulty: 1, DurationSec: 60, TargetCount: 5, Tips: "缓慢起身，拉伸腹部肌肉", Phase: domain.ExercisePhaseCooldown, Dimensions: []string{"flexibility"}},
}

// weakDimensionThreshold 弱项判定阈值，低于此分数的维度视为弱项
const weakDimensionThreshold = 5

// RecommendationService 训练计划推荐服务，根据体能评估和年龄组生成个性化训练计划
type RecommendationService struct{}

// NewRecommendationService 创建训练计划推荐服务实例
func NewRecommendationService() *RecommendationService {
	return &RecommendationService{}
}

// GeneratePlan 根据体能评估结果和年龄组生成个性化训练计划
// 包含热身、主体训练和放松三个阶段的项目选择，
// 并对评估中低于阈值的弱项维度增加针对性训练
func (s *RecommendationService) GeneratePlan(ctx context.Context, userID string, assessment *domain.FitnessAssessment, ageGroup string) (*domain.TrainingPlan, error) {
	config, ok := ageGroupConfigs[ageGroup]
	if !ok {
		return nil, fmt.Errorf("不支持的年龄组: %s", ageGroup)
	}

	plan := domain.NewTrainingPlan(userID, time.Now())
	plan.TotalDuration = config.TotalDurationSec

	// 分析弱项维度
	weakDimensions := s.analyzeWeakDimensions(assessment)

	// 选择各阶段运动项目
	warmupItems := s.selectExercisesByPhase(weakDimensions, ageGroup, domain.ExercisePhaseWarmup, config.WarmupCount)
	mainItems := s.selectExercisesByPhase(weakDimensions, ageGroup, domain.ExercisePhaseMain, config.MainCount)
	cooldownItems := s.selectExercisesByPhase(weakDimensions, ageGroup, domain.ExercisePhaseCooldown, config.CooldownCount)

	plan.WarmupItems = warmupItems
	plan.MainItems = mainItems
	plan.CooldownItems = cooldownItems

	return plan, nil
}

// analyzeWeakDimensions 从体能评估中识别低于阈值的弱项维度
func (s *RecommendationService) analyzeWeakDimensions(assessment *domain.FitnessAssessment) []string {
	if assessment == nil {
		return []string{}
	}
	var weakDims []string
	dimensionScores := map[string]int{
		"endurance":    assessment.Endurance,
		"agility":      assessment.Agility,
		"strength":     assessment.Strength,
		"speed":        assessment.Speed,
		"coordination": assessment.Coordination,
		"balance":      assessment.Balance,
		"flexibility":  assessment.Flexibility,
	}
	for dim, score := range dimensionScores {
		if score < weakDimensionThreshold {
			weakDims = append(weakDims, dim)
		}
	}
	return weakDims
}

// selectExercisesByPhase 根据指定阶段筛选并转换运动项目为领域模型
func (s *RecommendationService) selectExercisesByPhase(weakDimensions []string, ageGroup string, phase domain.ExercisePhase, count int) []domain.ExerciseItem {
	defs := s.selectExercises(weakDimensions, ageGroup, phase)
	if len(defs) > count {
		defs = defs[:count]
	}
	items := make([]domain.ExerciseItem, 0, len(defs))
	for i, def := range defs {
		item := domain.NewExerciseItem("", def.Type, def.Phase, i+1)
		item.Name = def.Name
		item.DurationSec = def.DurationSec
		item.TargetCount = def.TargetCount
		item.Difficulty = def.Difficulty
		item.Tips = def.Tips
		items = append(items, *item)
	}
	return items
}

// selectExercises 从预定义项目库中筛选适合指定阶段的运动项目
// 弱项维度的项目会被优先选中以加强对应能力
func (s *RecommendationService) selectExercises(weakDimensions []string, ageGroup string, phase domain.ExercisePhase) []exerciseDefinition {
	config := ageGroupConfigs[ageGroup]
	var candidates []exerciseDefinition
	var others []exerciseDefinition

	for _, ex := range exerciseLibrary {
		if ex.Phase != phase || ex.Difficulty > config.MaxDifficulty {
			continue
		}
		if s.matchesWeakDimension(ex.Dimensions, weakDimensions) {
			candidates = append(candidates, ex)
		} else {
			others = append(others, ex)
		}
	}

	result := candidates
	for _, ex := range others {
		result = append(result, ex)
	}
	return result
}

// matchesWeakDimension 判断运动项目的训练维度是否与弱项列表匹配
func (s *RecommendationService) matchesWeakDimension(dimensions, weakDimensions []string) bool {
	if len(weakDimensions) == 0 {
		return false
	}
	for _, d := range dimensions {
		for _, w := range weakDimensions {
			if d == w {
				return true
			}
		}
	}
	return false
}
