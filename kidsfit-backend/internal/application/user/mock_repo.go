package user

import (
	"context"
	"errors"
	"sync"

	"github.com/kidsfit/api/internal/domain/user"
)

// 模拟仓储层返回的错误
var (
	// ErrMockNotFound 模拟记录未找到错误
	ErrMockNotFound = errors.New("record not found")
)

// MockUserRepository 用户仓储的Mock实现，用于单元测试
type MockUserRepository struct {
	mu    sync.RWMutex
	users map[string]*user.User
	byPhone map[string]*user.User
	createErr error
}

// NewMockUserRepository 创建Mock用户仓储实例
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:   make(map[string]*user.User),
		byPhone: make(map[string]*user.User),
	}
}

// SetCreateError 设置Create方法返回的错误，用于模拟创建失败场景
func (m *MockUserRepository) SetCreateError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.createErr = err
}

// Create 模拟创建用户，将用户存入内存map
func (m *MockUserRepository) Create(ctx context.Context, u *user.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.createErr != nil {
		return m.createErr
	}
	m.users[u.ID] = u
	if u.Phone != nil {
		m.byPhone[*u.Phone] = u
	}
	return nil
}

// GetByID 模拟根据ID获取用户
func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	u, ok := m.users[id]
	if !ok {
		return nil, ErrMockNotFound
	}
	return u, nil
}

// GetByPhone 模拟根据手机号获取用户
func (m *MockUserRepository) GetByPhone(ctx context.Context, phone string) (*user.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	u, ok := m.byPhone[phone]
	if !ok {
		return nil, ErrMockNotFound
	}
	return u, nil
}

// Update 模拟更新用户信息
func (m *MockUserRepository) Update(ctx context.Context, u *user.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.users[u.ID] = u
	if u.Phone != nil {
		m.byPhone[*u.Phone] = u
	}
	return nil
}

// Delete 模拟软删除用户
func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.users, id)
	return nil
}

// List 模拟查询用户列表
func (m *MockUserRepository) List(ctx context.Context, filter user.UserFilter, pagination user.Pagination) (*user.PaginatedResult[user.User], error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	items := make([]user.User, 0)
	for _, u := range m.users {
		items = append(items, *u)
	}
	return &user.PaginatedResult[user.User]{
		Items:    items,
		Total:    int64(len(items)),
		Page:     pagination.Page,
		PageSize: pagination.PageSize,
	}, nil
}

// MockFamilyRepository 家庭关系仓储的Mock实现，用于单元测试
type MockFamilyRepository struct {
	mu       sync.RWMutex
	families map[string]*user.Family
	byParent map[string][]*user.Family
}

// NewMockFamilyRepository 创建Mock家庭关系仓储实例
func NewMockFamilyRepository() *MockFamilyRepository {
	return &MockFamilyRepository{
		families: make(map[string]*user.Family),
		byParent: make(map[string][]*user.Family),
	}
}

// Create 模拟创建家庭关系
func (m *MockFamilyRepository) Create(ctx context.Context, f *user.Family) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.families[f.ID] = f
	m.byParent[f.ParentID] = append(m.byParent[f.ParentID], f)
	return nil
}

// GetByParentID 模拟根据家长ID获取家庭关系列表
func (m *MockFamilyRepository) GetByParentID(ctx context.Context, parentID string) ([]*user.Family, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.byParent[parentID], nil
}

// GetByChildID 模拟根据儿童ID获取家庭关系
func (m *MockFamilyRepository) GetByChildID(ctx context.Context, childID string) ([]*user.Family, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []*user.Family
	for _, f := range m.families {
		if f.ChildID == childID {
			result = append(result, f)
		}
	}
	return result, nil
}

// Delete 模拟删除家庭关系
func (m *MockFamilyRepository) Delete(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.families, id)
	return nil
}

// MockParentSettingsRepository 家长设置仓储的Mock实现，用于单元测试
type MockParentSettingsRepository struct {
	mu       sync.RWMutex
	settings map[string]*user.ParentSettings
}

// NewMockParentSettingsRepository 创建Mock家长设置仓储实例
func NewMockParentSettingsRepository() *MockParentSettingsRepository {
	return &MockParentSettingsRepository{
		settings: make(map[string]*user.ParentSettings),
	}
}

// Create 模拟创建家长设置
func (m *MockParentSettingsRepository) Create(ctx context.Context, s *user.ParentSettings) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.settings[s.ParentID] = s
	return nil
}

// GetByParentID 模拟根据家长ID获取设置
func (m *MockParentSettingsRepository) GetByParentID(ctx context.Context, parentID string) (*user.ParentSettings, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.settings[parentID]
	if !ok {
		return nil, ErrMockNotFound
	}
	return s, nil
}

// Update 模拟更新家长设置
func (m *MockParentSettingsRepository) Update(ctx context.Context, s *user.ParentSettings) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.settings[s.ParentID] = s
	return nil
}
