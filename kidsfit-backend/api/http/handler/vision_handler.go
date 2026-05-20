package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/kidsfit/api/api/http/middleware"
	"github.com/kidsfit/api/internal/application/vision"
	appErrors "github.com/kidsfit/api/internal/pkg/errors"
	"github.com/kidsfit/api/internal/pkg/response"
)

// VisionHandler 视力HTTP处理器，处理视力相关的HTTP请求
type VisionHandler struct {
	visionAppService *vision.VisionAppService
}

// NewVisionHandler 创建视力HTTP处理器实例
// svc: 视力应用服务
func NewVisionHandler(svc *vision.VisionAppService) *VisionHandler {
	return &VisionHandler{
		visionAppService: svc,
	}
}

// CreateVisionRecord 创建视力记录
// POST /api/v1/vision-records
// 请求体包含儿童ID、检查日期、双眼数据等
func (h *VisionHandler) CreateVisionRecord(c *gin.Context) {
	var req vision.CreateVisionRecordRequest
	// 绑定并校验请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, appErrors.ErrBadRequest.WithMessage(err.Error()))
		return
	}

	// 调用应用服务创建视力记录
	result, err := h.visionAppService.CreateVisionRecord(c.Request.Context(), &req)
	if err != nil {
		if appErr, ok := err.(*appErrors.AppError); ok {
			response.Error(c, appErr)
			return
		}
		response.Error(c, appErrors.ErrInternal.WithMessage(err.Error()))
		return
	}

	response.Success(c, result)
}

// GetVisionRecords 分页查询儿童的视力记录
// GET /api/v1/vision-records
// 通过query参数指定child_id和分页信息
func (h *VisionHandler) GetVisionRecords(c *gin.Context) {
	childID := c.Query("child_id")
	if childID == "" {
		response.Error(c, appErrors.ErrBadRequest.WithMessage("child_id不能为空"))
		return
	}

	// 解析分页参数
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("page_size", "20"), 10, 64)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 调用应用服务查询视力记录
	result, pag, err := h.visionAppService.GetVisionRecords(c.Request.Context(), childID, page, pageSize)
	if err != nil {
		if appErr, ok := err.(*appErrors.AppError); ok {
			response.Error(c, appErr)
			return
		}
		response.Error(c, appErrors.ErrInternal.WithMessage(err.Error()))
		return
	}

	response.SuccessWithPage(c, result, pag)
}

// GetVisionTrend 获取儿童视力趋势数据
// GET /api/v1/vision/trend
// 通过query参数指定child_id
func (h *VisionHandler) GetVisionTrend(c *gin.Context) {
	childID := c.Query("child_id")
	if childID == "" {
		response.Error(c, appErrors.ErrBadRequest.WithMessage("child_id不能为空"))
		return
	}

	// 调用应用服务获取视力趋势
	result, err := h.visionAppService.GetVisionTrend(c.Request.Context(), childID)
	if err != nil {
		if appErr, ok := err.(*appErrors.AppError); ok {
			response.Error(c, appErr)
			return
		}
		response.Error(c, appErrors.ErrInternal.WithMessage(err.Error()))
		return
	}

	response.Success(c, result)
}

// GetTodayOutdoor 获取今日户外活动记录
// GET /api/v1/outdoor/today
func (h *VisionHandler) GetTodayOutdoor(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Error(c, appErrors.ErrUnauthorized)
		return
	}

	// 调用应用服务获取今日户外活动
	result, err := h.visionAppService.GetTodayOutdoor(c.Request.Context(), userID)
	if err != nil {
		if appErr, ok := err.(*appErrors.AppError); ok {
			response.Error(c, appErr)
			return
		}
		response.Error(c, appErrors.ErrInternal.WithMessage(err.Error()))
		return
	}

	response.Success(c, result)
}

// SyncOutdoorData 同步户外活动数据
// POST /api/v1/outdoor/sync
// 请求体包含新增的户外时长（分钟）
func (h *VisionHandler) SyncOutdoorData(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Error(c, appErrors.ErrUnauthorized)
		return
	}

	var req struct {
		DurationMin int `json:"duration_min" binding:"required,min=1"`
	}
	// 绑定并校验请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, appErrors.ErrBadRequest.WithMessage(err.Error()))
		return
	}

	// 调用应用服务同步户外活动数据
	result, err := h.visionAppService.SyncOutdoorData(c.Request.Context(), userID, req.DurationMin)
	if err != nil {
		if appErr, ok := err.(*appErrors.AppError); ok {
			response.Error(c, appErr)
			return
		}
		response.Error(c, appErrors.ErrInternal.WithMessage(err.Error()))
		return
	}

	response.Success(c, result)
}

// GetReminders 分页查询护眼提醒
// GET /api/v1/reminders
// 支持通过query参数分页
func (h *VisionHandler) GetReminders(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Error(c, appErrors.ErrUnauthorized)
		return
	}

	// 解析分页参数
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("page_size", "20"), 10, 64)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 调用应用服务查询护眼提醒
	result, err := h.visionAppService.GetReminders(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		if appErr, ok := err.(*appErrors.AppError); ok {
			response.Error(c, appErr)
			return
		}
		response.Error(c, appErrors.ErrInternal.WithMessage(err.Error()))
		return
	}

	response.Success(c, result)
}

// OCRVisionRecord 通过OCR识别验光单图片创建视力记录
// POST /api/v1/vision-records/ocr
// 接收 multipart/form-data 上传的验光单图片
func (h *VisionHandler) OCRVisionRecord(c *gin.Context) {
	// 获取上传的图片文件
	file, err := c.FormFile("image")
	if err != nil {
		response.Error(c, appErrors.ErrBadRequest.WithMessage("请上传验光单图片"))
		return
	}

	// 校验文件大小（最大10MB）
	if file.Size > 10*1024*1024 {
		response.Error(c, appErrors.ErrBadRequest.WithMessage("图片大小不能超过10MB"))
		return
	}

	// TODO: 接入OCR服务，识别验光单图片中的视力数据
	// 当前返回占位响应
	result := map[string]interface{}{
		"message": "OCR识别功能开发中",
		"file":    file.Filename,
		"size":    file.Size,
	}

	response.Success(c, result)
}

// AckReminder 确认护眼提醒
// POST /api/v1/reminders/:id/ack
// 路径参数id为提醒ID
func (h *VisionHandler) AckReminder(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Error(c, appErrors.ErrUnauthorized)
		return
	}

	reminderID := c.Param("id")
	if reminderID == "" {
		response.Error(c, appErrors.ErrBadRequest.WithMessage("提醒ID不能为空"))
		return
	}

	// 调用应用服务确认提醒
	err := h.visionAppService.AckReminder(c.Request.Context(), userID, reminderID)
	if err != nil {
		if appErr, ok := err.(*appErrors.AppError); ok {
			response.Error(c, appErr)
			return
		}
		response.Error(c, appErrors.ErrInternal.WithMessage(err.Error()))
		return
	}

	response.Success(c, nil)
}
