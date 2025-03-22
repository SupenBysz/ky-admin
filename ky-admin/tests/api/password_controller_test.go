package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/SupenBysz/ky-admin/internal/api/controller"
	"github.com/SupenBysz/ky-admin/internal/dto"
	"github.com/SupenBysz/ky-admin/internal/pkg/response"
	"github.com/SupenBysz/ky-admin/internal/service"
)

// 创建密码服务的模拟对象
type mockPasswordService struct {
	mock.Mock
}

func (m *mockPasswordService) ForgotPassword(dto *dto.ForgotPasswordDTO) error {
	args := m.Called(dto)
	return args.Error(0)
}

func (m *mockPasswordService) ResetPassword(dto *dto.ResetPasswordDTO) error {
	args := m.Called(dto)
	return args.Error(0)
}

func (m *mockPasswordService) ValidateResetToken(token string) (bool, error) {
	args := m.Called(token)
	return args.Bool(0), args.Error(1)
}

// PasswordControllerTestSuite 密码控制器测试套件
type PasswordControllerTestSuite struct {
	suite.Suite
	router       *gin.Engine
	mockService  *mockPasswordService
	passwordCtrl *controller.PasswordController
}

// SetupTest 在每个测试前设置测试环境
func (suite *PasswordControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	suite.mockService = new(mockPasswordService)
	suite.passwordCtrl = controller.NewPasswordController(suite.mockService)

	suite.router = gin.New()
	suite.router.Use(gin.Recovery())

	// 注册路由
	api := suite.router.Group("/api/v1")
	{
		password := api.Group("/password")
		{
			password.POST("/forgot", suite.passwordCtrl.ForgotPassword)
			password.GET("/validate-token", suite.passwordCtrl.ValidateResetToken)
			password.POST("/reset", suite.passwordCtrl.ResetPassword)
		}
	}
}

// TestForgotPassword_Success 测试忘记密码请求成功
func (suite *PasswordControllerTestSuite) TestForgotPassword_Success() {
	t := suite.T()

	// 准备请求数据
	reqData := dto.ForgotPasswordDTO{
		Email: "test@example.com",
	}
	jsonData, _ := json.Marshal(reqData)

	// 设置模拟服务的行为
	suite.mockService.On("ForgotPassword", mock.MatchedBy(
		func(dto *dto.ForgotPasswordDTO) bool {
			return dto.Email == reqData.Email
		})).Return(nil)

	// 创建请求
	req, _ := http.NewRequest("POST", "/api/v1/password/forgot", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// 处理请求
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// 断言响应
	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, 0, resp.Code)
	assert.Contains(t, resp.Message, "密码重置邮件已发送")

	// 验证模拟服务被正确调用
	suite.mockService.AssertExpectations(t)
}

// TestForgotPassword_InvalidEmail 测试提供无效邮箱的情况
func (suite *PasswordControllerTestSuite) TestForgotPassword_InvalidEmail() {
	t := suite.T()

	// 准备请求数据 - 无效邮箱格式
	reqData := map[string]string{
		"email": "invalid-email",
	}
	jsonData, _ := json.Marshal(reqData)

	// 创建请求
	req, _ := http.NewRequest("POST", "/api/v1/password/forgot", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// 处理请求
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// 断言响应 - 应返回400错误
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Message, "请求参数错误")
}

// TestForgotPassword_UserNotFound 测试用户不存在的情况
func (suite *PasswordControllerTestSuite) TestForgotPassword_UserNotFound() {
	t := suite.T()

	// 准备请求数据
	reqData := dto.ForgotPasswordDTO{
		Email: "nonexistent@example.com",
	}
	jsonData, _ := json.Marshal(reqData)

	// 设置模拟服务的行为 - 返回用户不存在错误
	suite.mockService.On("ForgotPassword", mock.MatchedBy(
		func(dto *dto.ForgotPasswordDTO) bool {
			return dto.Email == reqData.Email
		})).Return(service.ErrUserNotFound)

	// 创建请求
	req, _ := http.NewRequest("POST", "/api/v1/password/forgot", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// 处理请求
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// 断言响应 - 应返回500错误
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Contains(t, resp.Message, "请求密码重置失败")

	// 验证模拟服务被正确调用
	suite.mockService.AssertExpectations(t)
}

// TestValidateResetToken_Success 测试验证令牌成功的情况
func (suite *PasswordControllerTestSuite) TestValidateResetToken_Success() {
	t := suite.T()

	// 设置模拟服务的行为
	token := "valid-token"
	suite.mockService.On("ValidateResetToken", token).Return(true, nil)

	// 创建请求
	req, _ := http.NewRequest("GET", "/api/v1/password/validate-token?token="+token, nil)

	// 处理请求
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// 断言响应
	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, true, resp.Data)
	assert.Contains(t, resp.Message, "令牌验证成功")

	// 验证模拟服务被正确调用
	suite.mockService.AssertExpectations(t)
}

// TestValidateResetToken_Invalid 测试无效令牌的情况
func (suite *PasswordControllerTestSuite) TestValidateResetToken_Invalid() {
	t := suite.T()

	// 设置模拟服务的行为
	token := "invalid-token"
	suite.mockService.On("ValidateResetToken", token).Return(false, service.ErrInvalidResetToken)

	// 创建请求
	req, _ := http.NewRequest("GET", "/api/v1/password/validate-token?token="+token, nil)

	// 处理请求
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// 断言响应
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Message, "令牌验证失败")

	// 验证模拟服务被正确调用
	suite.mockService.AssertExpectations(t)
}

// TestValidateResetToken_Expired 测试过期令牌的情况
func (suite *PasswordControllerTestSuite) TestValidateResetToken_Expired() {
	t := suite.T()

	// 设置模拟服务的行为
	token := "expired-token"
	suite.mockService.On("ValidateResetToken", token).Return(false, service.ErrExpiredResetToken)

	// 创建请求
	req, _ := http.NewRequest("GET", "/api/v1/password/validate-token?token="+token, nil)

	// 处理请求
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// 断言响应
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Message, "令牌验证失败")
	assert.Contains(t, resp.Data.(string), "重置令牌已过期")

	// 验证模拟服务被正确调用
	suite.mockService.AssertExpectations(t)
}

// TestValidateResetToken_MissingToken 测试缺少令牌参数的情况
func (suite *PasswordControllerTestSuite) TestValidateResetToken_MissingToken() {
	t := suite.T()

	// 创建请求 - 不包含token参数
	req, _ := http.NewRequest("GET", "/api/v1/password/validate-token", nil)

	// 处理请求
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// 断言响应
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Message, "令牌不能为空")
}

// TestResetPassword_Success 测试重置密码成功的情况
func (suite *PasswordControllerTestSuite) TestResetPassword_Success() {
	t := suite.T()

	// 准备请求数据
	reqData := dto.ResetPasswordDTO{
		Token:       "valid-token",
		NewPassword: "newPassword123",
	}
	jsonData, _ := json.Marshal(reqData)

	// 设置模拟服务的行为
	suite.mockService.On("ResetPassword", mock.MatchedBy(
		func(dto *dto.ResetPasswordDTO) bool {
			return dto.Token == reqData.Token && dto.NewPassword == reqData.NewPassword
		})).Return(nil)

	// 创建请求
	req, _ := http.NewRequest("POST", "/api/v1/password/reset", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// 处理请求
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// 断言响应
	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, 0, resp.Code)
	assert.Contains(t, resp.Message, "密码已重置")

	// 验证模拟服务被正确调用
	suite.mockService.AssertExpectations(t)
}

// TestResetPassword_InvalidToken 测试使用无效令牌重置密码的情况
func (suite *PasswordControllerTestSuite) TestResetPassword_InvalidToken() {
	t := suite.T()

	// 准备请求数据
	reqData := dto.ResetPasswordDTO{
		Token:       "invalid-token",
		NewPassword: "newPassword123",
	}
	jsonData, _ := json.Marshal(reqData)

	// 设置模拟服务的行为
	suite.mockService.On("ResetPassword", mock.MatchedBy(
		func(dto *dto.ResetPasswordDTO) bool {
			return dto.Token == reqData.Token && dto.NewPassword == reqData.NewPassword
		})).Return(service.ErrInvalidResetToken)

	// 创建请求
	req, _ := http.NewRequest("POST", "/api/v1/password/reset", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// 处理请求
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// 断言响应
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Contains(t, resp.Message, "密码重置失败")
	assert.Contains(t, resp.Data.(string), "无效的重置令牌")

	// 验证模拟服务被正确调用
	suite.mockService.AssertExpectations(t)
}

// TestResetPassword_ExpiredToken 测试使用过期令牌重置密码的情况
func (suite *PasswordControllerTestSuite) TestResetPassword_ExpiredToken() {
	t := suite.T()

	// 准备请求数据
	reqData := dto.ResetPasswordDTO{
		Token:       "expired-token",
		NewPassword: "newPassword123",
	}
	jsonData, _ := json.Marshal(reqData)

	// 设置模拟服务的行为
	suite.mockService.On("ResetPassword", mock.MatchedBy(
		func(dto *dto.ResetPasswordDTO) bool {
			return dto.Token == reqData.Token && dto.NewPassword == reqData.NewPassword
		})).Return(service.ErrExpiredResetToken)

	// 创建请求
	req, _ := http.NewRequest("POST", "/api/v1/password/reset", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// 处理请求
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// 断言响应
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Contains(t, resp.Message, "密码重置失败")
	assert.Contains(t, resp.Data.(string), "重置令牌已过期")

	// 验证模拟服务被正确调用
	suite.mockService.AssertExpectations(t)
}

// TestResetPassword_UsedToken 测试使用已使用令牌重置密码的情况
func (suite *PasswordControllerTestSuite) TestResetPassword_UsedToken() {
	t := suite.T()

	// 准备请求数据
	reqData := dto.ResetPasswordDTO{
		Token:       "used-token",
		NewPassword: "newPassword123",
	}
	jsonData, _ := json.Marshal(reqData)

	// 设置模拟服务的行为
	suite.mockService.On("ResetPassword", mock.MatchedBy(
		func(dto *dto.ResetPasswordDTO) bool {
			return dto.Token == reqData.Token && dto.NewPassword == reqData.NewPassword
		})).Return(service.ErrTokenAlreadyUsed)

	// 创建请求
	req, _ := http.NewRequest("POST", "/api/v1/password/reset", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// 处理请求
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// 断言响应
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Contains(t, resp.Message, "密码重置失败")
	assert.Contains(t, resp.Data.(string), "令牌已被使用")

	// 验证模拟服务被正确调用
	suite.mockService.AssertExpectations(t)
}

// TestResetPassword_InvalidRequest 测试请求参数无效的情况
func (suite *PasswordControllerTestSuite) TestResetPassword_InvalidRequest() {
	t := suite.T()

	// 准备无效请求数据 - 缺少必要字段
	reqData := map[string]string{
		"token": "some-token",
		// 缺少new_password字段
	}
	jsonData, _ := json.Marshal(reqData)

	// 创建请求
	req, _ := http.NewRequest("POST", "/api/v1/password/reset", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// 处理请求
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// 断言响应
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Message, "请求参数错误")
}

// TestPasswordControllerTestSuite 运行测试套件
func TestPasswordControllerTestSuite(t *testing.T) {
	suite.Run(t, new(PasswordControllerTestSuite))
}
