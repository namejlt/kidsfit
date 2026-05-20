package vision

import "time"

// CreateVisionRecordRequest 创建视力记录请求DTO
type CreateVisionRecordRequest struct {
	// ChildID 儿童ID
	ChildID string `json:"child_id" binding:"required"`
	// Date 检查日期
	Date time.Time `json:"date" binding:"required"`
	// RightEye 右眼数据
	RightEye EyeDataDTO `json:"right_eye" binding:"required"`
	// LeftEye 左眼数据
	LeftEye EyeDataDTO `json:"left_eye" binding:"required"`
	// AxialLengthRight 右眼眼轴长度
	AxialLengthRight *float64 `json:"axial_length_right,omitempty"`
	// AxialLengthLeft 左眼眼轴长度
	AxialLengthLeft *float64 `json:"axial_length_left,omitempty"`
	// HyperopiaReserve 远视储备
	HyperopiaReserve *float64 `json:"hyperopia_reserve,omitempty"`
	// Source 数据来源（ocr/manual）
	Source string `json:"source" binding:"required"`
	// ImageURL 检查单图片URL
	ImageURL *string `json:"image_url,omitempty"`
}

// EyeDataDTO 眼睛数据DTO
type EyeDataDTO struct {
	// SPH 球镜度数
	SPH float64 `json:"sph"`
	// CYL 柱镜度数
	CYL float64 `json:"cyl"`
	// AXIS 轴位
	AXIS float64 `json:"axis"`
	// VA 矫正视力
	VA float64 `json:"va"`
}

// VisionRecordDTO 视力记录响应DTO
type VisionRecordDTO struct {
	// ID 记录ID
	ID string `json:"id"`
	// UserID 用户ID（家长）
	UserID string `json:"user_id"`
	// ChildID 儿童ID
	ChildID string `json:"child_id"`
	// Date 检查日期
	Date time.Time `json:"date"`
	// RightEye 右眼数据
	RightEye EyeDataDTO `json:"right_eye"`
	// LeftEye 左眼数据
	LeftEye EyeDataDTO `json:"left_eye"`
	// AxialLengthRight 右眼眼轴长度
	AxialLengthRight *float64 `json:"axial_length_right,omitempty"`
	// AxialLengthLeft 左眼眼轴长度
	AxialLengthLeft *float64 `json:"axial_length_left,omitempty"`
	// HyperopiaReserve 远视储备
	HyperopiaReserve *float64 `json:"hyperopia_reserve,omitempty"`
	// Source 数据来源
	Source string `json:"source"`
	// ImageURL 检查单图片URL
	ImageURL *string `json:"image_url,omitempty"`
	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`
	// VisionStatus 视力状态（good/medium/concern）
	VisionStatus string `json:"vision_status"`
	// RefractiveStatus 屈光状态（正视/近视/远视）
	RefractiveStatus string `json:"refractive_status"`
}

// OutdoorActivityDTO 户外活动响应DTO
type OutdoorActivityDTO struct {
	// ID 活动ID
	ID string `json:"id"`
	// UserID 用户ID
	UserID string `json:"user_id"`
	// Date 日期
	Date time.Time `json:"date"`
	// DurationMin 户外时长（分钟）
	DurationMin int `json:"duration_min"`
	// Segments 活动时段列表
	Segments []OutdoorSegmentDTO `json:"segments"`
	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`
	// IsTargetMet 是否达到每日目标
	IsTargetMet bool `json:"is_target_met"`
	// TargetProgress 目标完成进度（0-100）
	TargetProgress float64 `json:"target_progress"`
}

// OutdoorSegmentDTO 户外活动时段DTO
type OutdoorSegmentDTO struct {
	// ID 时段ID
	ID string `json:"id"`
	// StartTime 开始时间
	StartTime time.Time `json:"start_time"`
	// EndTime 结束时间
	EndTime time.Time `json:"end_time"`
	// DurationMin 时长（分钟）
	DurationMin int `json:"duration_min"`
	// Location 位置
	Location string `json:"location"`
}

// EyeReminderDTO 护眼提醒响应DTO
type EyeReminderDTO struct {
	// ID 提醒ID
	ID string `json:"id"`
	// UserID 用户ID
	UserID string `json:"user_id"`
	// Type 提醒类型
	Type string `json:"type"`
	// TriggeredAt 触发时间
	TriggeredAt time.Time `json:"triggered_at"`
	// Acknowledged 是否已确认
	Acknowledged bool `json:"acknowledged"`
	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`
}

// VisionTrendDTO 视力趋势响应DTO
type VisionTrendDTO struct {
	// ChildID 儿童ID
	ChildID string `json:"child_id"`
	// DataPoints 趋势数据点列表
	DataPoints []VisionTrendPoint `json:"data_points"`
}

// VisionTrendPoint 视力趋势数据点
type VisionTrendPoint struct {
	// Date 日期
	Date string `json:"date"`
	// RightEyeSPH 右眼球镜度数
	RightEyeSPH float64 `json:"right_eye_sph"`
	// LeftEyeSPH 左眼球镜度数
	LeftEyeSPH float64 `json:"left_eye_sph"`
	// RightEyeVA 右眼矫正视力
	RightEyeVA float64 `json:"right_eye_va"`
	// LeftEyeVA 左眼矫正视力
	LeftEyeVA float64 `json:"left_eye_va"`
	// AxialLengthRight 右眼眼轴长度
	AxialLengthRight *float64 `json:"axial_length_right,omitempty"`
	// AxialLengthLeft 左眼眼轴长度
	AxialLengthLeft *float64 `json:"axial_length_left,omitempty"`
}
