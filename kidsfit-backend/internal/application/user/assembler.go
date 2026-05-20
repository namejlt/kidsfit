package user

import (
	domain "github.com/kidsfit/api/internal/domain/user"
)

// UserToDTO 将用户领域模型转换为用户DTO
func UserToDTO(user *domain.User) *UserDTO {
	if user == nil {
		return nil
	}
	dto := &UserDTO{
		ID:        user.ID,
		Type:      string(user.Type),
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Status:    string(user.Status),
		CreatedAt: user.CreatedAt,
	}
	// 仅家长有手机号
	if user.Phone != nil {
		dto.Phone = *user.Phone
	}
	// 仅儿童有年龄
	if user.Age != nil {
		dto.Age = *user.Age
	}
	return dto
}

// DTOToUser 将用户DTO转换为用户领域模型
func DTOToUser(dto *UserDTO) *domain.User {
	if dto == nil {
		return nil
	}
	u := &domain.User{
		ID:       dto.ID,
		Type:     domain.UserType(dto.Type),
		Nickname: dto.Nickname,
		Avatar:   dto.Avatar,
		Status:   domain.UserStatus(dto.Status),
	}
	if dto.Phone != "" {
		phone := dto.Phone
		u.Phone = &phone
	}
	if dto.Age > 0 {
		age := dto.Age
		u.Age = &age
	}
	return u
}

// ChildToDTO 将用户领域模型转换为儿童DTO
func ChildToDTO(child *domain.User) *ChildDTO {
	if child == nil {
		return nil
	}
	dto := &ChildDTO{
		ID:       child.ID,
		Nickname: child.Nickname,
		Avatar:   child.Avatar,
		Status:   string(child.Status),
	}
	if child.Age != nil {
		dto.Age = *child.Age
		dto.AgeGroup = child.GetAgeGroup()
	}
	return dto
}

// FamilyToDTO 将家庭关系领域模型转换为家庭关系DTO
func FamilyToDTO(family *domain.Family) *FamilyDTO {
	if family == nil {
		return nil
	}
	return &FamilyDTO{
		ID:        family.ID,
		ParentID:  family.ParentID,
		ChildID:   family.ChildID,
		Relation:  string(family.Relation),
		CreatedAt: family.CreatedAt,
	}
}

// ParentSettingsToDTO 将家长设置领域模型转换为家长设置DTO
func ParentSettingsToDTO(settings *domain.ParentSettings) *ParentSettingsDTO {
	if settings == nil {
		return nil
	}
	return &ParentSettingsDTO{
		DailyLimitMin:   settings.DailyLimitMin,
		AvailableFrom:   settings.AvailableFrom,
		AvailableTo:     settings.AvailableTo,
		CameraAllowed:   settings.CameraAllowed,
		LocationAllowed: settings.LocationAllowed,
		DataUploadCloud: settings.DataUploadCloud,
	}
}

// DTOToParentSettings 将家长设置DTO转换为家长设置领域模型
func DTOToParentSettings(dto *ParentSettingsDTO) *domain.ParentSettings {
	if dto == nil {
		return nil
	}
	return &domain.ParentSettings{
		DailyLimitMin:   dto.DailyLimitMin,
		AvailableFrom:   dto.AvailableFrom,
		AvailableTo:     dto.AvailableTo,
		CameraAllowed:   dto.CameraAllowed,
		LocationAllowed: dto.LocationAllowed,
		DataUploadCloud: dto.DataUploadCloud,
	}
}

// RegisterRequestToUser 将注册请求DTO转换为用户领域模型
func RegisterRequestToUser(req *RegisterRequest) *domain.User {
	if req == nil {
		return nil
	}
	user := domain.NewUser(domain.UserTypeParent, req.Nickname)
	phone := req.Phone
	user.Phone = &phone
	return user
}

// AddChildRequestToUser 将添加儿童请求DTO转换为用户领域模型
func AddChildRequestToUser(req *AddChildRequest, parentID string) *domain.User {
	if req == nil {
		return nil
	}
	child := domain.NewUser(domain.UserTypeChild, req.Nickname)
	child.ParentID = &parentID
	age := req.Age
	child.Age = &age
	if req.Avatar != "" {
		child.Avatar = req.Avatar
	}
	return child
}
