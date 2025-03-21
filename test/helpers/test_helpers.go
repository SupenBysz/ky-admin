package helpers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/SupenBysz/ky-admin/pkg/config"
	"github.com/SupenBysz/ky-admin/pkg/logger"
)

// SetupTestEnvironment 初始化测试环境
func SetupTestEnvironment(t *testing.T) (*config.Config, *gorm.DB) {
	// 确保使用测试环境配置
	os.Setenv("KY_ADMIN_ENV", "test")

	// 获取项目根目录路径
	rootDir, err := findRootDir()
	require.NoError(t, err, "获取项目根目录失败")

	// 加载测试配置
	cfg, err := config.New(config.WithConfigPath(filepath.Join(rootDir, "configs")))
	require.NoError(t, err, "加载测试配置失败")
	require.NotNil(t, cfg, "获取配置失败")

	// 构建测试数据库DSN
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local",
		cfg.GetString("database.username"),
		cfg.GetString("database.password"),
		cfg.GetString("database.host"),
		cfg.GetInt("database.port"),
		cfg.GetString("database.name"),
		cfg.GetString("database.charset"),
	)

	// 连接测试数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	require.NoError(t, err, "连接测试数据库失败")
	require.NotNil(t, db, "获取数据库连接失败")

	return cfg, db
}

// CleanupTestEnvironment 清理测试环境
func CleanupTestEnvironment() {
	logger.Sync()
}

// GetTestTransaction 获取测试事务
func GetTestTransaction(t *testing.T, db *gorm.DB) *gorm.DB {
	tx := db.Begin()
	require.NoError(t, tx.Error, "开始测试事务失败")

	// 返回测试事务，调用方负责回滚
	return tx
}

// SetupControllerTest 设置控制器测试环境
func SetupControllerTest() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// CreateTestRequest 创建测试请求
func CreateTestRequest(method, url string, body string) (*http.Request, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, nil)
	if body != "" {
		req.Body = http.NoBody
	}
	return req, w
}

// WithTestContext 使用测试上下文执行函数
func WithTestContext(timeout time.Duration, fn func(ctx context.Context)) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	fn(ctx)
}

// ExecuteInTransaction 在事务中执行测试
func ExecuteInTransaction(t *testing.T, db *gorm.DB, fn func(tx *gorm.DB)) {
	// 开始事务
	tx := db.Begin()
	require.NoError(t, tx.Error, "开始事务失败")

	// 确保事务被回滚
	defer func() {
		tx.Rollback()
	}()

	// 执行测试函数
	fn(tx)
}

// MockHTTPResponse 模拟HTTP响应
func MockHTTPResponse(code int, body string) *http.Response {
	return &http.Response{
		StatusCode: code,
		Body:       nil, // 实际使用时应使用ioutil.NopCloser
	}
}

// GetTestDB 获取测试数据库连接
func GetTestDB(t *testing.T) *gorm.DB {
	_, db := SetupTestEnvironment(t)
	return db
}

// GetTestDBWithSQL 获取底层SQL数据库连接
func GetTestDBWithSQL(t *testing.T) (*gorm.DB, *sql.DB) {
	db := GetTestDB(t)
	sqlDB, err := db.DB()
	require.NoError(t, err, "获取SQL数据库连接失败")
	return db, sqlDB
}

// 查找项目根目录
func findRootDir() (string, error) {
	// 获取当前目录
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// 向上查找直到找到configs目录
	for {
		if _, err := os.Stat(filepath.Join(dir, "configs")); err == nil {
			return dir, nil
		}

		// 到达根目录，无法继续向上查找
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	// 如果当前不在项目目录中，尝试使用预设的路径
	baseDir := "/Volumes/DataDocument/CodeSpace/AI Solution/kysion.com/backend/ky-admin"
	if _, err := os.Stat(filepath.Join(baseDir, "configs")); err == nil {
		return baseDir, nil
	}

	return "", os.ErrNotExist
}
