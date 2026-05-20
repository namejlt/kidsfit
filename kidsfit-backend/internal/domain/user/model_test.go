package user

import (
	"testing"
)

// TestNewUser 测试创建家长和儿童用户
func TestNewUser(t *testing.T) {
	t.Run("创建家长用户", func(t *testing.T) {
		u := NewUser(UserTypeParent, "测试家长")
		if u.ID == "" {
			t.Error("用户ID不应为空")
		}
		if u.Type != UserTypeParent {
			t.Errorf("用户类型应为 parent，实际为 %s", u.Type)
		}
		if u.Nickname != "测试家长" {
			t.Errorf("昵称应为 测试家长，实际为 %s", u.Nickname)
		}
		if u.Status != UserStatusActive {
			t.Errorf("状态应为 active，实际为 %s", u.Status)
		}
		if u.CreatedAt.IsZero() {
			t.Error("创建时间不应为零值")
		}
		if u.UpdatedAt.IsZero() {
			t.Error("更新时间不应为零值")
		}
	})

	t.Run("创建儿童用户", func(t *testing.T) {
		u := NewUser(UserTypeChild, "测试儿童")
		if u.Type != UserTypeChild {
			t.Errorf("用户类型应为 child，实际为 %s", u.Type)
		}
		if u.Nickname != "测试儿童" {
			t.Errorf("昵称应为 测试儿童，实际为 %s", u.Nickname)
		}
	})
}

// TestUser_IsParent 测试用户是否为家长类型判断
func TestUser_IsParent(t *testing.T) {
	parent := NewUser(UserTypeParent, "家长")
	if !parent.IsParent() {
		t.Error("家长用户 IsParent() 应返回 true")
	}

	child := NewUser(UserTypeChild, "儿童")
	if child.IsParent() {
		t.Error("儿童用户 IsParent() 应返回 false")
	}
}

// TestUser_IsChild 测试用户是否为儿童类型判断
func TestUser_IsChild(t *testing.T) {
	child := NewUser(UserTypeChild, "儿童")
	if !child.IsChild() {
		t.Error("儿童用户 IsChild() 应返回 true")
	}

	parent := NewUser(UserTypeParent, "家长")
	if parent.IsChild() {
		t.Error("家长用户 IsChild() 应返回 false")
	}
}

// TestUser_IsActive 测试用户是否处于活跃状态判断
func TestUser_IsActive(t *testing.T) {
	u := NewUser(UserTypeParent, "用户")
	if !u.IsActive() {
		t.Error("新建用户 IsActive() 应返回 true")
	}

	u.Status = UserStatusInactive
	if u.IsActive() {
		t.Error("未激活用户 IsActive() 应返回 false")
	}

	u.Status = UserStatusDeleted
	if u.IsActive() {
		t.Error("已删除用户 IsActive() 应返回 false")
	}
}

// TestUser_ValidateAge 测试年龄校验
func TestUser_ValidateAge(t *testing.T) {
	t.Run("年龄为nil时不报错", func(t *testing.T) {
		u := NewUser(UserTypeChild, "儿童")
		if err := u.ValidateAge(); err != nil {
			t.Errorf("年龄为nil时应返回nil，实际返回 %v", err)
		}
	})

	t.Run("有效年龄3-12", func(t *testing.T) {
		validAges := []int{3, 6, 9, 12}
		for _, age := range validAges {
			u := NewUser(UserTypeChild, "儿童")
			u.Age = &[]int{age}[0]
			if err := u.ValidateAge(); err != nil {
				t.Errorf("年龄 %d 应为有效值，实际返回错误: %v", age, err)
			}
		}
	})

	t.Run("无效年龄2和13", func(t *testing.T) {
		invalidAges := []int{2, 13}
		for _, age := range invalidAges {
			u := NewUser(UserTypeChild, "儿童")
			u.Age = &[]int{age}[0]
			if err := u.ValidateAge(); err == nil {
				t.Errorf("年龄 %d 应为无效值，但未返回错误", age)
			}
		}
	})
}

// TestUser_GetAgeGroup 测试年龄组分类
func TestUser_GetAgeGroup(t *testing.T) {
	t.Run("年龄为nil时返回空字符串", func(t *testing.T) {
		u := NewUser(UserTypeChild, "儿童")
		if got := u.GetAgeGroup(); got != "" {
			t.Errorf("年龄为nil时应返回空字符串，实际返回 %s", got)
		}
	})

	t.Run("3-6岁年龄组", func(t *testing.T) {
		for _, age := range []int{3, 4, 5, 6} {
			u := NewUser(UserTypeChild, "儿童")
			u.Age = &[]int{age}[0]
			if got := u.GetAgeGroup(); got != "3-6" {
				t.Errorf("年龄 %d 应属于 3-6 组，实际为 %s", age, got)
			}
		}
	})

	t.Run("7-9岁年龄组", func(t *testing.T) {
		for _, age := range []int{7, 8, 9} {
			u := NewUser(UserTypeChild, "儿童")
			u.Age = &[]int{age}[0]
			if got := u.GetAgeGroup(); got != "7-9" {
				t.Errorf("年龄 %d 应属于 7-9 组，实际为 %s", age, got)
			}
		}
	})

	t.Run("10-12岁年龄组", func(t *testing.T) {
		for _, age := range []int{10, 11, 12} {
			u := NewUser(UserTypeChild, "儿童")
			u.Age = &[]int{age}[0]
			if got := u.GetAgeGroup(); got != "10-12" {
				t.Errorf("年龄 %d 应属于 10-12 组，实际为 %s", age, got)
			}
		}
	})

	t.Run("超出范围的年龄返回空字符串", func(t *testing.T) {
		for _, age := range []int{2, 13} {
			u := NewUser(UserTypeChild, "儿童")
			u.Age = &[]int{age}[0]
			if got := u.GetAgeGroup(); got != "" {
				t.Errorf("年龄 %d 应返回空字符串，实际为 %s", age, got)
			}
		}
	})
}

// TestNewFamily 测试创建家庭关系
func TestNewFamily(t *testing.T) {
	parentID := "parent-123"
	childID := "child-456"
	relation := RelationFather

	family := NewFamily(parentID, childID, relation)
	if family.ID == "" {
		t.Error("家庭关系ID不应为空")
	}
	if family.ParentID != parentID {
		t.Errorf("家长ID应为 %s，实际为 %s", parentID, family.ParentID)
	}
	if family.ChildID != childID {
		t.Errorf("儿童ID应为 %s，实际为 %s", childID, family.ChildID)
	}
	if family.Relation != relation {
		t.Errorf("关系类型应为 %s，实际为 %s", relation, family.Relation)
	}
	if family.CreatedAt.IsZero() {
		t.Error("创建时间不应为零值")
	}
}

// TestNewParentSettings 测试创建家长设置，验证默认值
func TestNewParentSettings(t *testing.T) {
	parentID := "parent-123"

	settings := NewParentSettings(parentID)
	if settings.ID == "" {
		t.Error("设置ID不应为空")
	}
	if settings.ParentID != parentID {
		t.Errorf("家长ID应为 %s，实际为 %s", parentID, settings.ParentID)
	}
	if settings.DailyLimitMin != 30 {
		t.Errorf("默认每日时长限制应为 30 分钟，实际为 %d", settings.DailyLimitMin)
	}
	if settings.CreatedAt.IsZero() {
		t.Error("创建时间不应为零值")
	}
	if settings.UpdatedAt.IsZero() {
		t.Error("更新时间不应为零值")
	}
}
