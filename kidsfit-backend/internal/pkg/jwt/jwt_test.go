package jwt

import (
	"testing"

	"github.com/kidsfit/api/internal/pkg/config"
)

// testJWTConfig 创建测试用的JWT配置
func testJWTConfig() *config.JWTConfig {
	return &config.JWTConfig{
		Secret:     "test-secret-key-for-unit-testing",
		AccessTTL:  3600,  // 1小时
		RefreshTTL: 86400, // 1天
		Issuer:     "kidsfit-test",
	}
}

// TestGenerateTokenPair 测试生成Token对
func TestGenerateTokenPair(t *testing.T) {
	cfg := testJWTConfig()

	pair, err := GenerateTokenPair("user-123", "parent", "", cfg)
	if err != nil {
		t.Fatalf("生成Token对失败: %v", err)
	}

	if pair.AccessToken == "" {
		t.Error("访问令牌不应为空")
	}
	if pair.RefreshToken == "" {
		t.Error("刷新令牌不应为空")
	}
	if pair.ExpiresIn != int64(cfg.AccessTTL) {
		t.Errorf("过期时间应为 %d，实际为 %d", cfg.AccessTTL, pair.ExpiresIn)
	}
	if pair.AccessToken == pair.RefreshToken {
		t.Error("访问令牌和刷新令牌不应相同")
	}
}

// TestParseToken_Valid 测试解析有效Token
func TestParseToken_Valid(t *testing.T) {
	cfg := testJWTConfig()

	pair, err := GenerateTokenPair("user-123", "parent", "parent-456", cfg)
	if err != nil {
		t.Fatalf("生成Token对失败: %v", err)
	}

	claims, err := ParseToken(pair.AccessToken, cfg.Secret)
	if err != nil {
		t.Fatalf("解析有效Token失败: %v", err)
	}

	if claims.UserID != "user-123" {
		t.Errorf("用户ID应为 user-123，实际为 %s", claims.UserID)
	}
	if claims.UserType != "parent" {
		t.Errorf("用户类型应为 parent，实际为 %s", claims.UserType)
	}
	if claims.ParentID != "parent-456" {
		t.Errorf("家长ID应为 parent-456，实际为 %s", claims.ParentID)
	}
	if claims.Issuer != "kidsfit-test" {
		t.Errorf("签发者应为 kidsfit-test，实际为 %s", claims.Issuer)
	}
}

// TestParseToken_Expired 测试解析过期Token
func TestParseToken_Expired(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:     "test-secret-key-for-unit-testing",
		AccessTTL:  -1, // 设置为负数使Token立即过期
		RefreshTTL: -1,
		Issuer:     "kidsfit-test",
	}

	pair, err := GenerateTokenPair("user-123", "parent", "", cfg)
	if err != nil {
		t.Fatalf("生成Token对失败: %v", err)
	}

	_, err = ParseToken(pair.AccessToken, cfg.Secret)
	if err == nil {
		t.Error("解析过期Token应返回错误")
	}
}

// TestParseToken_Invalid 测试解析无效Token
func TestParseToken_Invalid(t *testing.T) {
	cfg := testJWTConfig()

	t.Run("空字符串Token", func(t *testing.T) {
		_, err := ParseToken("", cfg.Secret)
		if err == nil {
			t.Error("解析空字符串Token应返回错误")
		}
	})

	t.Run("随机字符串Token", func(t *testing.T) {
		_, err := ParseToken("invalid.token.string", cfg.Secret)
		if err == nil {
			t.Error("解析无效Token应返回错误")
		}
	})

	t.Run("错误密钥解析", func(t *testing.T) {
		pair, err := GenerateTokenPair("user-123", "parent", "", cfg)
		if err != nil {
			t.Fatalf("生成Token对失败: %v", err)
		}
		_, err = ParseToken(pair.AccessToken, "wrong-secret-key")
		if err == nil {
			t.Error("使用错误密钥解析Token应返回错误")
		}
	})
}
