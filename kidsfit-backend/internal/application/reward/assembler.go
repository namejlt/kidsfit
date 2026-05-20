package reward

import (
	domain "github.com/kidsfit/api/internal/domain/reward"
)

// BadgeToDTO 将徽章领域模型转换为徽章DTO，earned参数标识用户是否已获得该徽章
func BadgeToDTO(badge *domain.Badge, earned bool) *BadgeDTO {
	if badge == nil {
		return nil
	}
	return &BadgeDTO{
		ID:          badge.ID,
		Code:        badge.Code,
		Name:        badge.Name,
		Description: badge.Description,
		Category:    string(badge.Category),
		Icon:        badge.Icon,
		Points:      badge.Points,
		Earned:      earned,
	}
}

// UserBadgeToDTO 将用户徽章领域模型转换为用户徽章DTO
func UserBadgeToDTO(userBadge *domain.UserBadge) *UserBadgeDTO {
	if userBadge == nil {
		return nil
	}
	return &UserBadgeDTO{
		ID:       userBadge.ID,
		UserID:   userBadge.UserID,
		BadgeID:  userBadge.BadgeID,
		EarnedAt: userBadge.EarnedAt,
	}
}

// PointRecordToDTO 将积分记录领域模型转换为积分记录DTO
func PointRecordToDTO(record *domain.PointRecord) *PointRecordDTO {
	if record == nil {
		return nil
	}
	dto := &PointRecordDTO{
		ID:          record.ID,
		UserID:      record.UserID,
		Points:      record.Points,
		Type:        string(record.Type),
		Description: record.Description,
		Balance:     record.Balance,
		CreatedAt:   record.CreatedAt,
	}
	if record.SourceID != "" {
		sourceID := record.SourceID
		dto.SourceID = &sourceID
	}
	if record.SourceType != "" {
		sourceType := record.SourceType
		dto.SourceType = &sourceType
	}
	return dto
}

// ChallengeToDTO 将挑战领域模型转换为挑战DTO
func ChallengeToDTO(challenge *domain.Challenge) *ChallengeDTO {
	if challenge == nil {
		return nil
	}
	return &ChallengeDTO{
		ID:             challenge.ID,
		Type:           string(challenge.Type),
		InitiatorID:    challenge.InitiatorID,
		AcceptorID:     challenge.AcceptorID,
		ExerciseType:   challenge.ExerciseType,
		TargetValue:    challenge.TargetValue,
		InitiatorScore: challenge.InitiatorScore,
		AcceptorScore:  challenge.AcceptorScore,
		WinnerID:       challenge.WinnerID,
		Status:         string(challenge.Status),
		ExpiresAt:      challenge.ExpiresAt,
		CompletedAt:    challenge.CompletedAt,
		CreatedAt:      challenge.CreatedAt,
	}
}

// CreateChallengeRequestToChallenge 将创建挑战请求DTO转换为挑战领域模型
func CreateChallengeRequestToChallenge(req *CreateChallengeRequest, initiatorID string) *domain.Challenge {
	if req == nil {
		return nil
	}
	challenge := domain.NewChallenge(
		domain.ChallengeType(req.Type),
		initiatorID,
		req.ExerciseType,
		req.TargetValue,
	)
	if req.AcceptorID != "" {
		acceptorID := req.AcceptorID
		challenge.AcceptorID = &acceptorID
	}
	return challenge
}

// LeaderboardEntryToDTO 将排行榜条目参数转换为排行榜DTO
func LeaderboardEntryToDTO(rank int64, userID, nickname, avatar string, score float64) *LeaderboardDTO {
	return &LeaderboardDTO{
		Rank:     rank,
		UserID:   userID,
		Nickname: nickname,
		Avatar:   avatar,
		Score:    score,
	}
}
