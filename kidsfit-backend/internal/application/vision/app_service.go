package vision

import (
	"context"
	"time"

	"github.com/kidsfit/api/internal/domain/vision"
	appErrors "github.com/kidsfit/api/internal/pkg/errors"
	"github.com/kidsfit/api/internal/pkg/response"
)

// VisionAppService 视力应用服务，负责视力相关的业务逻辑编排
type VisionAppService struct {
	visionRepo  vision.VisionRecordRepository
	outdoorRepo vision.OutdoorActivityRepository
	reminderRepo vision.EyeReminderRepository
}

// NewVisionAppService 创建视力应用服务实例
// visionRepo: 视力记录仓储，outdoorRepo: 户外活动仓储，reminderRepo: 护眼提醒仓储
func NewVisionAppService(
	visionRepo vision.VisionRecordRepository,
	outdoorRepo vision.OutdoorActivityRepository,
	reminderRepo vision.EyeReminderRepository,
) *VisionAppService {
	return &VisionAppService{
		visionRepo:  visionRepo,
		outdoorRepo: outdoorRepo,
		reminderRepo: reminderRepo,
	}
}

// CreateVisionRecord 创建视力记录
// 将请求DTO转换为领域模型并持久化
// ctx: 上下文，req: 创建视力记录请求DTO
func (s *VisionAppService) CreateVisionRecord(ctx context.Context, req *CreateVisionRecordRequest) (*VisionRecordDTO, error) {
	// 创建视力记录领域模型
	record := vision.NewVisionRecord("", req.ChildID, req.Date, vision.VisionDataSource(req.Source))

	// 填充眼睛数据
	record.RightEye = vision.EyeData{
		SPH:  req.RightEye.SPH,
		CYL:  req.RightEye.CYL,
		AXIS: req.RightEye.AXIS,
		VA:   req.RightEye.VA,
	}
	record.LeftEye = vision.EyeData{
		SPH:  req.LeftEye.SPH,
		CYL:  req.LeftEye.CYL,
		AXIS: req.LeftEye.AXIS,
		VA:   req.LeftEye.VA,
	}
	record.AxialLengthRight = req.AxialLengthRight
	record.AxialLengthLeft = req.AxialLengthLeft
	record.HyperopiaReserve = req.HyperopiaReserve
	record.ImageURL = req.ImageURL

	// 持久化记录
	if err := s.visionRepo.Create(ctx, record); err != nil {
		return nil, appErrors.ErrInternal.WithMessage("创建视力记录失败")
	}

	return toVisionRecordDTO(record), nil
}

// GetVisionRecords 分页查询儿童的视力记录
// ctx: 上下文，childID: 儿童ID，page: 页码，pageSize: 每页大小
func (s *VisionAppService) GetVisionRecords(ctx context.Context, childID string, page, pageSize int64) ([]*VisionRecordDTO, *response.Pagination, error) {
	pagination := vision.Pagination{
		Page:     int(page),
		PageSize: int(pageSize),
	}

	result, err := s.visionRepo.GetByChildID(ctx, childID, pagination)
	if err != nil {
		return nil, nil, appErrors.ErrInternal.WithMessage("查询视力记录失败")
	}

	// 转换为DTO
	dtos := make([]*VisionRecordDTO, 0, len(result.Items))
	for _, record := range result.Items {
		dtos = append(dtos, toVisionRecordDTO(&record))
	}

	pag := &response.Pagination{
		Page:       int64(result.Page),
		PageSize:   int64(result.PageSize),
		Total:      result.Total,
		TotalPages: int64(result.TotalPages),
	}

	return dtos, pag, nil
}

// GetVisionTrend 获取儿童视力趋势数据
// 查询最近6个月的视力记录，生成趋势数据点
// ctx: 上下文，childID: 儿童ID
func (s *VisionAppService) GetVisionTrend(ctx context.Context, childID string) (*VisionTrendDTO, error) {
	// 查询最近6个月的数据
	now := time.Now()
	sixMonthsAgo := now.AddDate(0, -6, 0)

	records, err := s.visionRepo.GetByChildIDAndDateRange(ctx, childID, sixMonthsAgo, now)
	if err != nil {
		return nil, appErrors.ErrInternal.WithMessage("查询视力趋势失败")
	}

	// 构建趋势数据
	trend := &VisionTrendDTO{
		ChildID:    childID,
		DataPoints: make([]VisionTrendPoint, 0, len(records)),
	}

	for _, record := range records {
		point := VisionTrendPoint{
			Date:             record.Date.Format("2006-01-02"),
			RightEyeSPH:      record.RightEye.SPH,
			LeftEyeSPH:       record.LeftEye.SPH,
			RightEyeVA:       record.RightEye.VA,
			LeftEyeVA:        record.LeftEye.VA,
			AxialLengthRight: record.AxialLengthRight,
			AxialLengthLeft:  record.AxialLengthLeft,
		}
		trend.DataPoints = append(trend.DataPoints, point)
	}

	return trend, nil
}

// GetTodayOutdoor 获取今日户外活动记录
// ctx: 上下文，userID: 用户ID
func (s *VisionAppService) GetTodayOutdoor(ctx context.Context, userID string) (*OutdoorActivityDTO, error) {
	today := time.Now().Truncate(24 * time.Hour)

	activity, err := s.outdoorRepo.GetByUserIDAndDate(ctx, userID, today)
	if err != nil {
		// 没有今日记录，返回空DTO
		return &OutdoorActivityDTO{
			UserID:         userID,
			Date:           today,
			DurationMin:    0,
			Segments:       []OutdoorSegmentDTO{},
			IsTargetMet:    false,
			TargetProgress: 0,
		}, nil
	}

	return toOutdoorActivityDTO(activity), nil
}

// SyncOutdoorData 同步户外活动数据，累加今日户外时长
// ctx: 上下文，userID: 用户ID，durationMin: 新增户外时长（分钟）
func (s *VisionAppService) SyncOutdoorData(ctx context.Context, userID string, durationMin int) (*OutdoorActivityDTO, error) {
	today := time.Now().Truncate(24 * time.Hour)

	// 查找今日已有记录
	activity, err := s.outdoorRepo.GetByUserIDAndDate(ctx, userID, today)
	if err != nil {
		// 创建新记录
		activity = vision.NewOutdoorActivity(userID, today)
		activity.DurationMin = durationMin

		if err := s.outdoorRepo.Create(ctx, activity); err != nil {
			return nil, appErrors.ErrInternal.WithMessage("创建户外活动记录失败")
		}
	} else {
		// 更新已有记录，累加时长
		activity.DurationMin += durationMin

		if err := s.outdoorRepo.Update(ctx, activity); err != nil {
			return nil, appErrors.ErrInternal.WithMessage("更新户外活动记录失败")
		}
	}

	return toOutdoorActivityDTO(activity), nil
}

// GetReminders 分页查询护眼提醒
// ctx: 上下文，userID: 用户ID，page: 页码，pageSize: 每页大小
func (s *VisionAppService) GetReminders(ctx context.Context, userID string, page, pageSize int64) ([]*EyeReminderDTO, error) {
	pagination := vision.Pagination{
		Page:     int(page),
		PageSize: int(pageSize),
	}

	result, err := s.reminderRepo.GetByUserID(ctx, userID, pagination)
	if err != nil {
		return nil, appErrors.ErrInternal.WithMessage("查询护眼提醒失败")
	}

	// 转换为DTO
	dtos := make([]*EyeReminderDTO, 0, len(result.Items))
	for _, reminder := range result.Items {
		dtos = append(dtos, toEyeReminderDTO(&reminder))
	}

	return dtos, nil
}

// AckReminder 确认护眼提醒
// ctx: 上下文，userID: 用户ID，reminderID: 提醒ID
func (s *VisionAppService) AckReminder(ctx context.Context, userID string, reminderID string) error {
	if err := s.reminderRepo.UpdateAcknowledged(ctx, reminderID, true); err != nil {
		return appErrors.ErrInternal.WithMessage("确认提醒失败")
	}
	return nil
}

// toVisionRecordDTO 将视力记录领域模型转换为DTO
func toVisionRecordDTO(r *vision.VisionRecord) *VisionRecordDTO {
	return &VisionRecordDTO{
		ID:               r.ID,
		UserID:           r.UserID,
		ChildID:          r.ChildID,
		Date:             r.Date,
		RightEye:         toEyeDataDTO(r.RightEye),
		LeftEye:          toEyeDataDTO(r.LeftEye),
		AxialLengthRight: r.AxialLengthRight,
		AxialLengthLeft:  r.AxialLengthLeft,
		HyperopiaReserve: r.HyperopiaReserve,
		Source:           string(r.Source),
		ImageURL:         r.ImageURL,
		CreatedAt:        r.CreatedAt,
		VisionStatus:     string(r.VisionStatus()),
		RefractiveStatus: r.RefractiveStatus(),
	}
}

// toEyeDataDTO 将眼睛数据值对象转换为DTO
func toEyeDataDTO(ed vision.EyeData) EyeDataDTO {
	return EyeDataDTO{
		SPH:  ed.SPH,
		CYL:  ed.CYL,
		AXIS: ed.AXIS,
		VA:   ed.VA,
	}
}

// toOutdoorActivityDTO 将户外活动领域模型转换为DTO
func toOutdoorActivityDTO(oa *vision.OutdoorActivity) *OutdoorActivityDTO {
	segments := make([]OutdoorSegmentDTO, 0, len(oa.Segments))
	for _, seg := range oa.Segments {
		segments = append(segments, OutdoorSegmentDTO{
			ID:          seg.ID,
			StartTime:   seg.StartTime,
			EndTime:     seg.EndTime,
			DurationMin: seg.DurationMin,
			Location:    seg.Location,
		})
	}

	return &OutdoorActivityDTO{
		ID:             oa.ID,
		UserID:         oa.UserID,
		Date:           oa.Date,
		DurationMin:    oa.DurationMin,
		Segments:       segments,
		CreatedAt:      oa.CreatedAt,
		IsTargetMet:    oa.IsTargetMet(),
		TargetProgress: oa.TargetProgress(),
	}
}

// toEyeReminderDTO 将护眼提醒领域模型转换为DTO
func toEyeReminderDTO(r *vision.EyeReminder) *EyeReminderDTO {
	return &EyeReminderDTO{
		ID:           r.ID,
		UserID:       r.UserID,
		Type:         string(r.Type),
		TriggeredAt:  r.TriggeredAt,
		Acknowledged: r.Acknowledged,
		CreatedAt:    r.CreatedAt,
	}
}
