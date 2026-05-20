package training

import (
	"testing"
)

// TestExerciseRecord_ScoreGrade 测试评分等级划分
func TestExerciseRecord_ScoreGrade(t *testing.T) {
	tests := []struct {
		name      string
		score     int
		wantGrade string
	}{
		{"S级-90分", 90, "S"},
		{"S级-95分", 95, "S"},
		{"S级-100分", 100, "S"},
		{"A级-80分", 80, "A"},
		{"A级-89分", 89, "A"},
		{"B级-70分", 70, "B"},
		{"B级-79分", 79, "B"},
		{"C级-60分", 60, "C"},
		{"C级-69分", 69, "C"},
		{"D级-0分", 0, "D"},
		{"D级-59分", 59, "D"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ExerciseRecord{Score: tt.score}
			if got := r.ScoreGrade(); got != tt.wantGrade {
				t.Errorf("分数 %d 等级应为 %s，实际为 %s", tt.score, tt.wantGrade, got)
			}
		})
	}
}

// TestExerciseRecord_ValidateScore 测试分数校验
func TestExerciseRecord_ValidateScore(t *testing.T) {
	t.Run("有效分数0-100", func(t *testing.T) {
		validScores := []int{0, 50, 100}
		for _, score := range validScores {
			r := &ExerciseRecord{Score: score}
			if err := r.ValidateScore(); err != nil {
				t.Errorf("分数 %d 应为有效值，实际返回错误: %v", score, err)
			}
		}
	})

	t.Run("无效分数-1和101", func(t *testing.T) {
		invalidScores := []int{-1, 101}
		for _, score := range invalidScores {
			r := &ExerciseRecord{Score: score}
			if err := r.ValidateScore(); err == nil {
				t.Errorf("分数 %d 应为无效值，但未返回错误", score)
			}
		}
	})
}

// TestExerciseItem_ValidateDifficulty 测试难度校验
func TestExerciseItem_ValidateDifficulty(t *testing.T) {
	t.Run("有效难度1-5", func(t *testing.T) {
		validDifficulties := []int{1, 2, 3, 4, 5}
		for _, d := range validDifficulties {
			ei := &ExerciseItem{Difficulty: d}
			if err := ei.ValidateDifficulty(); err != nil {
				t.Errorf("难度 %d 应为有效值，实际返回错误: %v", d, err)
			}
		}
	})

	t.Run("无效难度0和6", func(t *testing.T) {
		invalidDifficulties := []int{0, 6}
		for _, d := range invalidDifficulties {
			ei := &ExerciseItem{Difficulty: d}
			if err := ei.ValidateDifficulty(); err == nil {
				t.Errorf("难度 %d 应为无效值，但未返回错误", d)
			}
		}
	})
}

// TestFitnessAssessment_ValidateScores 测试体能评估分数校验
func TestFitnessAssessment_ValidateScores(t *testing.T) {
	t.Run("所有分数在1-10范围内有效", func(t *testing.T) {
		fa := &FitnessAssessment{
			Endurance:    5,
			Agility:      5,
			Strength:     5,
			Speed:        5,
			Coordination: 5,
			Balance:      5,
			Flexibility:  5,
		}
		if err := fa.ValidateScores(); err != nil {
			t.Errorf("有效分数应返回nil，实际返回错误: %v", err)
		}
	})

	t.Run("边界值1和10有效", func(t *testing.T) {
		fa := &FitnessAssessment{
			Endurance:    1,
			Agility:      10,
			Strength:     1,
			Speed:        10,
			Coordination: 1,
			Balance:      10,
			Flexibility:  1,
		}
		if err := fa.ValidateScores(); err != nil {
			t.Errorf("边界值1和10应有效，实际返回错误: %v", err)
		}
	})

	t.Run("分数为0时无效", func(t *testing.T) {
		fa := &FitnessAssessment{
			Endurance:    0,
			Agility:      5,
			Strength:     5,
			Speed:        5,
			Coordination: 5,
			Balance:      5,
			Flexibility:  5,
		}
		if err := fa.ValidateScores(); err == nil {
			t.Error("分数为0时应返回错误")
		}
	})

	t.Run("分数为11时无效", func(t *testing.T) {
		fa := &FitnessAssessment{
			Endurance:    5,
			Agility:      5,
			Strength:     5,
			Speed:        5,
			Coordination: 5,
			Balance:      5,
			Flexibility:  11,
		}
		if err := fa.ValidateScores(); err == nil {
			t.Error("分数为11时应返回错误")
		}
	})
}
