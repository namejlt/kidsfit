package vision

import (
	"math"
	"time"

	"github.com/google/uuid"
)

// VisionDataSource 视力数据来源枚举
type VisionDataSource string

const (
	// VisionDataSourceOCR OCR识别来源
	VisionDataSourceOCR VisionDataSource = "ocr"
	// VisionDataSourceManual 手动录入来源
	VisionDataSourceManual VisionDataSource = "manual"
)

// ReminderType 提醒类型枚举
type ReminderType string

const (
	// ReminderType202020 20-20-20护眼法则提醒
	ReminderType202020 ReminderType = "20_20_20"
	// ReminderTypeOutdoor 户外活动提醒
	ReminderTypeOutdoor ReminderType = "outdoor"
	// ReminderTypeBreak 休息提醒
	ReminderTypeBreak ReminderType = "break"
)

// VisionStatus 视力状态枚举
type VisionStatus string

const (
	// VisionStatusGood 视力良好
	VisionStatusGood VisionStatus = "good"
	// VisionStatusMedium 视力一般
	VisionStatusMedium VisionStatus = "medium"
	// VisionStatusConcern 视力需关注
	VisionStatusConcern VisionStatus = "concern"
)

// EyeData 眼睛数据值对象，包含球镜、柱镜、轴位和矫正视力
type EyeData struct {
	SPH  float64 `json:"sph" db:"sph"`
	CYL  float64 `json:"cyl" db:"cyl"`
	AXIS float64 `json:"axis" db:"axis"`
	VA   float64 `json:"va" db:"va"`
}

// IsMyopic 判断是否近视（球镜度数<-0.50D）
func (ed EyeData) IsMyopic() bool {
	return ed.SPH < -0.50
}

// IsHyperopic 判断是否远视（球镜度数>+0.50D）
func (ed EyeData) IsHyperopic() bool {
	return ed.SPH > 0.50
}

// HasAstigmatism 判断是否有散光（柱镜绝对值>0.50D）
func (ed EyeData) HasAstigmatism() bool {
	return math.Abs(ed.CYL) > 0.50
}

// SphericalEquivalent 计算等效球镜度数（球镜+柱镜/2）
func (ed EyeData) SphericalEquivalent() float64 {
	return ed.SPH + ed.CYL/2
}

// VisionRecord 视力记录领域模型
type VisionRecord struct {
	ID               string           `json:"id" db:"id"`
	UserID           string           `json:"user_id" db:"user_id"`
	ChildID          string           `json:"child_id" db:"child_id"`
	Date             time.Time        `json:"date" db:"date"`
	RightEye         EyeData          `json:"right_eye" db:"right_eye"`
	LeftEye          EyeData          `json:"left_eye" db:"left_eye"`
	AxialLengthRight *float64         `json:"axial_length_right,omitempty" db:"axial_length_right"`
	AxialLengthLeft  *float64         `json:"axial_length_left,omitempty" db:"axial_length_left"`
	HyperopiaReserve *float64         `json:"hyperopia_reserve,omitempty" db:"hyperopia_reserve"`
	Source           VisionDataSource `json:"source" db:"source"`
	ImageURL         *string          `json:"image_url,omitempty" db:"image_url"`
	CreatedAt        time.Time        `json:"created_at" db:"created_at"`
}

// NewVisionRecord 创建视力记录实例
func NewVisionRecord(userID, childID string, date time.Time, source VisionDataSource) *VisionRecord {
	return &VisionRecord{
		ID:        uuid.New().String(),
		UserID:    userID,
		ChildID:   childID,
		Date:      date,
		Source:    source,
		CreatedAt: time.Now(),
	}
}

// AverageSPH 计算双眼平均球镜度数
func (vr *VisionRecord) AverageSPH() float64 {
	return (vr.RightEye.SPH + vr.LeftEye.SPH) / 2
}

// AverageVA 计算双眼平均矫正视力
func (vr *VisionRecord) AverageVA() float64 {
	return (vr.RightEye.VA + vr.LeftEye.VA) / 2
}

// RefractiveStatus 根据等效球镜判断屈光状态
func (vr *VisionRecord) RefractiveStatus() string {
	avgSE := (vr.RightEye.SphericalEquivalent() + vr.LeftEye.SphericalEquivalent()) / 2
	switch {
	case avgSE < -0.50:
		return "近视"
	case avgSE > 0.50:
		return "远视"
	default:
		return "正视"
	}
}

// VisionStatus 根据屈光数据评估视力状态
func (vr *VisionRecord) VisionStatus() VisionStatus {
	avgSE := (vr.RightEye.SphericalEquivalent() + vr.LeftEye.SphericalEquivalent()) / 2
	switch {
	case avgSE >= -0.50 && avgSE <= 0.50:
		return VisionStatusGood
	case avgSE > 0.50 && avgSE <= 1.50, avgSE < -0.50 && avgSE >= -3.00:
		return VisionStatusMedium
	default:
		return VisionStatusConcern
	}
}

// OutdoorActivity 户外活动领域模型
type OutdoorActivity struct {
	ID          string           `json:"id" db:"id"`
	UserID      string           `json:"user_id" db:"user_id"`
	Date        time.Time        `json:"date" db:"date"`
	DurationMin int              `json:"duration_min" db:"duration_min"`
	Segments    []OutdoorSegment `json:"segments" db:"segments"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at"`
}

// NewOutdoorActivity 创建户外活动实例
func NewOutdoorActivity(userID string, date time.Time) *OutdoorActivity {
	return &OutdoorActivity{
		ID:        uuid.New().String(),
		UserID:    userID,
		Date:      date,
		Segments:  []OutdoorSegment{},
		CreatedAt: time.Now(),
	}
}

// IsTargetMet 判断是否达到每日120分钟户外活动目标
func (oa *OutdoorActivity) IsTargetMet() bool {
	return oa.DurationMin >= 120
}

// TargetProgress 计算户外活动目标完成进度（0-100百分比）
func (oa *OutdoorActivity) TargetProgress() float64 {
	if oa.DurationMin >= 120 {
		return 100.0
	}
	return float64(oa.DurationMin) / 120.0 * 100.0
}

// OutdoorSegment 户外活动时段值对象
type OutdoorSegment struct {
	ID          string    `json:"id" db:"id"`
	ActivityID  string    `json:"activity_id" db:"activity_id"`
	StartTime   time.Time `json:"start_time" db:"start_time"`
	EndTime     time.Time `json:"end_time" db:"end_time"`
	DurationMin int       `json:"duration_min" db:"duration_min"`
	Location    string    `json:"location" db:"location"`
}

// NewOutdoorSegment 创建户外活动时段实例
func NewOutdoorSegment(activityID string, startTime, endTime time.Time, location string) *OutdoorSegment {
	duration := int(endTime.Sub(startTime).Minutes())
	return &OutdoorSegment{
		ID:          uuid.New().String(),
		ActivityID:  activityID,
		StartTime:   startTime,
		EndTime:     endTime,
		DurationMin: duration,
		Location:    location,
	}
}

// EyeReminder 护眼提醒领域模型
type EyeReminder struct {
	ID           string       `json:"id" db:"id"`
	UserID       string       `json:"user_id" db:"user_id"`
	Type         ReminderType `json:"type" db:"type"`
	TriggeredAt  time.Time    `json:"triggered_at" db:"triggered_at"`
	Acknowledged bool         `json:"acknowledged" db:"acknowledged"`
	CreatedAt    time.Time    `json:"created_at" db:"created_at"`
}

// NewEyeReminder 创建护眼提醒实例
func NewEyeReminder(userID string, reminderType ReminderType) *EyeReminder {
	return &EyeReminder{
		ID:          uuid.New().String(),
		UserID:      userID,
		Type:        reminderType,
		TriggeredAt: time.Now(),
		CreatedAt:   time.Now(),
	}
}
