package unit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/SupenBysz/ky-admin/internal/model"
	"github.com/SupenBysz/ky-admin/internal/service"
	"github.com/SupenBysz/ky-admin/tests/mocks"
)

// UserServiceTestSuite 是用户服务测试套件
type UserServiceTestSuite struct {
	suite.Suite
	mockUserRepo *mocks.UserRepository
	userService  service.UserService
}

// SetupSuite 在所有测试前设置测试环境
func (suite *UserServiceTestSuite) SetupSuite() {
	// 创建模拟对象
	suite.mockUserRepo = new(mocks.UserRepository)

	// 创建服务实例，注入模拟对象
	suite.userService = service.NewUserService(suite.mockUserRepo)
}

// TestGetUsersByPage 测试分页获取用户
func (suite *UserServiceTestSuite) TestGetUsersByPage() {
	t := suite.T()

	// 设置模拟预期
	mockUsers := []model.User{
		{
			ID:        1,
			Username:  "admin",
			Nickname:  "系统管理员",
			Email:     "admin@example.com",
			Phone:     "13800138000",
			Status:    1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			Username:  "user1",
			Nickname:  "测试用户1",
			Email:     "user1@example.com",
			Phone:     "13800138001",
			Status:    1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	var total int64 = 2
	page := 1
	pageSize := 10

	// 设置模拟行为
	suite.mockUserRepo.On("FindAll", page, pageSize).Return(mockUsers, total, nil)

	// 执行被测试函数
	users, count, err := suite.userService.GetUsers(page, pageSize)

	// 断言结果
	assert.NoError(t, err)
	assert.Equal(t, total, count)
	assert.Len(t, users, 2)
	assert.Equal(t, uint(1), users[0].ID)
	assert.Equal(t, "admin", users[0].Username)
	assert.Equal(t, uint(2), users[1].ID)
	assert.Equal(t, "user1", users[1].Username)

	// 验证模拟对象的方法是否被调用
	suite.mockUserRepo.AssertExpectations(t)
}

// TestGetUserByID 测试根据ID获取用户
func (suite *UserServiceTestSuite) TestGetUserByID() {
	t := suite.T()

	// 设置模拟数据
	mockUser := model.User{
		ID:        1,
		Username:  "admin",
		Nickname:  "系统管理员",
		Email:     "admin@example.com",
		Phone:     "13800138000",
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 设置模拟行为
	suite.mockUserRepo.On("FindByID", uint(1)).Return(&mockUser, nil)

	// 执行被测试函数
	user, err := suite.userService.GetByID(1)

	// 断言结果
	assert.NoError(t, err)
	assert.Equal(t, uint(1), user.ID)
	assert.Equal(t, "admin", user.Username)
	assert.Equal(t, "系统管理员", user.Nickname)

	// 验证模拟对象的方法是否被调用
	suite.mockUserRepo.AssertExpectations(t)
}

// TestCreateUser 测试创建用户
func (suite *UserServiceTestSuite) TestCreateUser() {
	t := suite.T()

	// 准备测试数据
	newUser := &model.User{
		Username: "newuser",
		Password: "Test@123",
		Nickname: "新用户",
		Email:    "newuser@example.com",
		Phone:    "13800138002",
		Status:   1,
	}

	// 设置模拟数据
	createdUser := *newUser
	createdUser.ID = 3
	createdUser.CreatedAt = time.Now()
	createdUser.UpdatedAt = time.Now()

	// 设置模拟行为
	suite.mockUserRepo.On("Create", mock.AnythingOfType("*model.User")).Return(nil)

	// 执行被测试函数
	err := suite.userService.CreateUser(newUser)

	// 断言结果
	assert.NoError(t, err)

	// 验证模拟对象的方法是否被调用
	suite.mockUserRepo.AssertExpectations(t)
}

// TestUpdateUser 测试更新用户
func (suite *UserServiceTestSuite) TestUpdateUser() {
	t := suite.T()

	// 设置模拟数据
	originalUser := model.User{
		ID:        1,
		Username:  "admin",
		Nickname:  "系统管理员",
		Email:     "admin@example.com",
		Phone:     "13800138000",
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	updatedUser := originalUser
	updatedUser.Nickname = "更新的昵称"
	updatedUser.Email = "updated@example.com"
	updatedUser.Phone = "13900139000"

	// 设置模拟行为
	suite.mockUserRepo.On("FindByID", uint(1)).Return(&originalUser, nil)
	suite.mockUserRepo.On("Update", mock.AnythingOfType("*model.User")).Return(nil)

	// 执行被测试函数 - 直接使用更新后的用户对象
	err := suite.userService.UpdateUser(&updatedUser)

	// 断言结果
	assert.NoError(t, err)

	// 验证模拟对象的方法是否被调用
	suite.mockUserRepo.AssertExpectations(t)
}

// TestDeleteUser 测试删除用户
func (suite *UserServiceTestSuite) TestDeleteUser() {
	t := suite.T()

	// 设置模拟行为
	suite.mockUserRepo.On("Delete", uint(1)).Return(nil)

	// 执行被测试函数
	err := suite.userService.DeleteUser(1)

	// 断言结果
	assert.NoError(t, err)

	// 验证模拟对象的方法是否被调用
	suite.mockUserRepo.AssertExpectations(t)
}

// TestUserServiceTestSuite 运行用户服务测试套件
func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
