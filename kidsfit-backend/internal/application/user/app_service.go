package user

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/kidsfit/api/internal/domain/user"
	"github.com/kidsfit/api/internal/infrastructure/persistence/redis"
	appErrors "github.com/kidsfit/api/internal/pkg/errors"
	"github.com/kidsfit/api/internal/pkg/jwt"
	"github.com/kidsfit/api/internal/pkg/config"
	"github.com/kidsfit/api/internal/pkg/crypto"
)

// 手机号正则表达式，匹配中国大陆11位手机号
var phoneRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)

// 用户缓存键前缀
const (
	// userCacheKeyPrefix 用户缓存键前缀
	userCacheKeyPrefix = "user:profile:"
	// userCacheTTL 用户缓存过期时间
	userCacheTTL = 30 * time.Minute
)

// UserAppService 用户应用服务，负责用户相关的业务逻辑编排
type UserAppService struct {
	userRepo     user.UserRepository
	familyRepo   user.FamilyRepository
	settingsRepo user.ParentSettingsRepository
	redisCache   *redis.RedisCache
	jwtCfg       *config.JWTConfig
}

// NewUserAppService 创建用户应用服务实例
// userRepo: 用户仓储，familyRepo: 家庭关系仓储，settingsRepo: 家长设置仓储
// redisCache: Redis缓存客户端，jwtCfg: JWT配置
func NewUserAppService(
	userRepo user.UserRepository,
	familyRepo user.FamilyRepository,
	settingsRepo user.ParentSettingsRepository,
	redisCache *redis.RedisCache,
	jwtCfg *config.JWTConfig,
) *UserAppService {
	return &UserAppService{
		userRepo:     userRepo,
		familyRepo:   familyRepo,
		settingsRepo: settingsRepo,
		redisCache:   redisCache,
		jwtCfg:       jwtCfg,
	}
}

// Register 注册家长账号
// 校验手机号格式、检查手机号是否已注册、密码bcrypt加密、创建用户、生成JWT令牌
// ctx: 上下文，req: 注册请求DTO
func (s *UserAppService) Register(ctx context.Context, req *RegisterRequest) (*LoginResponse, error) {
	// 校验手机号格式
	if !phoneRegex.MatchString(req.Phone) {
		return nil, appErrors.ErrInvalidPhone
	}

	// 检查手机号是否已注册
	existingUser, err := s.userRepo.GetByPhone(ctx, req.Phone)
	if err == nil && existingUser != nil {
		return nil, appErrors.ErrPhoneExists
	}

	// 密码bcrypt加密
	passwordHash, err := crypto.HashPassword(req.Password)
	if err != nil {
		return nil, appErrors.ErrInternal.WithMessage("密码加密失败")
	}

	// 创建用户
	newUser := user.NewUser(user.UserTypeParent, req.Nickname)
	newUser.Phone = &req.Phone
	newUser.PasswordHash = &passwordHash

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, appErrors.ErrInternal.WithMessage("创建用户失败")
	}

	// 创建默认家长设置
	settings := user.NewParentSettings(newUser.ID)
	if err := s.settingsRepo.Create(ctx, settings); err != nil {
		return nil, appErrors.ErrInternal.WithMessage("创建家长设置失败")
	}

	// 生成JWT令牌
	tokenPair, err := jwt.GenerateTokenPair(newUser.ID, string(newUser.Type), "", s.jwtCfg)
	if err != nil {
		return nil, appErrors.ErrInternal.WithMessage("生成令牌失败")
	}

	return &LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		User:         toUserDTO(newUser),
	}, nil
}

// Login 用户登录
// 查找用户、验证密码、生成JWT令牌
// ctx: 上下文，req: 登录请求DTO
func (s *UserAppService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// 根据手机号查找用户
	u, err := s.userRepo.GetByPhone(ctx, req.Phone)
	if err != nil {
		return nil, appErrors.ErrUserNotFound
	}

	// 验证密码
	if u.PasswordHash == nil || !crypto.CheckPassword(req.Password, *u.PasswordHash) {
		return nil, appErrors.ErrWrongPassword
	}

	// 检查用户状态
	if !u.IsActive() {
		return nil, appErrors.ErrForbidden.WithMessage("用户已被禁用")
	}

	// 获取家长ID（用于儿童账号）
	parentID := ""
	if u.IsChild() && u.ParentID != nil {
		parentID = *u.ParentID
	}

	// 生成JWT令牌
	tokenPair, err := jwt.GenerateTokenPair(u.ID, string(u.Type), parentID, s.jwtCfg)
	if err != nil {
		return nil, appErrors.ErrInternal.WithMessage("生成令牌失败")
	}

	return &LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		User:         toUserDTO(u),
	}, nil
}

// RefreshToken 刷新访问令牌
// 解析刷新令牌，验证有效性后生成新的令牌对
// ctx: 上下文，refreshToken: 刷新令牌字符串
func (s *UserAppService) RefreshToken(ctx context.Context, refreshToken string) (*LoginResponse, error) {
	// 解析刷新令牌
	claims, err := jwt.ParseToken(refreshToken, s.jwtCfg.Secret)
	if err != nil {
		return nil, appErrors.ErrInvalidToken
	}

	// 验证用户是否仍然有效
	u, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, appErrors.ErrUserNotFound
	}

	if !u.IsActive() {
		return nil, appErrors.ErrForbidden.WithMessage("用户已被禁用")
	}

	// 获取家长ID
	parentID := ""
	if u.IsChild() && u.ParentID != nil {
		parentID = *u.ParentID
	}

	// 生成新的令牌对
	tokenPair, err := jwt.GenerateTokenPair(u.ID, string(u.Type), parentID, s.jwtCfg)
	if err != nil {
		return nil, appErrors.ErrInternal.WithMessage("生成令牌失败")
	}

	return &LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		User:         toUserDTO(u),
	}, nil
}

// GetCurrentUser 获取当前用户信息，使用Cache-Aside模式缓存用户资料
// 先从缓存获取，缓存未命中则从数据库加载并写入缓存
// ctx: 上下文，userID: 用户ID
func (s *UserAppService) GetCurrentUser(ctx context.Context, userID string) (*UserDTO, error) {
	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("%s%s", userCacheKeyPrefix, userID)
	var dto UserDTO
	if err := s.redisCache.GetJSON(ctx, cacheKey, &dto); err == nil && dto.ID != "" {
		return &dto, nil
	}

	// 缓存未命中，从数据库获取
	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, appErrors.ErrUserNotFound
	}

	// 写入缓存
	userDTO := toUserDTO(u)
	_ = s.redisCache.SetJSON(ctx, cacheKey, userDTO, userCacheTTL)

	return userDTO, nil
}

// UpdateUser 更新用户信息，更新后清除缓存
// ctx: 上下文，userID: 用户ID，req: 更新请求DTO
func (s *UserAppService) UpdateUser(ctx context.Context, userID string, req *UpdateUserRequest) (*UserDTO, error) {
	// 获取用户
	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, appErrors.ErrUserNotFound
	}

	// 更新字段
	if req.Nickname != "" {
		u.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		u.Avatar = req.Avatar
	}

	// 持久化更新
	if err := s.userRepo.Update(ctx, u); err != nil {
		return nil, appErrors.ErrInternal.WithMessage("更新用户失败")
	}

	// 清除缓存
	cacheKey := fmt.Sprintf("%s%s", userCacheKeyPrefix, userID)
	_ = s.redisCache.Delete(ctx, cacheKey)

	return toUserDTO(u), nil
}

// AddChild 添加儿童账号，创建儿童用户和家庭关系
// ctx: 上下文，parentID: 家长ID，req: 添加儿童请求DTO
func (s *UserAppService) AddChild(ctx context.Context, parentID string, req *AddChildRequest) (*ChildDTO, error) {
	// 验证家长是否存在
	parent, err := s.userRepo.GetByID(ctx, parentID)
	if err != nil {
		return nil, appErrors.ErrUserNotFound
	}
	if !parent.IsParent() {
		return nil, appErrors.ErrForbidden.WithMessage("仅家长可添加儿童")
	}

	// 检查儿童数量上限
	families, err := s.familyRepo.GetByParentID(ctx, parentID)
	if err != nil {
		return nil, appErrors.ErrInternal.WithMessage("查询家庭关系失败")
	}
	if len(families) >= 5 {
		return nil, appErrors.ErrChildLimitExceeded
	}

	// 创建儿童用户
	child := user.NewUser(user.UserTypeChild, req.Nickname)
	child.ParentID = &parentID
	child.Age = &req.Age
	if req.Avatar != "" {
		child.Avatar = req.Avatar
	}

	// 校验年龄
	if err := child.ValidateAge(); err != nil {
		return nil, appErrors.ErrBadRequest.WithMessage(err.Error())
	}

	if err := s.userRepo.Create(ctx, child); err != nil {
		return nil, appErrors.ErrInternal.WithMessage("创建儿童用户失败")
	}

	// 创建家庭关系
	family := user.NewFamily(parentID, child.ID, user.RelationOther)
	if err := s.familyRepo.Create(ctx, family); err != nil {
		return nil, appErrors.ErrInternal.WithMessage("创建家庭关系失败")
	}

	return toChildDTO(child), nil
}

// GetChildren 获取家长下的所有儿童列表
// ctx: 上下文，parentID: 家长ID
func (s *UserAppService) GetChildren(ctx context.Context, parentID string) ([]*ChildDTO, error) {
	// 查询家庭关系
	families, err := s.familyRepo.GetByParentID(ctx, parentID)
	if err != nil {
		return nil, appErrors.ErrInternal.WithMessage("查询家庭关系失败")
	}

	// 获取每个儿童的详细信息
	children := make([]*ChildDTO, 0, len(families))
	for _, family := range families {
		child, err := s.userRepo.GetByID(ctx, family.ChildID)
		if err != nil {
			continue
		}
		children = append(children, toChildDTO(child))
	}

	return children, nil
}

// GetParentSettings 获取家长设置
// ctx: 上下文，parentID: 家长ID
func (s *UserAppService) GetParentSettings(ctx context.Context, parentID string) (*ParentSettingsDTO, error) {
	settings, err := s.settingsRepo.GetByParentID(ctx, parentID)
	if err != nil {
		return nil, appErrors.ErrNotFound.WithMessage("家长设置不存在")
	}

	return toParentSettingsDTO(settings), nil
}

// UpdateParentSettings 更新家长设置
// ctx: 上下文，parentID: 家长ID，dto: 家长设置DTO
func (s *UserAppService) UpdateParentSettings(ctx context.Context, parentID string, dto *ParentSettingsDTO) (*ParentSettingsDTO, error) {
	// 获取现有设置
	settings, err := s.settingsRepo.GetByParentID(ctx, parentID)
	if err != nil {
		return nil, appErrors.ErrNotFound.WithMessage("家长设置不存在")
	}

	// 更新字段
	settings.DailyLimitMin = dto.DailyLimitMin
	settings.AvailableFrom = dto.AvailableFrom
	settings.AvailableTo = dto.AvailableTo
	settings.CameraAllowed = dto.CameraAllowed
	settings.LocationAllowed = dto.LocationAllowed
	settings.DataUploadCloud = dto.DataUploadCloud

	// 校验每日时长限制
	if err := settings.ValidateDailyLimit(); err != nil {
		return nil, appErrors.ErrBadRequest.WithMessage(err.Error())
	}

	// 持久化更新
	if err := s.settingsRepo.Update(ctx, settings); err != nil {
		return nil, appErrors.ErrInternal.WithMessage("更新家长设置失败")
	}

	return toParentSettingsDTO(settings), nil
}

// toUserDTO 将用户领域模型转换为DTO
func toUserDTO(u *user.User) *UserDTO {
	dto := &UserDTO{
		ID:        u.ID,
		Type:      string(u.Type),
		Nickname:  u.Nickname,
		Avatar:    u.Avatar,
		Status:    string(u.Status),
		CreatedAt: u.CreatedAt,
	}
	if u.Phone != nil {
		dto.Phone = *u.Phone
	}
	if u.Age != nil {
		dto.Age = *u.Age
	}
	return dto
}

// toChildDTO 将用户领域模型转换为儿童DTO
func toChildDTO(u *user.User) *ChildDTO {
	dto := &ChildDTO{
		ID:       u.ID,
		Nickname: u.Nickname,
		Avatar:   u.Avatar,
		Status:   string(u.Status),
		AgeGroup: u.GetAgeGroup(),
	}
	if u.Age != nil {
		dto.Age = *u.Age
	}
	return dto
}

// toParentSettingsDTO 将家长设置领域模型转换为DTO
func toParentSettingsDTO(ps *user.ParentSettings) *ParentSettingsDTO {
	return &ParentSettingsDTO{
		DailyLimitMin:   ps.DailyLimitMin,
		AvailableFrom:   ps.AvailableFrom,
		AvailableTo:     ps.AvailableTo,
		CameraAllowed:   ps.CameraAllowed,
		LocationAllowed: ps.LocationAllowed,
		DataUploadCloud: ps.DataUploadCloud,
	}
}
