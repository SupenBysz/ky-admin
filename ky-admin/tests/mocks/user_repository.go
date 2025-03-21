package mocks

import (
	"github.com/SupenBysz/ky-admin/internal/dto"
	"github.com/SupenBysz/ky-admin/internal/model"
	"github.com/stretchr/testify/mock"
)

// UserRepository 用户仓库模拟
type UserRepository struct {
	mock.Mock
}

// FindByID 根据ID查找用户
func (m *UserRepository) FindByID(id uint) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

// FindByUsername 根据用户名查找用户
func (m *UserRepository) FindByUsername(username string) (*model.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

// FindAll 查找所有用户
func (m *UserRepository) FindAll(page, pageSize int) ([]model.User, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]model.User), args.Get(1).(int64), args.Error(2)
}

// Create 创建用户
func (m *UserRepository) Create(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// Update 更新用户
func (m *UserRepository) Update(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// Delete 删除用户
func (m *UserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// FindByEmail 根据邮箱查找用户
func (m *UserRepository) FindByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

// FindByPage 分页查询用户
func (m *UserRepository) FindByPage(query *dto.UserPageQueryDTO) ([]model.User, int64, error) {
	args := m.Called(query)
	return args.Get(0).([]model.User), args.Get(1).(int64), args.Error(2)
}

// AssignRoles 分配角色
func (m *UserRepository) AssignRoles(userID uint, roleIDs []uint) error {
	args := m.Called(userID, roleIDs)
	return args.Error(0)
}

// FindUserRoles 查找用户角色
func (m *UserRepository) FindUserRoles(userID uint) ([]model.Role, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.Role), args.Error(1)
}
