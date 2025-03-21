// Package api_test 包含API测试
package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestPingRoute 测试服务器状态接口
func TestPingRoute(t *testing.T) {
	// 设置测试模式
	gin.SetMode(gin.TestMode)

	// 创建一个新的gin引擎
	router := gin.New()

	// 添加一个简单的ping路由
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// 创建一个测试请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	// 断言返回状态码为200
	assert.Equal(t, http.StatusOK, w.Code)
	// 断言返回的JSON包含预期的消息
	assert.Contains(t, w.Body.String(), "pong")
}
