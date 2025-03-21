package api_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/SupenBysz/ky-admin/test/helpers"
)

func TestDatabaseConnection(t *testing.T) {
	// 不再跳过测试
	// t.Skip("需要真实数据库连接")

	// 获取项目根目录路径
	_, err := findRootDir()
	assert.NoError(t, err, "查找项目根目录失败")

	// 设置测试环境
	os.Setenv("KY_ADMIN_ENV", "test")

	// 构建DSN - 从配置文件中读取
	dbConfig := helpers.DatabaseConfig{
		Driver:   "mysql",
		Host:     "10.68.73.213",
		Port:     3306,
		Username: "kysion_test",
		Password: "Qm2bhXr5dFWMJ5yt",
		Database: "kysiondb_test",
		Charset:  "utf8mb4",
	}
	dsn := formatDSN(dbConfig)

	// 连接测试数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	assert.NoError(t, err, "数据库连接失败")
	assert.NotNil(t, db, "数据库连接应该不为nil")

	// 健康检查
	sqlDB, err := db.DB()
	assert.NoError(t, err, "获取底层数据库连接失败")
	err = sqlDB.Ping()
	assert.NoError(t, err, "数据库Ping失败")
}

func TestDatabaseTransactions(t *testing.T) {
	// 不再跳过测试
	// t.Skip("需要真实数据库连接")

	// 获取项目根目录路径
	_, err := findRootDir()
	assert.NoError(t, err, "查找项目根目录失败")

	// 设置测试环境
	os.Setenv("KY_ADMIN_ENV", "test")

	// 构建DSN
	dbConfig := helpers.DatabaseConfig{
		Driver:   "mysql",
		Host:     "10.68.73.213",
		Port:     3306,
		Username: "kysion_test",
		Password: "Qm2bhXr5dFWMJ5yt",
		Database: "kysiondb_test",
		Charset:  "utf8mb4",
	}
	dsn := formatDSN(dbConfig)

	// 连接测试数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	assert.NoError(t, err, "数据库连接失败")

	// 测试事务
	err = helpers.WithTransaction(db, func(tx interface{}) error {
		// 在事务内执行一些简单操作，如查询版本
		var version string
		result := tx.(*gorm.DB).Raw("SELECT VERSION()").Scan(&version)
		assert.NoError(t, result.Error, "查询版本失败")
		assert.NotEmpty(t, version, "版本不应为空")
		return nil
	})
	assert.NoError(t, err, "事务执行失败")

	// 测试事务回滚
	err = helpers.WithTransaction(db, func(tx interface{}) error {
		// 执行一些操作，然后返回错误以触发回滚
		return assert.AnError
	})
	assert.Error(t, err, "应该返回错误以触发事务回滚")
}

func TestDatabaseConfig(t *testing.T) {
	// 创建测试配置
	dbConfig := helpers.DatabaseConfig{
		Driver:          "mysql",
		Host:            "localhost",
		Port:            3306,
		Username:        "testuser",
		Password:        "testpass",
		Database:        "testdb",
		Charset:         "utf8mb4",
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: 3600,
		LogLevel:        4,
	}

	// 测试配置有效性
	dsn := formatDSN(dbConfig)
	assert.Contains(t, dsn, "testuser:testpass@tcp(localhost:3306)/testdb", "DSN格式不正确")
	assert.Contains(t, dsn, "charset=utf8mb4", "DSN中的字符集不正确")
}

func formatDSN(config helpers.DatabaseConfig) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.Charset,
	)
}

func TestConfigLoading(t *testing.T) {
	// 保存原始环境变量
	originalEnv := os.Getenv("KY_ADMIN_ENV")
	defer os.Setenv("KY_ADMIN_ENV", originalEnv)

	// 设置测试环境
	os.Setenv("KY_ADMIN_ENV", "test")

	// 使用相对路径加载配置文件，跳过测试
	t.Skip("配置文件正在调整中")

	// 尝试查找项目根目录
	rootDir, err := helpers.FindRootDir()
	assert.NoError(t, err, "查找项目根目录失败")
	assert.NotEmpty(t, rootDir, "项目根目录不应为空")
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
