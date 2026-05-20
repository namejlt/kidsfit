package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/kidsfit/api/internal/pkg/config"
)

// Claims 自定义JWT声明结构体，包含用户身份信息
type Claims struct {
	UserID   string `json:"user_id"`
	UserType string `json:"user_type"`
	ParentID string `json:"parent_id"`
	jwt.RegisteredClaims
}

// TokenPair 令牌对结构体，包含访问令牌和刷新令牌
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // 访问令牌过期时间（秒）
}

// GenerateTokenPair 根据用户信息生成访问令牌和刷新令牌对
// userID: 用户ID，userType: 用户类型，parentID: 家长ID，cfg: JWT配置
func GenerateTokenPair(userID, userType, parentID string, cfg *config.JWTConfig) (*TokenPair, error) {
	now := time.Now()

	// 生成访问令牌
	accessClaims := Claims{
		UserID:   userID,
		UserType: userType,
		ParentID: parentID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.Issuer,
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(cfg.AccessTTL) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	accessToken, err := generateToken(accessClaims, cfg.Secret)
	if err != nil {
		return nil, fmt.Errorf("生成访问令牌失败: %w", err)
	}

	// 生成刷新令牌
	refreshClaims := Claims{
		UserID:   userID,
		UserType: userType,
		ParentID: parentID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.Issuer,
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(cfg.RefreshTTL) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	refreshToken, err := generateToken(refreshClaims, cfg.Secret)
	if err != nil {
		return nil, fmt.Errorf("生成刷新令牌失败: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(cfg.AccessTTL),
	}, nil
}

// generateToken 根据声明和密钥生成JWT令牌字符串
func generateToken(claims Claims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken 解析并验证JWT令牌字符串，返回令牌中的声明信息
// tokenString: JWT令牌字符串，secret: 签名密钥
func ParseToken(tokenString string, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("不支持的签名算法: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("解析令牌失败: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("无效的令牌声明")
	}

	return claims, nil
}
