package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/kidsfit/api/api/http/middleware"
	"github.com/kidsfit/api/internal/application/reward"
	appErrors "github.com/kidsfit/api/internal/pkg/errors"
	"github.com/kidsfit/api/internal/pkg/response"
)

// RewardHandler 激励HTTP处理器，处理徽章、积分、挑战和排行榜相关的HTTP请求
type RewardHandler struct {
	rewardAppService *reward.RewardAppService
}

// NewRewardHandler 创建激励HTTP处理器实例
// svc: 激励应用服务
func NewRewardHandler(svc *reward.RewardAppService) *RewardHandler {
	return &RewardHandler{
		rewardAppService: svc,
	}
}

// GetBadges 按类别查询徽章列表，并标记用户是否已获得
// GET /api/v1/badges
// 支持通过query参数category过滤徽章类别
func (h *RewardHandler) GetBadges(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Error(c, appErrors.ErrUnauthorized)
		return
	}

	// 获取徽章类别过滤参数（可选）
	category := c.Query("category")

	// 调用应用服务查询徽章列表
	result, err := h.rewardAppService.GetBadges(c.Request.Context(), userID, category)
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

// GetMyBadges 获取当前用户已获得的所有徽章
// GET /api/v1/badges/my
func (h *RewardHandler) GetMyBadges(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Error(c, appErrors.ErrUnauthorized)
		return
	}

	// 调用应用服务获取用户已获得的徽章
	result, err := h.rewardAppService.GetMyBadges(c.Request.Context(), userID)
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

// GetPoints 分页查询用户积分记录
// GET /api/v1/points
// 支持通过query参数分页
func (h *RewardHandler) GetPoints(c *gin.Context) {
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

	// 调用应用服务查询积分记录
	result, pag, err := h.rewardAppService.GetPoints(c.Request.Context(), userID, page, pageSize)
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

// GetPointsBalance 获取用户当前积分余额
// GET /api/v1/points/balance
func (h *RewardHandler) GetPointsBalance(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Error(c, appErrors.ErrUnauthorized)
		return
	}

	// 调用应用服务获取积分余额
	balance, err := h.rewardAppService.GetPointsBalance(c.Request.Context(), userID)
	if err != nil {
		if appErr, ok := err.(*appErrors.AppError); ok {
			response.Error(c, appErr)
			return
		}
		response.Error(c, appErrors.ErrInternal.WithMessage(err.Error()))
		return
	}

	response.Success(c, map[string]int64{"balance": balance})
}

// CreateChallenge 创建挑战
// POST /api/v1/challenges
// 请求体包含挑战类型、运动类型、目标值和可选的接受者ID
func (h *RewardHandler) CreateChallenge(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Error(c, appErrors.ErrUnauthorized)
		return
	}

	var req reward.CreateChallengeRequest
	// 绑定并校验请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, appErrors.ErrBadRequest.WithMessage(err.Error()))
		return
	}

	// 调用应用服务创建挑战
	result, err := h.rewardAppService.CreateChallenge(c.Request.Context(), userID, &req)
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

// GetChallenges 获取挑战列表
// GET /api/v1/challenges
// 查询当前用户发起或接受的挑战
func (h *RewardHandler) GetChallenges(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Error(c, appErrors.ErrUnauthorized)
		return
	}

	// 调用应用服务获取用户发起的挑战
	initiatedChallenges, err := h.rewardAppService.GetChallenges(c.Request.Context(), userID)
	if err != nil {
		if appErr, ok := err.(*appErrors.AppError); ok {
			response.Error(c, appErr)
			return
		}
		response.Error(c, appErrors.ErrInternal.WithMessage(err.Error()))
		return
	}

	response.Success(c, initiatedChallenges)
}

// AcceptChallenge 接受挑战
// POST /api/v1/challenges/:id/accept
// 路径参数id为挑战ID
func (h *RewardHandler) AcceptChallenge(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Error(c, appErrors.ErrUnauthorized)
		return
	}

	challengeID := c.Param("id")
	if challengeID == "" {
		response.Error(c, appErrors.ErrBadRequest.WithMessage("挑战ID不能为空"))
		return
	}

	// 调用应用服务接受挑战
	result, err := h.rewardAppService.AcceptChallenge(c.Request.Context(), userID, challengeID)
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

// SubmitChallengeScore 提交挑战成绩
// POST /api/v1/challenges/:id/submit
// 路径参数id为挑战ID，请求体包含成绩
func (h *RewardHandler) SubmitChallengeScore(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Error(c, appErrors.ErrUnauthorized)
		return
	}

	challengeID := c.Param("id")
	if challengeID == "" {
		response.Error(c, appErrors.ErrBadRequest.WithMessage("挑战ID不能为空"))
		return
	}

	var req struct {
		Score int `json:"score" binding:"required,min=0"`
	}
	// 绑定并校验请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, appErrors.ErrBadRequest.WithMessage(err.Error()))
		return
	}

	// 调用应用服务提交挑战成绩
	result, err := h.rewardAppService.SubmitChallengeScore(c.Request.Context(), userID, challengeID, req.Score)
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

// GetFamilyLeaderboard 获取家庭排行榜
// GET /api/v1/leaderboard/family
// 查询当前家长下所有家庭成员的排名
func (h *RewardHandler) GetFamilyLeaderboard(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Error(c, appErrors.ErrUnauthorized)
		return
	}

	// 调用应用服务获取家庭排行榜
	result, err := h.rewardAppService.GetFamilyLeaderboard(c.Request.Context(), userID)
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

// GetGlobalLeaderboard 获取全局排行榜
// GET /api/v1/leaderboard/global
// 通过query参数指定运动类型和获取数量
func (h *RewardHandler) GetGlobalLeaderboard(c *gin.Context) {
	exerciseType := c.Query("exercise_type")
	if exerciseType == "" {
		response.Error(c, appErrors.ErrBadRequest.WithMessage("exercise_type不能为空"))
		return
	}

	// 解析获取数量参数
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64)
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	// 调用应用服务获取全局排行榜
	result, err := h.rewardAppService.GetGlobalLeaderboard(c.Request.Context(), exerciseType, limit)
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
