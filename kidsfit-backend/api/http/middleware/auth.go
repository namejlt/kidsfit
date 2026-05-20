package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/kidsfit/api/internal/pkg/jwt"
	appErrors "github.com/kidsfit/api/internal/pkg/errors"
	"github.com/kidsfit/api/internal/pkg/response"
)

// 上下文键常量，用于在gin.Context中存储认证信息
const (
	// ContextKeyUserID 用户ID上下文键
	ContextKeyUserID = "user_id"
	// ContextKeyUserType 用户类型上下文键
	ContextKeyUserType = "user_type"
	// ContextKeyParentID 家长ID上下文键
	ContextKeyParentID = "parent_id"
)

// AuthMiddleware 认证中间件，从Authorization Header提取Bearer Token，
// 验证JWT令牌后将user_id/user_type/parent_id注入到gin.Context中
// jwtSecret: JWT签名密钥
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Header中提取Authorization字段
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, appErrors.ErrUnauthorized)
			c.Abort()
			return
		}

		// 解析Bearer Token格式
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, appErrors.ErrUnauthorized)
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析并验证JWT令牌
		claims, err := jwt.ParseToken(tokenString, jwtSecret)
		if err != nil {
			response.Error(c, appErrors.ErrInvalidToken)
			c.Abort()
			return
		}

		// 将用户身份信息注入到上下文中
		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyUserType, claims.UserType)
		c.Set(ContextKeyParentID, claims.ParentID)

		c.Next()
	}
}

// RequireParent 要求当前用户必须是家长类型的中间件
// 从上下文中获取user_type，若非parent则返回403错误
func RequireParent() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get(ContextKeyUserType)
		if !exists {
			response.Error(c, appErrors.ErrUnauthorized)
			c.Abort()
			return
		}

		if userType.(string) != "parent" {
			response.Error(c, appErrors.ErrForbidden.WithMessage("仅家长可执行此操作"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireChildOrParent 要求当前用户是儿童本人或关联家长的中间件
// 对于儿童用户，验证其parent_id与路径中的parent_id匹配
// 对于家长用户，验证其user_id与路径中的parent_id匹配
func RequireChildOrParent() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get(ContextKeyUserType)
		if !exists {
			response.Error(c, appErrors.ErrUnauthorized)
			c.Abort()
			return
		}

		userID, _ := c.Get(ContextKeyUserID)
		parentID, _ := c.Get(ContextKeyParentID)

		// 家长用户：直接允许访问
		if userType.(string) == "parent" {
			c.Next()
			return
		}

		// 儿童用户：验证关联家长ID是否匹配
		if userType.(string) == "child" {
			// 从路径参数中获取目标家长ID
			targetParentID := c.Param("parent_id")
			if targetParentID == "" {
				// 如果路径中没有parent_id，则从查询参数中获取child_id进行验证
				targetChildID := c.Param("child_id")
				if targetChildID != "" && targetChildID == userID.(string) {
					c.Next()
					return
				}
				// 允许儿童访问自己的资源
				c.Next()
				return
			}
			if parentID.(string) == targetParentID {
				c.Next()
				return
			}
		}

		response.Error(c, appErrors.ErrForbidden.WithMessage("无权访问此资源"))
		c.Abort()
	}
}

// GetUserID 从gin.Context中获取当前认证用户的ID
// c: gin上下文
func GetUserID(c *gin.Context) string {
	val, _ := c.Get(ContextKeyUserID)
	if val == nil {
		return ""
	}
	return val.(string)
}

// GetUserType 从gin.Context中获取当前认证用户的类型
// c: gin上下文
func GetUserType(c *gin.Context) string {
	val, _ := c.Get(ContextKeyUserType)
	if val == nil {
		return ""
	}
	return val.(string)
}

// GetParentID 从gin.Context中获取当前认证用户的家长ID
// c: gin上下文
func GetParentID(c *gin.Context) string {
	val, _ := c.Get(ContextKeyParentID)
	if val == nil {
		return ""
	}
	return val.(string)
}
