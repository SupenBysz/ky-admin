package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	internalapi "github.com/SupenBysz/ky-admin/internal/api"
	pkgapi "github.com/SupenBysz/ky-admin/pkg/api"
)

func TestHealthCheck(t *testing.T) {
	// 设置Gin测试模式
	gin.SetMode(gin.TestMode)

	// 创建路由器
	router := gin.New()
	router.Use(gin.Recovery())

	// 初始化健康检查处理器
	healthHandler := internalapi.NewHealthController()

	// 注册健康检查路由
	router.GET("/health", healthHandler.Health)
	router.GET("/livez", healthHandler.LivenessProbe)
	router.GET("/readyz", healthHandler.ReadinessProbe)

	// 测试健康检查
	t.Run("HealthCheck", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 验证响应状态码
		assert.Equal(t, http.StatusOK, w.Code)

		// 解析响应
		var resp pkgapi.Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		// 验证响应内容
		assert.Equal(t, 0, resp.Code)
		assert.Equal(t, "操作成功", resp.Message)
		assert.NotNil(t, resp.Data)

		// 验证健康检查数据
		data, ok := resp.Data.(map[string]interface{})
		require.True(t, ok)
		assert.Contains(t, data, "status")
		assert.Equal(t, "ok", data["status"])
	})

	// 测试存活检查
	t.Run("LivezCheck", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/livez", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 验证响应状态码
		assert.Equal(t, http.StatusOK, w.Code)

		// 解析响应
		var resp pkgapi.Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		// 验证响应内容
		assert.Equal(t, 0, resp.Code)
		assert.Equal(t, "操作成功", resp.Message)
		assert.NotNil(t, resp.Data)

		// 验证存活检查数据
		data, ok := resp.Data.(map[string]interface{})
		require.True(t, ok)
		assert.Contains(t, data, "status")
		assert.Equal(t, "alive", data["status"])
	})

	// 测试就绪检查
	t.Run("ReadyzCheck", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/readyz", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 验证响应状态码
		assert.Equal(t, http.StatusOK, w.Code)

		// 解析响应
		var resp pkgapi.Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		// 验证响应内容
		assert.Equal(t, 0, resp.Code)
		assert.Equal(t, "操作成功", resp.Message)
		assert.NotNil(t, resp.Data)

		// 验证就绪检查数据
		data, ok := resp.Data.(map[string]interface{})
		require.True(t, ok)
		assert.Contains(t, data, "status")
		assert.Equal(t, "ready", data["status"])
	})
}
