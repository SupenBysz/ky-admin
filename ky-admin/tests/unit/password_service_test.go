package unit

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/SupenBysz/ky-admin/internal/dto"
	"github.com/SupenBysz/ky-admin/internal/model"
	"github.com/SupenBysz/ky-admin/internal/service"
)

// 创建模拟的用户仓库
type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) FindByID(id uint) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *mockUserRepository) FindByUsername(username string) (*model.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *mockUserRepository) FindByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *mockUserRepository) FindAll(page, pageSize int) ([]model.User, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]model.User), args.Get(1).(int64), args.Error(2)
}

func (m *mockUserRepository) Create(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *mockUserRepository) Update(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *mockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockUserRepository) FindByPage(query *dto.UserPageQueryDTO) ([]model.User, int64, error) {
	args := m.Called(query)
	return args.Get(0).([]model.User), args.Get(1).(int64), args.Error(2)
}

func (m *mockUserRepository) AssignRoles(userID uint, roleIDs []uint) error {
	args := m.Called(userID, roleIDs)
	return args.Error(0)
}

func (m *mockUserRepository) FindUserRoles(userID uint) ([]model.Role, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.Role), args.Error(1)
}

// 创建模拟的密码重置仓库
type mockPasswordResetRepository struct {
	mock.Mock
}

func (m *mockPasswordResetRepository) Create(userID uint, email, token string, expiration time.Time) (*model.PasswordReset, error) {
	args := m.Called(userID, email, token, expiration)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.PasswordReset), args.Error(1)
}

func (m *mockPasswordResetRepository) FindByToken(token string) (*model.PasswordReset, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.PasswordReset), args.Error(1)
}

func (m *mockPasswordResetRepository) MarkTokenAsUsed(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *mockPasswordResetRepository) DeleteExpired() error {
	args := m.Called()
	return args.Error(0)
}

// 创建模拟的邮件服务
type mockEmailService struct {
	mock.Mock
}

func (m *mockEmailService) SendEmail(to []string, subject, content string) error {
	args := m.Called(to, subject, content)
	return args.Error(0)
}

func (m *mockEmailService) SendHTMLEmail(to []string, subject, htmlContent string) error {
	args := m.Called(to, subject, htmlContent)
	return args.Error(0)
}

func (m *mockEmailService) SendTemplateEmail(to []string, subject, templateName string, data interface{}) error {
	args := m.Called(to, subject, templateName, data)
	return args.Error(0)
}

func (m *mockEmailService) SendPasswordResetEmail(to string, resetLink string, username string) error {
	args := m.Called(to, resetLink, username)
	return args.Error(0)
}

// PasswordServiceTestSuite 密码服务测试套件
type PasswordServiceTestSuite struct {
	suite.Suite
	mockUserRepo          *mockUserRepository
	mockPasswordResetRepo *mockPasswordResetRepository
	mockEmailSvc          *mockEmailService
	passwordService       service.PasswordService
}

// SetupTest 在每个测试前设置测试环境
func (suite *PasswordServiceTestSuite) SetupTest() {
	suite.mockUserRepo = new(mockUserRepository)
	suite.mockPasswordResetRepo = new(mockPasswordResetRepository)
	suite.mockEmailSvc = new(mockEmailService)

	// 创建服务实例，注入模拟对象
	suite.passwordService = service.NewPasswordService(
		suite.mockUserRepo,
		suite.mockPasswordResetRepo,
		suite.mockEmailSvc,
		"http://example.com",
	)
}

// TestForgotPassword_Success 测试忘记密码成功的情况
func (suite *PasswordServiceTestSuite) TestForgotPassword_Success() {
	t := suite.T()

	// 1. 设置测试数据
	email := "test@example.com"
	user := &model.User{
		ID:       1,
		Username: "testuser",
		Email:    email,
	}
	resetToken := &model.PasswordReset{
		ID:        1,
		UserID:    user.ID,
		Email:     email,
		Token:     "valid-token",
		Used:      false,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	// 2. 设置模拟行为
	suite.mockUserRepo.On("FindByEmail", email).Return(user, nil)
	suite.mockPasswordResetRepo.On("Create",
		user.ID, email, mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).
		Return(resetToken, nil)
	suite.mockEmailSvc.On("SendPasswordResetEmail",
		email, mock.AnythingOfType("string"), user.Username).
		Return(nil)

	// 3. 执行测试
	err := suite.passwordService.ForgotPassword(&dto.ForgotPasswordDTO{Email: email})

	// 4. 断言
	assert.NoError(t, err)
	suite.mockUserRepo.AssertExpectations(t)
	suite.mockPasswordResetRepo.AssertExpectations(t)
	suite.mockEmailSvc.AssertExpectations(t)
}

// TestForgotPassword_UserNotFound 测试邮箱不存在的情况
func (suite *PasswordServiceTestSuite) TestForgotPassword_UserNotFound() {
	t := suite.T()

	// 1. 设置测试数据
	email := "nonexistent@example.com"

	// 2. 设置模拟行为
	suite.mockUserRepo.On("FindByEmail", email).Return(nil, errors.New("用户不存在"))

	// 3. 执行测试
	err := suite.passwordService.ForgotPassword(&dto.ForgotPasswordDTO{Email: email})

	// 4. 断言
	assert.Error(t, err)
	assert.Equal(t, service.ErrUserNotFound, err)
	suite.mockUserRepo.AssertExpectations(t)
}

// TestValidateResetToken_Success 测试验证重置令牌成功的情况
func (suite *PasswordServiceTestSuite) TestValidateResetToken_Success() {
	t := suite.T()

	// 1. 设置测试数据
	token := "valid-token"
	resetToken := &model.PasswordReset{
		ID:        1,
		UserID:    1,
		Email:     "test@example.com",
		Token:     token,
		Used:      false,
		ExpiresAt: time.Now().Add(1 * time.Hour), // 未过期
	}

	// 2. 设置模拟行为
	suite.mockPasswordResetRepo.On("FindByToken", token).Return(resetToken, nil)

	// 3. 执行测试
	valid, err := suite.passwordService.ValidateResetToken(token)

	// 4. 断言
	assert.NoError(t, err)
	assert.True(t, valid)
	suite.mockPasswordResetRepo.AssertExpectations(t)
}

// TestValidateResetToken_InvalidToken 测试无效令牌的情况
func (suite *PasswordServiceTestSuite) TestValidateResetToken_InvalidToken() {
	t := suite.T()

	// 1. 设置测试数据
	token := "invalid-token"

	// 2. 设置模拟行为
	suite.mockPasswordResetRepo.On("FindByToken", token).Return(nil, errors.New("令牌不存在"))

	// 3. 执行测试
	valid, err := suite.passwordService.ValidateResetToken(token)

	// 4. 断言
	assert.Error(t, err)
	assert.False(t, valid)
	suite.mockPasswordResetRepo.AssertExpectations(t)
}

// TestValidateResetToken_ExpiredToken 测试过期令牌的情况
func (suite *PasswordServiceTestSuite) TestValidateResetToken_ExpiredToken() {
	t := suite.T()

	// 1. 设置测试数据
	token := "expired-token"
	resetToken := &model.PasswordReset{
		ID:        1,
		UserID:    1,
		Email:     "test@example.com",
		Token:     token,
		Used:      false,
		ExpiresAt: time.Now().Add(-1 * time.Hour), // 已过期
	}

	// 2. 设置模拟行为
	suite.mockPasswordResetRepo.On("FindByToken", token).Return(resetToken, nil)

	// 3. 执行测试
	valid, err := suite.passwordService.ValidateResetToken(token)

	// 4. 断言
	assert.Error(t, err)
	assert.False(t, valid)
	assert.Equal(t, service.ErrExpiredResetToken, err)
	suite.mockPasswordResetRepo.AssertExpectations(t)
}

// TestValidateResetToken_UsedToken 测试已使用令牌的情况
func (suite *PasswordServiceTestSuite) TestValidateResetToken_UsedToken() {
	t := suite.T()

	// 1. 设置测试数据
	token := "used-token"
	resetToken := &model.PasswordReset{
		ID:        1,
		UserID:    1,
		Email:     "test@example.com",
		Token:     token,
		Used:      true, // 已使用
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	// 2. 设置模拟行为
	suite.mockPasswordResetRepo.On("FindByToken", token).Return(resetToken, nil)

	// 3. 执行测试
	valid, err := suite.passwordService.ValidateResetToken(token)

	// 4. 断言
	assert.Error(t, err)
	assert.False(t, valid)
	assert.Equal(t, service.ErrTokenAlreadyUsed, err)
	suite.mockPasswordResetRepo.AssertExpectations(t)
}

// TestResetPassword_Success 测试重置密码成功的情况
func (suite *PasswordServiceTestSuite) TestResetPassword_Success() {
	t := suite.T()

	// 1. 设置测试数据
	token := "valid-token"
	newPassword := "newPassword123"
	resetToken := &model.PasswordReset{
		ID:        1,
		UserID:    1,
		Email:     "test@example.com",
		Token:     token,
		Used:      false,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	user := &model.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Password: "oldhashedpassword",
	}

	// 2. 设置模拟行为
	suite.mockPasswordResetRepo.On("FindByToken", token).Return(resetToken, nil)
	suite.mockUserRepo.On("FindByID", uint(1)).Return(user, nil)
	suite.mockUserRepo.On("Update", mock.AnythingOfType("*model.User")).Return(nil)
	suite.mockPasswordResetRepo.On("MarkTokenAsUsed", token).Return(nil)

	// 3. 执行测试
	err := suite.passwordService.ResetPassword(&dto.ResetPasswordDTO{
		Token:       token,
		NewPassword: newPassword,
	})

	// 4. 断言
	assert.NoError(t, err)
	suite.mockPasswordResetRepo.AssertExpectations(t)
	suite.mockUserRepo.AssertExpectations(t)
}

// TestResetPassword_InvalidToken 测试无效令牌重置密码的情况
func (suite *PasswordServiceTestSuite) TestResetPassword_InvalidToken() {
	t := suite.T()

	// 1. 设置测试数据
	token := "invalid-token"
	newPassword := "newPassword123"

	// 2. 设置模拟行为
	suite.mockPasswordResetRepo.On("FindByToken", token).Return(nil, errors.New("令牌不存在"))

	// 3. 执行测试
	err := suite.passwordService.ResetPassword(&dto.ResetPasswordDTO{
		Token:       token,
		NewPassword: newPassword,
	})

	// 4. 断言
	assert.Error(t, err)
	assert.Equal(t, service.ErrInvalidResetToken, err)
	suite.mockPasswordResetRepo.AssertExpectations(t)
}

// TestResetPassword_ExpiredToken 测试使用过期令牌重置密码的情况
func (suite *PasswordServiceTestSuite) TestResetPassword_ExpiredToken() {
	t := suite.T()

	// 1. 设置测试数据
	token := "expired-token"
	newPassword := "newPassword123"
	resetToken := &model.PasswordReset{
		ID:        1,
		UserID:    1,
		Email:     "test@example.com",
		Token:     token,
		Used:      false,
		ExpiresAt: time.Now().Add(-1 * time.Hour), // 已过期
	}

	// 2. 设置模拟行为
	suite.mockPasswordResetRepo.On("FindByToken", token).Return(resetToken, nil)

	// 3. 执行测试
	err := suite.passwordService.ResetPassword(&dto.ResetPasswordDTO{
		Token:       token,
		NewPassword: newPassword,
	})

	// 4. 断言
	assert.Error(t, err)
	assert.Equal(t, service.ErrExpiredResetToken, err)
	suite.mockPasswordResetRepo.AssertExpectations(t)
}

// TestResetPassword_UsedToken 测试使用已使用令牌重置密码的情况
func (suite *PasswordServiceTestSuite) TestResetPassword_UsedToken() {
	t := suite.T()

	// 1. 设置测试数据
	token := "used-token"
	newPassword := "newPassword123"
	resetToken := &model.PasswordReset{
		ID:        1,
		UserID:    1,
		Email:     "test@example.com",
		Token:     token,
		Used:      true, // 已使用
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	// 2. 设置模拟行为
	suite.mockPasswordResetRepo.On("FindByToken", token).Return(resetToken, nil)

	// 3. 执行测试
	err := suite.passwordService.ResetPassword(&dto.ResetPasswordDTO{
		Token:       token,
		NewPassword: newPassword,
	})

	// 4. 断言
	assert.Error(t, err)
	assert.Equal(t, service.ErrTokenAlreadyUsed, err)
	suite.mockPasswordResetRepo.AssertExpectations(t)
}

// TestPasswordServiceTestSuite 运行测试套件
func TestPasswordServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PasswordServiceTestSuite))
}
