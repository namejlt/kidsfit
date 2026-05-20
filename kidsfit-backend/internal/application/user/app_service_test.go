package user

import (
	"context"
	"testing"

	"github.com/kidsfit/api/internal/domain/user"
	"github.com/kidsfit/api/internal/pkg/config"
	"github.com/kidsfit/api/internal/pkg/crypto"
	appErrors "github.com/kidsfit/api/internal/pkg/errors"
)

// testJWTConfig 创建测试用的JWT配置
func testJWTConfig() *config.JWTConfig {
	return &config.JWTConfig{
		Secret:     "test-secret-key-for-unit-testing",
		AccessTTL:  3600,
		RefreshTTL: 86400,
		Issuer:     "kidsfit-test",
	}
}

// newTestAppService 创建测试用的用户应用服务实例
// Register和Login方法不使用redisCache，因此传nil
func newTestAppService() (*UserAppService, *MockUserRepository, *MockFamilyRepository, *MockParentSettingsRepository) {
	userRepo := NewMockUserRepository()
	familyRepo := NewMockFamilyRepository()
	settingsRepo := NewMockParentSettingsRepository()
	jwtCfg := testJWTConfig()

	svc := NewUserAppService(userRepo, familyRepo, settingsRepo, nil, jwtCfg)
	return svc, userRepo, familyRepo, settingsRepo
}

// createTestParentUser 创建测试用的家长用户，密码使用bcrypt哈希
func createTestParentUser(phone, password, nickname string) *user.User {
	u := user.NewUser(user.UserTypeParent, nickname)
	u.Phone = &phone
	hash, _ := crypto.HashPassword(password)
	u.PasswordHash = &hash
	return u
}

// TestUserAppService_Register_Success 测试注册成功
func TestUserAppService_Register_Success(t *testing.T) {
	svc, _, _, _ := newTestAppService()

	resp, err := svc.Register(context.Background(), &RegisterRequest{
		Phone:    "13800138000",
		Password: "password123",
		Nickname: "测试家长",
	})
	if err != nil {
		t.Fatalf("注册应成功，实际返回错误: %v", err)
	}
	if resp.AccessToken == "" {
		t.Error("注册成功应返回访问令牌")
	}
	if resp.RefreshToken == "" {
		t.Error("注册成功应返回刷新令牌")
	}
	if resp.User == nil {
		t.Error("注册成功应返回用户信息")
	}
	if resp.User.Type != "parent" {
		t.Errorf("注册用户类型应为 parent，实际为 %s", resp.User.Type)
	}
	if resp.User.Nickname != "测试家长" {
		t.Errorf("昵称应为 测试家长，实际为 %s", resp.User.Nickname)
	}
}

// TestUserAppService_Register_PhoneExists 测试手机号已注册
func TestUserAppService_Register_PhoneExists(t *testing.T) {
	svc, userRepo, _, _ := newTestAppService()

	// 先注册一个用户
	phone := "13800138000"
	existingUser := createTestParentUser(phone, "password123", "已有用户")
	_ = userRepo.Create(context.Background(), existingUser)

	// 再用相同手机号注册应失败
	_, err := svc.Register(context.Background(), &RegisterRequest{
		Phone:    phone,
		Password: "password456",
		Nickname: "新用户",
	})
	if err == nil {
		t.Error("手机号已注册时应返回错误")
	}
	if err != appErrors.ErrPhoneExists {
		t.Errorf("应返回 ErrPhoneExists 错误，实际返回: %v", err)
	}
}

// TestUserAppService_Register_InvalidPhone 测试无效手机号
func TestUserAppService_Register_InvalidPhone(t *testing.T) {
	svc, _, _, _ := newTestAppService()

	invalidPhones := []string{
		"12345678901",   // 不以1开头第二位3-9
		"1380013800",    // 10位数字
		"138001380001",  // 12位数字
		"abc12345678",   // 包含字母
		"",              // 空字符串
	}

	for _, phone := range invalidPhones {
		_, err := svc.Register(context.Background(), &RegisterRequest{
			Phone:    phone,
			Password: "password123",
			Nickname: "测试用户",
		})
		if err != appErrors.ErrInvalidPhone {
			t.Errorf("手机号 %s 应返回 ErrInvalidPhone 错误，实际返回: %v", phone, err)
		}
	}
}

// TestUserAppService_Login_Success 测试登录成功
func TestUserAppService_Login_Success(t *testing.T) {
	svc, userRepo, _, _ := newTestAppService()

	// 先注册用户
	phone := "13800138000"
	password := "password123"
	existingUser := createTestParentUser(phone, password, "测试家长")
	_ = userRepo.Create(context.Background(), existingUser)

	// 登录
	resp, err := svc.Login(context.Background(), &LoginRequest{
		Phone:    phone,
		Password: password,
	})
	if err != nil {
		t.Fatalf("登录应成功，实际返回错误: %v", err)
	}
	if resp.AccessToken == "" {
		t.Error("登录成功应返回访问令牌")
	}
	if resp.User == nil {
		t.Error("登录成功应返回用户信息")
	}
	if resp.User.Phone != phone {
		t.Errorf("手机号应为 %s，实际为 %s", phone, resp.User.Phone)
	}
}

// TestUserAppService_Login_WrongPassword 测试密码错误
func TestUserAppService_Login_WrongPassword(t *testing.T) {
	svc, userRepo, _, _ := newTestAppService()

	// 先注册用户
	phone := "13800138000"
	existingUser := createTestParentUser(phone, "correctPassword", "测试家长")
	_ = userRepo.Create(context.Background(), existingUser)

	// 使用错误密码登录
	_, err := svc.Login(context.Background(), &LoginRequest{
		Phone:    phone,
		Password: "wrongPassword",
	})
	if err == nil {
		t.Error("密码错误时应返回错误")
	}
	if err != appErrors.ErrWrongPassword {
		t.Errorf("应返回 ErrWrongPassword 错误，实际返回: %v", err)
	}
}

// TestUserAppService_Login_UserNotFound 测试用户不存在
func TestUserAppService_Login_UserNotFound(t *testing.T) {
	svc, _, _, _ := newTestAppService()

	_, err := svc.Login(context.Background(), &LoginRequest{
		Phone:    "13800138000",
		Password: "password123",
	})
	if err == nil {
		t.Error("用户不存在时应返回错误")
	}
	if err != appErrors.ErrUserNotFound {
		t.Errorf("应返回 ErrUserNotFound 错误，实际返回: %v", err)
	}
}
