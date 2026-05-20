package training

import (
	domain "github.com/kidsfit/api/internal/domain/training"
)

// ExerciseRecordToDTO 将运动记录领域模型转换为运动记录DTO
func ExerciseRecordToDTO(record *domain.ExerciseRecord) *ExerciseDTO {
	if record == nil {
		return nil
	}
	return &ExerciseDTO{
		ID:               record.ID,
		UserID:           record.UserID,
		Type:             string(record.Type),
		DurationSeconds:  record.DurationSeconds,
		Count:            record.Count,
		Score:            record.Score,
		RhythmScore:      record.RhythmScore,
		AmplitudeScore:   record.AmplitudeScore,
		SymmetryScore:    record.SymmetryScore,
		ContinuityScore:  record.ContinuityScore,
		Corrections:      record.Corrections,
		IsOffline:        record.IsOffline,
		StartedAt:        record.StartedAt,
		CompletedAt:      record.CompletedAt,
		CreatedAt:        record.CreatedAt,
	}
}

// DTOToExerciseRecord 将运动记录DTO转换为运动记录领域模型
func DTOToExerciseRecord(dto *ExerciseDTO) *domain.ExerciseRecord {
	if dto == nil {
		return nil
	}
	return &domain.ExerciseRecord{
		ID:              dto.ID,
		UserID:          dto.UserID,
		Type:            domain.ExerciseType(dto.Type),
		DurationSeconds: dto.DurationSeconds,
		Count:           dto.Count,
		Score:           dto.Score,
		RhythmScore:     dto.RhythmScore,
		AmplitudeScore:  dto.AmplitudeScore,
		SymmetryScore:   dto.SymmetryScore,
		ContinuityScore: dto.ContinuityScore,
		Corrections:     dto.Corrections,
		IsOffline:       dto.IsOffline,
		StartedAt:       dto.StartedAt,
		CompletedAt:     dto.CompletedAt,
		CreatedAt:       dto.CreatedAt,
	}
}

// CreateExerciseRequestToRecord 将创建运动记录请求DTO转换为运动记录领域模型
func CreateExerciseRequestToRecord(req *CreateExerciseRequest, userID string) *domain.ExerciseRecord {
	if req == nil {
		return nil
	}
	record := domain.NewExerciseRecord(userID, domain.ExerciseType(req.Type))
	record.DurationSeconds = req.DurationSeconds
	record.Count = req.Count
	record.Score = req.Score
	record.RhythmScore = req.RhythmScore
	record.AmplitudeScore = req.AmplitudeScore
	record.SymmetryScore = req.SymmetryScore
	record.ContinuityScore = req.ContinuityScore
	record.Corrections = req.Corrections
	record.IsOffline = req.IsOffline
	record.StartedAt = req.StartedAt
	record.CompletedAt = req.CompletedAt
	return record
}

// TrainingPlanToDTO 将训练计划领域模型转换为训练计划DTO
func TrainingPlanToDTO(plan *domain.TrainingPlan) *TrainingPlanDTO {
	if plan == nil {
		return nil
	}
	dto := &TrainingPlanDTO{
		ID:             plan.ID,
		Date:           plan.Date,
		Status:         string(plan.Status),
		TotalDuration:  plan.TotalDuration,
		ActualDuration: plan.ActualDuration,
	}
	// 转换所有阶段的运动项目
	items := make([]ExerciseItemDTO, 0)
	for _, item := range plan.WarmupItems {
		items = append(items, *ExerciseItemToDTO(&item))
	}
	for _, item := range plan.MainItems {
		items = append(items, *ExerciseItemToDTO(&item))
	}
	for _, item := range plan.CooldownItems {
		items = append(items, *ExerciseItemToDTO(&item))
	}
	dto.Items = items
	return dto
}

// ExerciseItemToDTO 将运动项目领域模型转换为运动项目DTO
func ExerciseItemToDTO(item *domain.ExerciseItem) *ExerciseItemDTO {
	if item == nil {
		return nil
	}
	return &ExerciseItemDTO{
		ID:           item.ID,
		Type:         string(item.Type),
		Name:         item.Name,
		DurationSec:  item.DurationSec,
		TargetCount:  item.TargetCount,
		Difficulty:   item.Difficulty,
		Tips:         item.Tips,
		Order:        item.Order,
		Phase:        string(item.Phase),
	}
}

// FitnessAssessmentToDTO 将体能评估领域模型转换为体能评估DTO
func FitnessAssessmentToDTO(assessment *domain.FitnessAssessment) *FitnessAssessmentDTO {
	if assessment == nil {
		return nil
	}
	return &FitnessAssessmentDTO{
		ID:           assessment.ID,
		UserID:       assessment.UserID,
		Endurance:    assessment.Endurance,
		Agility:      assessment.Agility,
		Strength:     assessment.Strength,
		Speed:        assessment.Speed,
		Coordination: assessment.Coordination,
		Balance:      assessment.Balance,
		Flexibility:  assessment.Flexibility,
		AssessedAt:   assessment.AssessedAt,
		CreatedAt:    assessment.CreatedAt,
	}
}
