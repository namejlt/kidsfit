package errors

// AppError 应用统一错误结构体，包含错误码和错误信息
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// New 创建一个新的AppError实例
func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Error 实现error接口，返回错误码和错误信息的字符串表示
func (e *AppError) Error() string {
	return e.Message
}

// WithMessage 返回一个具有相同错误码但不同错误信息的AppError副本
// 用于在保持错误码不变的情况下自定义错误信息
func (e *AppError) WithMessage(msg string) *AppError {
	return &AppError{
		Code:    e.Code,
		Message: msg,
	}
}

// 通用错误码 0-999
var (
	ErrSuccess         = New(0, "success")
	ErrBadRequest      = New(400, "请求参数错误")
	ErrUnauthorized    = New(401, "未授权")
	ErrForbidden       = New(403, "无权限")
	ErrNotFound        = New(404, "资源不存在")
	ErrTooManyRequests = New(429, "请求过于频繁")
	ErrInternal        = New(500, "服务器内部错误")
)

// 用户服务错误码 1000-1999
var (
	ErrUserNotFound       = New(1001, "用户不存在")
	ErrUserAlreadyExists  = New(1002, "用户已存在")
	ErrWrongPassword      = New(1003, "密码错误")
	ErrInvalidToken       = New(1004, "Token无效")
	ErrTokenExpired       = New(1005, "Token已过期")
	ErrPhoneExists        = New(1006, "手机号已注册")
	ErrInvalidPhone       = New(1007, "手机号无效")
	ErrChildLimitExceeded = New(1008, "儿童数量已达上限")
)

// 训练服务错误码 2000-2999
var (
	ErrExerciseNotFound    = New(2001, "运动记录不存在")
	ErrPlanNotFound        = New(2002, "训练计划不存在")
	ErrInvalidExerciseType = New(2003, "无效运动类型")
	ErrPlanAlreadyExists   = New(2004, "训练计划已存在")
)

// 视力服务错误码 3000-3999
var (
	ErrVisionNotFound    = New(3001, "视力记录不存在")
	ErrInvalidVisionData = New(3002, "无效视力数据")
	ErrOCRFailed         = New(3003, "OCR识别失败")
)

// 激励服务错误码 4000-4999
var (
	ErrBadgeNotFound      = New(4001, "勋章不存在")
	ErrBadgeAlreadyEarned = New(4002, "勋章已获得")
	ErrChallengeNotFound  = New(4003, "挑战不存在")
	ErrChallengeExpired   = New(4004, "挑战已过期")
	ErrInsufficientPoints = New(4005, "积分不足")
)
