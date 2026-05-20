package reward

import (
	"context"
	"time"

	"github.com/kidsfit/api/internal/domain/reward"
	"github.com/kidsfit/api/internal/domain/user"
	"github.com/kidsfit/api/internal/infrastructure/persistence/redis"
	appErrors "github.com/kidsfit/api/internal/pkg/errors"
	"github.com/kidsfit/api/internal/pkg/response"
)

// RewardAppService 激励应用服务，负责徽章、积分、挑战和排行榜相关的业务逻辑编排
type RewardAppService struct {
	badgeRepo     reward.BadgeRepository
	userBadgeRepo reward.UserBadgeRepository
	pointRepo     reward.PointRecordRepository
	challengeRepo reward.ChallengeRepository
	userRepo      user.UserRepository
	familyRepo    user.FamilyRepository
	leaderboard   *redis.Leaderboard
}

// NewRewardAppService 创建激励应用服务实例
// badgeRepo: 徽章仓储，userBadgeRepo: 用户徽章仓储，pointRepo: 积分记录仓储
// challengeRepo: 挑战仓储，userRepo: 用户仓储，familyRepo: 家庭关系仓储
// leaderboard: Redis排行榜
func NewRewardAppService(
	badgeRepo reward.BadgeRepository,
	userBadgeRepo reward.UserBadgeRepository,
	pointRepo reward.PointRecordRepository,
	challengeRepo reward.ChallengeRepository,
	userRepo user.UserRepository,
	familyRepo user.FamilyRepository,
	leaderboard *redis.Leaderboard,
) *RewardAppService {
	return &RewardAppService{
		badgeRepo:     badgeRepo,
		userBadgeRepo: userBadgeRepo,
		pointRepo:     pointRepo,
		challengeRepo: challengeRepo,
		userRepo:      userRepo,
		familyRepo:    familyRepo,
		leaderboard:   leaderboard,
	}
}

// GetBadges 按类别查询徽章列表，并标记用户是否已获得
// ctx: 上下文，userID: 用户ID，category: 徽章类别（空字符串表示全部）
func (s *RewardAppService) GetBadges(ctx context.Context, userID string, category string) ([]*BadgeDTO, error) {
	// 查询徽章列表
	var categoryFilter *reward.BadgeCategory
	if category != "" {
		cat := reward.BadgeCategory(category)
		categoryFilter = &cat
	}

	badges, err := s.badgeRepo.List(ctx, categoryFilter)
	if err != nil {
		return nil, appErrors.ErrInternal.WithMessage("查询徽章列表失败")
	}

	// 查询用户已获得的徽章
	userBadges, err := s.userBadgeRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, appErrors.ErrInternal.WithMessage("查询用户徽章失败")
	}

	// 构建已获得徽章的集合
	earnedSet := make(map[string]bool)
	for _, ub := range userBadges {
		earnedSet[ub.BadgeID] = true
	}

	// 转换为DTO并标记是否已获得
	dtos := make([]*BadgeDTO, 0, len(badges))
	for _, badge := range badges {
		dto := toBadgeDTO(badge)
		dto.Earned = earnedSet[badge.ID]
		dtos = append(dtos, dto)
	}

	return dtos, nil
}

// GetMyBadges 获取用户已获得的所有徽章
// ctx: 上下文，userID: 用户ID
func (s *RewardAppService) GetMyBadges(ctx context.Context, userID string) ([]*BadgeDTO, error) {
	// 查询用户已获得的徽章
	userBadges, err := s.userBadgeRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, appErrors.ErrInternal.WithMessage("查询用户徽章失败")
	}

	dtos := make([]*BadgeDTO, 0, len(userBadges))
	for _, ub := range userBadges {
		badge, err := s.badgeRepo.GetByID(ctx, ub.BadgeID)
		if err != nil {
			continue
		}
		dto := toBadgeDTO(badge)
		dto.Earned = true
		dtos = append(dtos, dto)
	}

	return dtos, nil
}

// GetPoints 分页查询用户积分记录
// ctx: 上下文，userID: 用户ID，page: 页码，pageSize: 每页大小
func (s *RewardAppService) GetPoints(ctx context.Context, userID string, page, pageSize int64) ([]*PointRecordDTO, *response.Pagination, error) {
	pagination := reward.Pagination{
		Page:     int(page),
		PageSize: int(pageSize),
	}

	result, err := s.pointRepo.GetByUserID(ctx, userID, pagination)
	if err != nil {
		return nil, nil, appErrors.ErrInternal.WithMessage("查询积分记录失败")
	}

	dtos := make([]*PointRecordDTO, 0, len(result.Items))
	for _, record := range result.Items {
		dtos = append(dtos, toPointRecordDTO(&record))
	}

	pag := &response.Pagination{
		Page:       int64(result.Page),
		PageSize:   int64(result.PageSize),
		Total:      result.Total,
		TotalPages: int64(result.TotalPages),
	}

	return dtos, pag, nil
}

// GetPointsBalance 获取用户当前积分余额
// ctx: 上下文，userID: 用户ID
func (s *RewardAppService) GetPointsBalance(ctx context.Context, userID string) (int64, error) {
	balance, err := s.pointRepo.GetBalanceByUserID(ctx, userID)
	if err != nil {
		return 0, appErrors.ErrInternal.WithMessage("查询积分余额失败")
	}
	return int64(balance), nil
}

// CreateChallenge 创建挑战
// ctx: 上下文，initiatorID: 发起者ID，req: 创建挑战请求DTO
func (s *RewardAppService) CreateChallenge(ctx context.Context, initiatorID string, req *CreateChallengeRequest) (*ChallengeDTO, error) {
	// 创建挑战领域模型
	challenge := reward.NewChallenge(
		reward.ChallengeType(req.Type),
		initiatorID,
		req.ExerciseType,
		req.TargetValue,
	)

	// 设置接受者（如果有）
	if req.AcceptorID != "" {
		challenge.AcceptorID = &req.AcceptorID
	}

	// 设置过期时间（异步挑战24小时后过期）
	if challenge.Type == reward.ChallengeTypeAsync {
		expiresAt := time.Now().Add(24 * time.Hour)
		challenge.ExpiresAt = &expiresAt
	}

	// 持久化挑战
	if err := s.challengeRepo.Create(ctx, challenge); err != nil {
		return nil, appErrors.ErrInternal.WithMessage("创建挑战失败")
	}

	return toChallengeDTO(challenge), nil
}

// AcceptChallenge 接受挑战
// 校验挑战状态和过期时间后更新挑战状态
// ctx: 上下文，userID: 接受者ID，challengeID: 挑战ID
func (s *RewardAppService) AcceptChallenge(ctx context.Context, userID string, challengeID string) (*ChallengeDTO, error) {
	// 获取挑战
	challenge, err := s.challengeRepo.GetByID(ctx, challengeID)
	if err != nil {
		return nil, appErrors.ErrChallengeNotFound
	}

	// 校验是否可以接受
	if !challenge.CanAccept() {
		if challenge.IsExpired() {
			return nil, appErrors.ErrChallengeExpired
		}
		return nil, appErrors.ErrBadRequest.WithMessage("挑战无法被接受")
	}

	// 不能接受自己发起的挑战
	if challenge.InitiatorID == userID {
		return nil, appErrors.ErrBadRequest.WithMessage("不能接受自己发起的挑战")
	}

	// 更新挑战状态
	challenge.AcceptorID = &userID
	challenge.Status = reward.ChallengeStatusAccepted

	if err := s.challengeRepo.Update(ctx, challenge); err != nil {
		return nil, appErrors.ErrInternal.WithMessage("接受挑战失败")
	}

	return toChallengeDTO(challenge), nil
}

// SubmitChallengeScore 提交挑战成绩
// 根据提交者身份更新对应成绩，双方都提交后判定胜负
// ctx: 上下文，userID: 提交者ID，challengeID: 挑战ID，score: 成绩
func (s *RewardAppService) SubmitChallengeScore(ctx context.Context, userID string, challengeID string, score int) (*ChallengeDTO, error) {
	// 获取挑战
	challenge, err := s.challengeRepo.GetByID(ctx, challengeID)
	if err != nil {
		return nil, appErrors.ErrChallengeNotFound
	}

	// 校验是否可以提交成绩
	if !challenge.CanSubmit() {
		if challenge.IsExpired() {
			return nil, appErrors.ErrChallengeExpired
		}
		return nil, appErrors.ErrBadRequest.WithMessage("挑战无法提交成绩")
	}

	// 根据提交者身份更新成绩
	now := time.Now()
	if userID == challenge.InitiatorID {
		challenge.InitiatorScore = &score
	} else if challenge.AcceptorID != nil && userID == *challenge.AcceptorID {
		challenge.AcceptorScore = &score
	} else {
		return nil, appErrors.ErrBadRequest.WithMessage("非挑战参与者")
	}

	// 判断是否双方都已提交成绩
	if challenge.InitiatorScore != nil && challenge.AcceptorScore != nil {
		challenge.Status = reward.ChallengeStatusCompleted
		challenge.CompletedAt = &now

		// 判定胜负
		if *challenge.InitiatorScore > *challenge.AcceptorScore {
			challenge.WinnerID = &challenge.InitiatorID
		} else if *challenge.AcceptorScore > *challenge.InitiatorScore {
			challenge.WinnerID = challenge.AcceptorID
		}
		// 平局不设置WinnerID
	}

	if err := s.challengeRepo.Update(ctx, challenge); err != nil {
		return nil, appErrors.ErrInternal.WithMessage("提交挑战成绩失败")
	}

	return toChallengeDTO(challenge), nil
}

// GetChallenges 获取用户相关的挑战列表
// 查询用户发起和接受的挑战，合并后按创建时间降序排列
// ctx: 上下文，userID: 用户ID
func (s *RewardAppService) GetChallenges(ctx context.Context, userID string) ([]*ChallengeDTO, error) {
	// 查询用户发起的挑战
	initiated, err := s.challengeRepo.GetByInitiatorID(ctx, userID)
	if err != nil {
		return nil, appErrors.ErrInternal.WithMessage("查询发起的挑战失败")
	}

	// 查询用户接受的挑战
	accepted, err := s.challengeRepo.GetByAcceptorID(ctx, userID)
	if err != nil {
		return nil, appErrors.ErrInternal.WithMessage("查询接受的挑战失败")
	}

	// 合并并去重
	seen := make(map[string]bool)
	dtos := make([]*ChallengeDTO, 0, len(initiated)+len(accepted))

	for _, c := range initiated {
		if !seen[c.ID] {
			seen[c.ID] = true
			dtos = append(dtos, toChallengeDTO(c))
		}
	}
	for _, c := range accepted {
		if !seen[c.ID] {
			seen[c.ID] = true
			dtos = append(dtos, toChallengeDTO(c))
		}
	}

	return dtos, nil
}

// GetFamilyLeaderboard 获取家庭排行榜
// 查询家庭下所有成员，从Redis ZSet获取排名
// ctx: 上下文，parentID: 家长ID
func (s *RewardAppService) GetFamilyLeaderboard(ctx context.Context, parentID string) ([]*LeaderboardDTO, error) {
	// 查询家庭关系
	families, err := s.familyRepo.GetByParentID(ctx, parentID)
	if err != nil {
		return nil, appErrors.ErrInternal.WithMessage("查询家庭关系失败")
	}

	// 收集所有家庭成员ID（包含家长自己）
	memberIDs := []string{parentID}
	for _, family := range families {
		memberIDs = append(memberIDs, family.ChildID)
	}

	// 从Redis排行榜获取排名
	leaderboardKey := "leaderboard:family:" + parentID
	memberScores, err := s.leaderboard.GetTopN(ctx, leaderboardKey, int64(len(memberIDs)))
	if err != nil {
		return nil, appErrors.ErrInternal.WithMessage("获取家庭排行榜失败")
	}

	// 构建排行榜DTO
	dtos := make([]*LeaderboardDTO, 0, len(memberScores))
	for i, ms := range memberScores {
		// 获取用户信息
		u, err := s.userRepo.GetByID(ctx, ms.Member)
		if err != nil {
			continue
		}
		dtos = append(dtos, &LeaderboardDTO{
			Rank:     int64(i) + 1,
			UserID:   u.ID,
			Nickname: u.Nickname,
			Avatar:   u.Avatar,
			Score:    ms.Score,
		})
	}

	return dtos, nil
}

// GetGlobalLeaderboard 获取全局排行榜
// 从Redis ZSet获取指定运动类型的全球排名
// ctx: 上下文，exerciseType: 运动类型，limit: 获取数量
func (s *RewardAppService) GetGlobalLeaderboard(ctx context.Context, exerciseType string, limit int64) ([]*LeaderboardDTO, error) {
	// 从Redis排行榜获取排名
	leaderboardKey := "leaderboard:global:" + exerciseType
	memberScores, err := s.leaderboard.GetTopN(ctx, leaderboardKey, limit)
	if err != nil {
		return nil, appErrors.ErrInternal.WithMessage("获取全局排行榜失败")
	}

	// 构建排行榜DTO
	dtos := make([]*LeaderboardDTO, 0, len(memberScores))
	for i, ms := range memberScores {
		// 获取用户信息
		u, err := s.userRepo.GetByID(ctx, ms.Member)
		if err != nil {
			continue
		}
		dtos = append(dtos, &LeaderboardDTO{
			Rank:     int64(i) + 1,
			UserID:   u.ID,
			Nickname: u.Nickname,
			Avatar:   u.Avatar,
			Score:    ms.Score,
		})
	}

	return dtos, nil
}

// toBadgeDTO 将徽章领域模型转换为DTO
func toBadgeDTO(b *reward.Badge) *BadgeDTO {
	return &BadgeDTO{
		ID:          b.ID,
		Code:        b.Code,
		Name:        b.Name,
		Description: b.Description,
		Category:    string(b.Category),
		Icon:        b.Icon,
		Points:      b.Points,
	}
}

// toPointRecordDTO 将积分记录领域模型转换为DTO
func toPointRecordDTO(r *reward.PointRecord) *PointRecordDTO {
	return &PointRecordDTO{
		ID:          r.ID,
		UserID:      r.UserID,
		Points:      r.Points,
		Type:        string(r.Type),
		SourceID:    stringToPtr(r.SourceID),
		SourceType:  stringToPtr(r.SourceType),
		Description: r.Description,
		Balance:     r.Balance,
		CreatedAt:   r.CreatedAt,
	}
}

// toChallengeDTO 将挑战领域模型转换为DTO
func toChallengeDTO(c *reward.Challenge) *ChallengeDTO {
	return &ChallengeDTO{
		ID:             c.ID,
		Type:           string(c.Type),
		InitiatorID:    c.InitiatorID,
		AcceptorID:     c.AcceptorID,
		ExerciseType:   c.ExerciseType,
		TargetValue:    c.TargetValue,
		InitiatorScore: c.InitiatorScore,
		AcceptorScore:  c.AcceptorScore,
		WinnerID:       c.WinnerID,
		Status:         string(c.Status),
		ExpiresAt:      c.ExpiresAt,
		CompletedAt:    c.CompletedAt,
		CreatedAt:      c.CreatedAt,
	}
}

// stringToPtr 将字符串转换为字符串指针，空字符串返回nil
func stringToPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
