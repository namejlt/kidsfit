package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	appErrors "github.com/kidsfit/api/internal/pkg/errors"
)

// Response 统一API响应结构体
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Pagination 分页信息结构体
type Pagination struct {
	Page       int64 `json:"page"`
	PageSize   int64 `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"total_pages"`
}

// PageResponse 分页响应结构体
type PageResponse struct {
	List       interface{} `json:"list"`
	Pagination Pagination  `json:"pagination"`
}

// Success 返回成功的API响应
// data: 响应数据，可以为nil
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    appErrors.ErrSuccess.Code,
		Message: appErrors.ErrSuccess.Message,
		Data:    data,
	})
}

// Error 返回错误的API响应
// 根据错误码自动映射HTTP状态码
func Error(c *gin.Context, err *appErrors.AppError) {
	httpStatus := mapCodeToHTTPStatus(err.Code)
	c.JSON(httpStatus, Response{
		Code:    err.Code,
		Message: err.Message,
	})
}

// SuccessWithPage 返回带分页信息的成功API响应
// list: 数据列表，pagination: 分页信息
func SuccessWithPage(c *gin.Context, list interface{}, pagination *Pagination) {
	c.JSON(http.StatusOK, Response{
		Code:    appErrors.ErrSuccess.Code,
		Message: appErrors.ErrSuccess.Message,
		Data: PageResponse{
			List:       list,
			Pagination: *pagination,
		},
	})
}

// mapCodeToHTTPStatus 将应用错误码映射为HTTP状态码
func mapCodeToHTTPStatus(code int) int {
	switch {
	case code == 0:
		return http.StatusOK
	case code >= 400 && code < 500:
		return code
	case code >= 1000:
		// 业务错误码统一返回400
		return http.StatusBadRequest
	case code >= 500:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
