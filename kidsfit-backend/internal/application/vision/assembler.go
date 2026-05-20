package vision

import (
	domain "github.com/kidsfit/api/internal/domain/vision"
)

// VisionRecordToDTO 将视力记录领域模型转换为视力记录DTO
func VisionRecordToDTO(record *domain.VisionRecord) *VisionRecordDTO {
	if record == nil {
		return nil
	}
	return &VisionRecordDTO{
		ID:               record.ID,
		UserID:           record.UserID,
		ChildID:          record.ChildID,
		Date:             record.Date,
		RightEye:         *EyeDataToDTO(record.RightEye),
		LeftEye:          *EyeDataToDTO(record.LeftEye),
		AxialLengthRight: record.AxialLengthRight,
		AxialLengthLeft:  record.AxialLengthLeft,
		HyperopiaReserve: record.HyperopiaReserve,
		Source:           string(record.Source),
		ImageURL:         record.ImageURL,
		CreatedAt:        record.CreatedAt,
		VisionStatus:     string(record.VisionStatus()),
		RefractiveStatus: record.RefractiveStatus(),
	}
}

// EyeDataToDTO 将眼睛数据值对象转换为眼睛数据DTO
func EyeDataToDTO(eyeData domain.EyeData) *EyeDataDTO {
	return &EyeDataDTO{
		SPH:  eyeData.SPH,
		CYL:  eyeData.CYL,
		AXIS: eyeData.AXIS,
		VA:   eyeData.VA,
	}
}

// DTOToEyeData 将眼睛数据DTO转换为眼睛数据值对象
func DTOToEyeData(dto *EyeDataDTO) domain.EyeData {
	if dto == nil {
		return domain.EyeData{}
	}
	return domain.EyeData{
		SPH:  dto.SPH,
		CYL:  dto.CYL,
		AXIS: dto.AXIS,
		VA:   dto.VA,
	}
}

// CreateVisionRecordRequestToRecord 将创建视力记录请求DTO转换为视力记录领域模型
func CreateVisionRecordRequestToRecord(req *CreateVisionRecordRequest, userID string) *domain.VisionRecord {
	if req == nil {
		return nil
	}
	record := domain.NewVisionRecord(userID, req.ChildID, req.Date, domain.VisionDataSource(req.Source))
	record.RightEye = DTOToEyeData(&req.RightEye)
	record.LeftEye = DTOToEyeData(&req.LeftEye)
	record.AxialLengthRight = req.AxialLengthRight
	record.AxialLengthLeft = req.AxialLengthLeft
	record.HyperopiaReserve = req.HyperopiaReserve
	record.ImageURL = req.ImageURL
	return record
}

// OutdoorActivityToDTO 将户外活动领域模型转换为户外活动DTO
func OutdoorActivityToDTO(activity *domain.OutdoorActivity) *OutdoorActivityDTO {
	if activity == nil {
		return nil
	}
	dto := &OutdoorActivityDTO{
		ID:             activity.ID,
		UserID:         activity.UserID,
		Date:           activity.Date,
		DurationMin:    activity.DurationMin,
		CreatedAt:      activity.CreatedAt,
		IsTargetMet:    activity.IsTargetMet(),
		TargetProgress: activity.TargetProgress(),
	}
	// 转换活动时段列表
	segments := make([]OutdoorSegmentDTO, 0, len(activity.Segments))
	for _, seg := range activity.Segments {
		segments = append(segments, OutdoorSegmentDTO{
			ID:          seg.ID,
			StartTime:   seg.StartTime,
			EndTime:     seg.EndTime,
			DurationMin: seg.DurationMin,
			Location:    seg.Location,
		})
	}
	dto.Segments = segments
	return dto
}

// EyeReminderToDTO 将护眼提醒领域模型转换为护眼提醒DTO
func EyeReminderToDTO(reminder *domain.EyeReminder) *EyeReminderDTO {
	if reminder == nil {
		return nil
	}
	return &EyeReminderDTO{
		ID:           reminder.ID,
		UserID:       reminder.UserID,
		Type:         string(reminder.Type),
		TriggeredAt:  reminder.TriggeredAt,
		Acknowledged: reminder.Acknowledged,
		CreatedAt:    reminder.CreatedAt,
	}
}
