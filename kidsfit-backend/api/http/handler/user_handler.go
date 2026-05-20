package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kidsfit/api/api/http/middleware"
	"github.com/kidsfit/api/internal/application/user"
	appErrors "github.com/kidsfit/api/internal/pkg/errors"
	"github.com/kidsfit/api/internal/pkg/response"
)

// UserHandler 用户HTTP处理器，处理用户相关的HTTP请求
type UserHandler struct {
	userAppService *user.UserAppService
}

// NewUserHandler 创建用户HTTP处理器实例
// svc: 用户应用服务
func NewUserHandler(svc *user.UserAppService) *UserHandler {
	return &UserHandler{
		userAppService: svc,
	}
}

// Register 注册家长账号
// POST /api/v1/auth/register
// 请求体包含手机号、密码和昵称
func (h *UserHandler) Register(c *gin.Context) {
	var req user.RegisterRequest
	// 绑定并校验请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, appErrors.ErrBadRequest.WithMessage(err.Error()))
		return
	}

	// 调用应用服务执行注册
	result, err := h.userAppService.Register(c.Request.Context(), &req)
	if err != nil {
		if appErr, ok := err.(*appErrors.AppError); ok {
			response.Error(c, appErr)
			return
		}
		response.Error(c, appErrors.ErrInternal.WithMessage(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.Response{
		Code:    0,
		Message: "success",
		Data:    result,
	})
}

// Login 用户登录
// POST /api/v1/auth/login
// 请求体包含手机号和密码
func (h *UserHandler) Login(c *gin.Context) {
	var req user.LoginRequest
	// 绑定并校验请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, appErrors.ErrBadRequest.WithMessage(err.Error()))
		return
	}

	// 调用应用服务执行登录
	result, err := h.userAppService.Login(c.Request.Context(), &req)
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

// RefreshToken 刷新访问令牌
// POST /api/v1/auth/refresh
// 请求体包含refresh_token
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	// 绑定并校验请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, appErrors.ErrBadRequest.WithMessage(err.Error()))
		return
	}

	// 调用应用服务刷新令牌
	result, err := h.userAppService.RefreshToken(c.Request.Context(), req.RefreshToken)
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

// Logout 用户登出
// POST /api/v1/auth/logout
// 当前为空操作，后续可加入令牌黑名单逻辑
func (h *UserHandler) Logout(c *gin.Context) {
	// TODO: 将当前令牌加入黑名单
	response.Success(c, nil)
}

// GetCurrentUser 获取当前登录用户信息
// GET /api/v1/users/me
// 从认证中间件注入的上下文中获取用户ID
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Error(c, appErrors.ErrUnauthorized)
		return
	}

	// 调用应用服务获取用户信息
	result, err := h.userAppService.GetCurrentUser(c.Request.Context(), userID)
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

// UpdateUser 更新当前用户信息
// PUT /api/v1/users/me
// 请求体包含待更新的昵称和头像
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Error(c, appErrors.ErrUnauthorized)
		return
	}

	var req user.UpdateUserRequest
	// 绑定并校验请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, appErrors.ErrBadRequest.WithMessage(err.Error()))
		return
	}

	// 调用应用服务更新用户信息
	result, err := h.userAppService.UpdateUser(c.Request.Context(), userID, &req)
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

// AddChild 添加儿童账号
// POST /api/v1/children
// 仅家长可操作，请求体包含儿童昵称、年龄和头像
func (h *UserHandler) AddChild(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Error(c, appErrors.ErrUnauthorized)
		return
	}

	var req user.AddChildRequest
	// 绑定并校验请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, appErrors.ErrBadRequest.WithMessage(err.Error()))
		return
	}

	// 调用应用服务添加儿童
	result, err := h.userAppService.AddChild(c.Request.Context(), userID, &req)
	if err != nil {
		if appErr, ok := err.(*appErrors.AppError); ok {
			response.Error(c, appErr)
			return
		}
		response.Error(c, appErrors.ErrInternal.WithMessage(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.Response{
		Code:    0,
		Message: "success",
		Data:    result,
	})
}

// GetChildren 获取当前家长下的所有儿童列表
// GET /api/v1/children
// 仅家长可操作
func (h *UserHandler) GetChildren(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Error(c, appErrors.ErrUnauthorized)
		return
	}

	// 调用应用服务获取儿童列表
	result, err := h.userAppService.GetChildren(c.Request.Context(), userID)
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

// GetParentSettings 获取家长设置
// GET /api/v1/settings
// 仅家长可操作
func (h *UserHandler) GetParentSettings(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Error(c, appErrors.ErrUnauthorized)
		return
	}

	// 调用应用服务获取家长设置
	result, err := h.userAppService.GetParentSettings(c.Request.Context(), userID)
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

// UploadAvatar 上传用户头像
// POST /api/v1/users/avatar
// 接收 multipart/form-data 上传的头像图片
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Error(c, appErrors.ErrUnauthorized)
		return
	}

	// 获取上传的头像文件
	file, err := c.FormFile("avatar")
	if err != nil {
		response.Error(c, appErrors.ErrBadRequest.WithMessage("请上传头像图片"))
		return
	}

	// 校验文件大小（最大5MB）
	if file.Size > 5*1024*1024 {
		response.Error(c, appErrors.ErrBadRequest.WithMessage("头像大小不能超过5MB"))
		return
	}

	// TODO: 将文件上传至对象存储，获取公开URL
	// 当前返回占位URL
	placeholderURL := "https://static.kidsfit.example.com/avatars/placeholder.png"

	// 更新用户头像URL
	req := user.UpdateUserRequest{
		Avatar: placeholderURL,
	}
	result, err := h.userAppService.UpdateUser(c.Request.Context(), userID, &req)
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

// UpdateParentSettings 更新家长设置
// PUT /api/v1/settings
// 仅家长可操作，请求体包含待更新的设置项
func (h *UserHandler) UpdateParentSettings(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Error(c, appErrors.ErrUnauthorized)
		return
	}

	var dto user.ParentSettingsDTO
	// 绑定并校验请求参数
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, appErrors.ErrBadRequest.WithMessage(err.Error()))
		return
	}

	// 调用应用服务更新家长设置
	result, err := h.userAppService.UpdateParentSettings(c.Request.Context(), userID, &dto)
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
