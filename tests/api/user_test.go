package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SupenBysz/ky-admin/internal/api"
	pkgapi "github.com/SupenBysz/ky-admin/pkg/api"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUserController_Login(t *testing.T) {
	// 设置测试模式
	gin.SetMode(gin.TestMode)

	// 创建用户控制器
	userController := api.NewUserController()

	// 定义测试用例
	tests := []struct {
		name           string
		requestBody    api.LoginRequest
		expectedStatus int
		expectedCode   int
		shouldHaveData bool
	}{
		{
			name: "成功登录",
			requestBody: api.LoginRequest{
				Username: "admin",
				Password: "password123",
			},
			expectedStatus: http.StatusOK,
			expectedCode:   pkgapi.CodeSuccess,
			shouldHaveData: true,
		},
		{
			name: "用户名为空",
			requestBody: api.LoginRequest{
				Username: "",
				Password: "password123",
			},
			expectedStatus: http.StatusOK,
			expectedCode:   pkgapi.CodeParamError,
			shouldHaveData: false,
		},
		{
			name: "密码为空",
			requestBody: api.LoginRequest{
				Username: "admin",
				Password: "",
			},
			expectedStatus: http.StatusOK,
			expectedCode:   pkgapi.CodeParamError,
			shouldHaveData: false,
		},
	}

	// 运行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建新的记录器和上下文
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// 序列化请求体
			jsonData, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			// 调用登录处理函数
			userController.Login(c)

			// 验证响应状态码
			assert.Equal(t, tt.expectedStatus, w.Code)

			// 解析响应体
			var resp pkgapi.Response
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)

			// 验证响应结果
			assert.Equal(t, tt.expectedCode, resp.Code)

			if tt.shouldHaveData {
				// 验证响应数据
				data, ok := resp.Data.(map[string]interface{})
				assert.True(t, ok)
				assert.Contains(t, data, "token")
				assert.Contains(t, data, "userId")
			}
		})
	}
}

func TestUserController_GetUsers(t *testing.T) {
	// 设置测试模式
	gin.SetMode(gin.TestMode)

	// 创建用户控制器
	userController := api.NewUserController()

	// 创建新的记录器和上下文
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 创建请求
	req := httptest.NewRequest("GET", "/api/users", nil)
	c.Request = req

	// 调用获取用户列表处理函数
	userController.GetUsers(c)

	// 验证响应状态码
	assert.Equal(t, http.StatusOK, w.Code)

	// 解析响应体
	var resp pkgapi.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	// 验证响应结果
	assert.Equal(t, pkgapi.CodeSuccess, resp.Code)

	// 验证响应数据
	users, ok := resp.Data.([]interface{})
	assert.True(t, ok)
	assert.Len(t, users, 2) // 预期有两个用户

	// 验证第一个用户的数据
	user1, ok := users[0].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, float64(1), user1["id"])
	assert.Equal(t, "admin", user1["username"])
	assert.Equal(t, "管理员", user1["nickname"])
	assert.Equal(t, "admin@example.com", user1["email"])
}

// 性能测试示例
func BenchmarkUserController_Login(b *testing.B) {
	// 设置测试模式
	gin.SetMode(gin.TestMode)

	// 创建用户控制器
	userController := api.NewUserController()

	// 准备请求数据
	requestBody := api.LoginRequest{
		Username: "admin",
		Password: "password123",
	}
	jsonData, _ := json.Marshal(requestBody)

	// 运行性能测试
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		userController.Login(c)
	}
}
