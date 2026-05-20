package training

import (
	"context"
	"time"

	"github.com/kidsfit/api/internal/domain/training"
	appErrors "github.com/kidsfit/api/internal/pkg/errors"
	"github.com/kidsfit/api/internal/pkg/response"
)

// TrainingAppService 训练应用服务，负责训练相关的业务逻辑编排
type TrainingAppService struct {
	exerciseRepo  training.ExerciseRecordRepository
	planRepo      training.TrainingPlanRepository
	assessmentRepo training.FitnessAssessmentRepository
}

// NewTrainingAppService 创建训练应用服务实例
// exerciseRepo: 运动记录仓储，planRepo: 训练计划仓储，assessmentRepo: 体能评估仓储
func NewTrainingAppService(
	exerciseRepo training.ExerciseRecordRepository,
	planRepo training.TrainingPlanRepository,
	assessmentRepo training.FitnessAssessmentRepository,
) *TrainingAppService {
	return &TrainingAppService{
		exerciseRepo:  exerciseRepo,
		planRepo:      planRepo,
		assessmentRepo: assessmentRepo,
	}
}

// CreateExerciseRecord 创建运动记录
// 校验运动类型、创建记录、检查是否打破个人纪录
// ctx: 上下文，userID: 用户ID，req: 创建运动记录请求DTO
func (s *TrainingAppService) CreateExerciseRecord(ctx context.Context, userID string, req *CreateExerciseRequest) (*ExerciseDTO, error) {
	// 校验运动类型
	exerciseType := training.ExerciseType(req.Type)
	if !isValidExerciseType(exerciseType) {
		return nil, appErrors.ErrInvalidExerciseType
	}

	// 创建运动记录
	record := training.NewExerciseRecord(userID, exerciseType)
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

	// 校验分数
	if err := record.ValidateScore(); err != nil {
		return nil, appErrors.ErrBadRequest.WithMessage(err.Error())
	}

	// 持久化记录
	if err := s.exerciseRepo.Create(ctx, record); err != nil {
		return nil, appErrors.ErrInternal.WithMessage("创建运动记录失败")
	}

	// 检查是否打破个人纪录
	isRecordBroken := false
	personalBest, err := s.exerciseRepo.GetPersonalBest(ctx, userID, exerciseType)
	if err == nil && personalBest != nil {
		if record.Count > personalBest.Count || record.Score > personalBest.Score {
			isRecordBroken = true
		}
	} else if err != nil {
		// 没有历史记录，视为首次打破纪录
		isRecordBroken = true
	}

	dto := toExerciseDTO(record)
	dto.IsRecordBroken = isRecordBroken
	// 根据评分计算积分（评分/10，最低1分）
	dto.PointsEarned = max(record.Score/10, 1)

	return dto, nil
}

// GetExerciseRecords 分页查询运动记录
// ctx: 上下文，userID: 用户ID，page: 页码，pageSize: 每页大小，exerciseType: 运动类型（可选过滤）
func (s *TrainingAppService) GetExerciseRecords(ctx context.Context, userID string, page, pageSize int64, exerciseType string) ([]*ExerciseDTO, *response.Pagination, error) {
	pagination := training.Pagination{
		Page:     int(page),
		PageSize: int(pageSize),
	}

	var result *training.PaginatedResult[training.ExerciseRecord]
	var err error

	// 根据是否有运动类型过滤选择不同的查询方法
	if exerciseType != "" {
		result, err = s.exerciseRepo.GetByUserIDAndType(ctx, userID, training.ExerciseType(exerciseType), pagination)
	} else {
		filter := training.ExerciseRecordFilter{}
		result, err = s.exerciseRepo.GetByUserID(ctx, userID, filter, pagination)
	}

	if err != nil {
		return nil, nil, appErrors.ErrInternal.WithMessage("查询运动记录失败")
	}

	// 转换为DTO
	dtos := make([]*ExerciseDTO, 0, len(result.Items))
	for _, record := range result.Items {
		dtos = append(dtos, toExerciseDTO(&record))
	}

	pag := &response.Pagination{
		Page:       int64(result.Page),
		PageSize:   int64(result.PageSize),
		Total:      result.Total,
		TotalPages: int64(result.TotalPages),
	}

	return dtos, pag, nil
}

// GetPersonalBest 获取用户所有运动类型的个人最佳记录
// ctx: 上下文，userID: 用户ID
func (s *TrainingAppService) GetPersonalBest(ctx context.Context, userID string) ([]*PersonalBestDTO, error) {
	var bests []*PersonalBestDTO

	// 遍历所有运动类型查询个人最佳
	exerciseTypes := []training.ExerciseType{
		training.ExerciseTypeJumpRope,
		training.ExerciseTypeJumpingJack,
		training.ExerciseTypeSquat,
		training.ExerciseTypeSitUp,
		training.ExerciseTypeHighKnee,
		training.ExerciseTypePushUp,
	}

	for _, et := range exerciseTypes {
		record, err := s.exerciseRepo.GetPersonalBest(ctx, userID, et)
		if err != nil || record == nil {
			continue
		}
		bests = append(bests, &PersonalBestDTO{
			Type:      string(record.Type),
			BestScore: record.Score,
			BestCount: record.Count,
			AchievedAt: record.CreatedAt,
		})
	}

	return bests, nil
}

// GetTodayPlan 获取今日训练计划
// ctx: 上下文，userID: 用户ID
func (s *TrainingAppService) GetTodayPlan(ctx context.Context, userID string) (*TrainingPlanDTO, error) {
	today := time.Now().Truncate(24 * time.Hour)

	plan, err := s.planRepo.GetByUserIDAndDate(ctx, userID, today)
	if err != nil {
		return nil, appErrors.ErrPlanNotFound
	}

	return toTrainingPlanDTO(plan), nil
}

// CompletePlan 完成训练计划
// ctx: 上下文，userID: 用户ID，planID: 计划ID
func (s *TrainingAppService) CompletePlan(ctx context.Context, userID string, planID string) error {
	plan, err := s.planRepo.GetByID(ctx, planID)
	if err != nil {
		return appErrors.ErrPlanNotFound
	}

	// 验证计划属于该用户
	if plan.UserID != userID {
		return appErrors.ErrForbidden
	}

	// 更新计划状态为已完成
	now := time.Now()
	plan.Status = training.PlanStatusCompleted
	plan.CompletedAt = &now

	if err := s.planRepo.Update(ctx, plan); err != nil {
		return appErrors.ErrInternal.WithMessage("更新训练计划失败")
	}

	return nil
}

// CreateFitnessAssessment 创建体能评估
// ctx: 上下文，userID: 用户ID，dto: 体能评估DTO
func (s *TrainingAppService) CreateFitnessAssessment(ctx context.Context, userID string, dto *FitnessAssessmentDTO) (*FitnessAssessmentDTO, error) {
	assessment := training.NewFitnessAssessment(userID)
	assessment.Endurance = dto.Endurance
	assessment.Agility = dto.Agility
	assessment.Strength = dto.Strength
	assessment.Speed = dto.Speed
	assessment.Coordination = dto.Coordination
	assessment.Balance = dto.Balance
	assessment.Flexibility = dto.Flexibility
	assessment.AssessedAt = dto.AssessedAt

	// 校验评估分数
	if err := assessment.ValidateScores(); err != nil {
		return nil, appErrors.ErrBadRequest.WithMessage(err.Error())
	}

	if err := s.assessmentRepo.Create(ctx, assessment); err != nil {
		return nil, appErrors.ErrInternal.WithMessage("创建体能评估失败")
	}

	return toFitnessAssessmentDTO(assessment), nil
}

// GetLatestAssessment 获取用户最新的体能评估
// ctx: 上下文，userID: 用户ID
func (s *TrainingAppService) GetLatestAssessment(ctx context.Context, userID string) (*FitnessAssessmentDTO, error) {
	assessment, err := s.assessmentRepo.GetLatestByUserID(ctx, userID)
	if err != nil {
		return nil, appErrors.ErrNotFound.WithMessage("体能评估不存在")
	}

	return toFitnessAssessmentDTO(assessment), nil
}

// GetWeeklyStats 获取用户本周运动统计
// ctx: 上下文，userID: 用户ID
func (s *TrainingAppService) GetWeeklyStats(ctx context.Context, userID string) (*WeeklyStatsDTO, error) {
	now := time.Now()
	// 计算本周开始时间（周一）
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	weekStart := now.AddDate(0, 0, -(weekday - 1)).Truncate(24 * time.Hour)

	filter := training.ExerciseRecordFilter{
		FromDate: &weekStart,
		ToDate:   &now,
	}

	result, err := s.exerciseRepo.GetByUserID(ctx, userID, filter, training.Pagination{Page: 1, PageSize: 100})
	if err != nil {
		return nil, appErrors.ErrInternal.WithMessage("查询运动记录失败")
	}

	// 计算统计数据
	stats := &WeeklyStatsDTO{}
	activeDaysMap := make(map[string]bool)
	dailyMap := make(map[string]*DailyStats)

	for _, record := range result.Items {
		stats.TotalExercises++
		stats.TotalDuration += record.DurationSeconds
		stats.TotalCount += record.Count
		stats.AverageScore += float64(record.Score)

		dateStr := record.CreatedAt.Format("2006-01-02")
		activeDaysMap[dateStr] = true

		if daily, ok := dailyMap[dateStr]; ok {
			daily.ExerciseCount++
			daily.Duration += record.DurationSeconds
			daily.Count += record.Count
			daily.Score += float64(record.Score)
		} else {
			dailyMap[dateStr] = &DailyStats{
				Date:          dateStr,
				ExerciseCount: 1,
				Duration:      record.DurationSeconds,
				Count:         record.Count,
				Score:         float64(record.Score),
			}
		}
	}

	if stats.TotalExercises > 0 {
		stats.AverageScore /= float64(stats.TotalExercises)
	}
	stats.ActiveDays = len(activeDaysMap)

	// 构建每日明细
	for _, daily := range dailyMap {
		if daily.ExerciseCount > 0 {
			daily.Score /= float64(daily.ExerciseCount)
		}
		stats.DailyBreakdown = append(stats.DailyBreakdown, *daily)
	}

	return stats, nil
}

// GetMonthlyStats 获取用户本月运动统计
// ctx: 上下文，userID: 用户ID
func (s *TrainingAppService) GetMonthlyStats(ctx context.Context, userID string) (*MonthlyStatsDTO, error) {
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	filter := training.ExerciseRecordFilter{
		FromDate: &monthStart,
		ToDate:   &now,
	}

	result, err := s.exerciseRepo.GetByUserID(ctx, userID, filter, training.Pagination{Page: 1, PageSize: 500})
	if err != nil {
		return nil, appErrors.ErrInternal.WithMessage("查询运动记录失败")
	}

	// 计算统计数据
	stats := &MonthlyStatsDTO{}
	activeDaysMap := make(map[string]bool)
	weeklyMap := make(map[string]*WeeklyBreakdown)

	for _, record := range result.Items {
		stats.TotalExercises++
		stats.TotalDuration += record.DurationSeconds
		stats.TotalCount += record.Count
		stats.AverageScore += float64(record.Score)

		dateStr := record.CreatedAt.Format("2006-01-02")
		activeDaysMap[dateStr] = true

		// 计算所属周
		weekday := int(record.CreatedAt.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		weekStart := record.CreatedAt.AddDate(0, 0, -(weekday - 1)).Format("2006-01-02")

		if weekly, ok := weeklyMap[weekStart]; ok {
			weekly.ExerciseCount++
			weekly.Duration += record.DurationSeconds
			weekly.Count += record.Count
			weekly.Score += float64(record.Score)
		} else {
			weeklyMap[weekStart] = &WeeklyBreakdown{
				WeekStart:     weekStart,
				ExerciseCount: 1,
				Duration:      record.DurationSeconds,
				Count:         record.Count,
				Score:         float64(record.Score),
			}
		}
	}

	if stats.TotalExercises > 0 {
		stats.AverageScore /= float64(stats.TotalExercises)
	}
	stats.ActiveDays = len(activeDaysMap)

	// 构建每周明细
	for _, weekly := range weeklyMap {
		if weekly.ExerciseCount > 0 {
			weekly.Score /= float64(weekly.ExerciseCount)
		}
		stats.WeeklyBreakdown = append(stats.WeeklyBreakdown, *weekly)
	}

	return stats, nil
}

// isValidExerciseType 校验运动类型是否合法
func isValidExerciseType(t training.ExerciseType) bool {
	switch t {
	case training.ExerciseTypeJumpRope, training.ExerciseTypeJumpingJack,
		training.ExerciseTypeSquat, training.ExerciseTypeSitUp,
		training.ExerciseTypeHighKnee, training.ExerciseTypePushUp:
		return true
	default:
		return false
	}
}

// toExerciseDTO 将运动记录领域模型转换为DTO
func toExerciseDTO(r *training.ExerciseRecord) *ExerciseDTO {
	return &ExerciseDTO{
		ID:              r.ID,
		UserID:          r.UserID,
		Type:            string(r.Type),
		DurationSeconds: r.DurationSeconds,
		Count:           r.Count,
		Score:           r.Score,
		RhythmScore:     r.RhythmScore,
		AmplitudeScore:  r.AmplitudeScore,
		SymmetryScore:   r.SymmetryScore,
		ContinuityScore: r.ContinuityScore,
		Corrections:     r.Corrections,
		IsOffline:       r.IsOffline,
		StartedAt:       r.StartedAt,
		CompletedAt:     r.CompletedAt,
		CreatedAt:       r.CreatedAt,
	}
}

// toTrainingPlanDTO 将训练计划领域模型转换为DTO
func toTrainingPlanDTO(p *training.TrainingPlan) *TrainingPlanDTO {
	dto := &TrainingPlanDTO{
		ID:             p.ID,
		Date:           p.Date,
		Status:         string(p.Status),
		TotalDuration:  p.TotalDuration,
		ActualDuration: p.ActualDuration,
		Items:          make([]ExerciseItemDTO, 0),
	}

	// 合并所有阶段的项目
	for _, item := range p.WarmupItems {
		dto.Items = append(dto.Items, toExerciseItemDTO(&item))
	}
	for _, item := range p.MainItems {
		dto.Items = append(dto.Items, toExerciseItemDTO(&item))
	}
	for _, item := range p.CooldownItems {
		dto.Items = append(dto.Items, toExerciseItemDTO(&item))
	}

	return dto
}

// toExerciseItemDTO 将运动项目领域模型转换为DTO
func toExerciseItemDTO(item *training.ExerciseItem) ExerciseItemDTO {
	return ExerciseItemDTO{
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

// toFitnessAssessmentDTO 将体能评估领域模型转换为DTO
func toFitnessAssessmentDTO(a *training.FitnessAssessment) *FitnessAssessmentDTO {
	return &FitnessAssessmentDTO{
		ID:           a.ID,
		UserID:       a.UserID,
		Endurance:    a.Endurance,
		Agility:      a.Agility,
		Strength:     a.Strength,
		Speed:        a.Speed,
		Coordination: a.Coordination,
		Balance:      a.Balance,
		Flexibility:  a.Flexibility,
		AssessedAt:   a.AssessedAt,
		CreatedAt:    a.CreatedAt,
	}
}
