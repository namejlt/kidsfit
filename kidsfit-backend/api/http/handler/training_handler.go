package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/kidsfit/api/api/http/middleware"
	"github.com/kidsfit/api/internal/application/training"
	appErrors "github.com/kidsfit/api/internal/pkg/errors"
	"github.com/kidsfit/api/internal/pkg/response"
)

// TrainingHandler 训练HTTP处理器，处理训练相关的HTTP请求
type TrainingHandler struct {
	trainingAppService *training.TrainingAppService
}

// NewTrainingHandler 创建训练HTTP处理器实例
// svc: 训练应用服务
func NewTrainingHandler(svc *training.TrainingAppService) *TrainingHandler {
	return &TrainingHandler{
		trainingAppService: svc,
	}
}

// CreateExerciseRecord 创建运动记录
// POST /api/v1/exercises
// 请求体包含运动类型、时长、次数、评分等信息
func (h *TrainingHandler) CreateExerciseRecord(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Error(c, appErrors.ErrUnauthorized)
		return
	}

	var req training.CreateExerciseRequest
	// 绑定并校验请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, appErrors.ErrBadRequest.WithMessage(err.Error()))
		return
	}

	// 调用应用服务创建运动记录
	result, err := h.trainingAppService.CreateExerciseRecord(c.Request.Context(), userID, &req)
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

// GetExerciseRecords 分页查询运动记录
// GET /api/v1/exercises
// 支持通过query参数过滤运动类型和分页
func (h *TrainingHandler) GetExerciseRecords(c *gin.Context) {
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

	// 获取运动类型过滤参数（可选）
	exerciseType := c.Query("type")

	// 调用应用服务查询运动记录
	result, pag, err := h.trainingAppService.GetExerciseRecords(c.Request.Context(), userID, page, pageSize, exerciseType)
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

// GetPersonalBest 获取用户所有运动类型的个人最佳记录
// GET /api/v1/exercises/personal-best
func (h *TrainingHandler) GetPersonalBest(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Error(c, appErrors.ErrUnauthorized)
		return
	}

	// 调用应用服务获取个人最佳记录
	result, err := h.trainingAppService.GetPersonalBest(c.Request.Context(), userID)
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

// GetTodayPlan 获取今日训练计划
// GET /api/v1/plans/today
func (h *TrainingHandler) GetTodayPlan(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Error(c, appErrors.ErrUnauthorized)
		return
	}

	// 调用应用服务获取今日训练计划
	result, err := h.trainingAppService.GetTodayPlan(c.Request.Context(), userID)
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

// CompletePlan 完成训练计划
// POST /api/v1/plans/:id/complete
// 路径参数id为训练计划ID
func (h *TrainingHandler) CompletePlan(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Error(c, appErrors.ErrUnauthorized)
		return
	}

	planID := c.Param("id")
	if planID == "" {
		response.Error(c, appErrors.ErrBadRequest.WithMessage("计划ID不能为空"))
		return
	}

	// 调用应用服务完成训练计划
	err := h.trainingAppService.CompletePlan(c.Request.Context(), userID, planID)
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

// CreateFitnessAssessment 创建体能评估
// POST /api/v1/assessments
// 请求体包含各维度评分和评估时间
func (h *TrainingHandler) CreateFitnessAssessment(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Error(c, appErrors.ErrUnauthorized)
		return
	}

	var dto training.FitnessAssessmentDTO
	// 绑定并校验请求参数
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, appErrors.ErrBadRequest.WithMessage(err.Error()))
		return
	}

	// 调用应用服务创建体能评估
	result, err := h.trainingAppService.CreateFitnessAssessment(c.Request.Context(), userID, &dto)
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

// GetLatestAssessment 获取用户最新的体能评估
// GET /api/v1/assessments/latest
func (h *TrainingHandler) GetLatestAssessment(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Error(c, appErrors.ErrUnauthorized)
		return
	}

	// 调用应用服务获取最新体能评估
	result, err := h.trainingAppService.GetLatestAssessment(c.Request.Context(), userID)
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

// GetWeeklyStats 获取用户本周运动统计
// GET /api/v1/stats/weekly
func (h *TrainingHandler) GetWeeklyStats(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Error(c, appErrors.ErrUnauthorized)
		return
	}

	// 调用应用服务获取周统计
	result, err := h.trainingAppService.GetWeeklyStats(c.Request.Context(), userID)
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

// GetMonthlyStats 获取用户本月运动统计
// GET /api/v1/stats/monthly
func (h *TrainingHandler) GetMonthlyStats(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Error(c, appErrors.ErrUnauthorized)
		return
	}

	// 调用应用服务获取月统计
	result, err := h.trainingAppService.GetMonthlyStats(c.Request.Context(), userID)
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
