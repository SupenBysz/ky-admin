// Package pure_api_test 提供纯API测试，不依赖项目内部代码
package pure_api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// 简单响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// 设置测试路由
func setupTestRouter() *gin.Engine {
	// 设置测试模式
	gin.SetMode(gin.TestMode)

	// 创建一个新的Gin引擎
	router := gin.New()
	router.Use(gin.Recovery())

	// 添加测试路由
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/api/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, Response{
			Code:    200,
			Message: "系统正常运行",
			Data:    gin.H{"status": "healthy"},
		})
	})

	router.POST("/api/data", func(c *gin.Context) {
		var data map[string]interface{}
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Code:    400,
				Message: "无效的请求数据",
				Data:    nil,
			})
			return
		}

		c.JSON(http.StatusCreated, Response{
			Code:    201,
			Message: "数据创建成功",
			Data:    data,
		})
	})

	router.GET("/api/protected", func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "Bearer valid-token" {
			c.JSON(http.StatusUnauthorized, Response{
				Code:    401,
				Message: "未授权访问",
				Data:    nil,
			})
			return
		}

		c.JSON(http.StatusOK, Response{
			Code:    200,
			Message: "访问成功",
			Data:    gin.H{"user": "测试用户"},
		})
	})

	return router
}

// TestPingRoute 测试简单的ping路由
func TestPingRoute(t *testing.T) {
	router := setupTestRouter()

	// 创建一个测试请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	// 断言返回状态码为200
	assert.Equal(t, http.StatusOK, w.Code)
	// 断言返回的JSON包含预期的消息
	assert.Contains(t, w.Body.String(), "pong")
}

// TestStatusAPI 测试状态API
func TestStatusAPI(t *testing.T) {
	router := setupTestRouter()

	// 创建一个测试请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/status", nil)
	router.ServeHTTP(w, req)

	// 断言返回状态码为200
	assert.Equal(t, http.StatusOK, w.Code)

	// 解析响应
	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 验证响应内容
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "系统正常运行", response.Message)

	// 验证data字段
	data, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "healthy", data["status"])
}

// TestPostData 测试POST数据
func TestPostData(t *testing.T) {
	router := setupTestRouter()

	// 准备请求数据
	requestBody := `{"name": "测试名称", "value": 123}`

	// 创建一个测试请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/data", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// 断言返回状态码为201
	assert.Equal(t, http.StatusCreated, w.Code)

	// 解析响应
	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 验证响应内容
	assert.Equal(t, 201, response.Code)
	assert.Equal(t, "数据创建成功", response.Message)

	// 验证返回的数据
	data, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "测试名称", data["name"])
	assert.Equal(t, float64(123), data["value"])
}

// TestProtectedRoute 测试需要认证的路由
func TestProtectedRoute(t *testing.T) {
	router := setupTestRouter()

	// 测试未授权访问
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/api/protected", nil)
	router.ServeHTTP(w1, req1)

	// 断言返回状态码为401
	assert.Equal(t, http.StatusUnauthorized, w1.Code)

	// 解析响应
	var response1 Response
	err := json.Unmarshal(w1.Body.Bytes(), &response1)
	assert.NoError(t, err)
	assert.Equal(t, 401, response1.Code)
	assert.Equal(t, "未授权访问", response1.Message)

	// 测试授权访问
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/api/protected", nil)
	req2.Header.Set("Authorization", "Bearer valid-token")
	router.ServeHTTP(w2, req2)

	// 断言返回状态码为200
	assert.Equal(t, http.StatusOK, w2.Code)

	// 解析响应
	var response2 Response
	err = json.Unmarshal(w2.Body.Bytes(), &response2)
	assert.NoError(t, err)
	assert.Equal(t, 200, response2.Code)
	assert.Equal(t, "访问成功", response2.Message)

	// 验证返回的数据
	data, ok := response2.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "测试用户", data["user"])
}
