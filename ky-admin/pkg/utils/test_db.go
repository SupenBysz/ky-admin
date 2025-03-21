package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/SupenBysz/ky-admin/pkg/config"
)

// SetupTestDB 设置测试数据库
func SetupTestDB() (*gorm.DB, error) {
	// 创建内存数据库
	gormConfig := &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Info,
				Colorful:      true,
			},
		),
	}

	// 使用SQLite内存数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("连接测试数据库失败: %w", err)
	}

	return db, nil
}

// CreateReportDir 创建测试报告目录
func CreateReportDir() (string, error) {
	testConfig, err := config.GetTestConfig()
	if err != nil {
		return "", fmt.Errorf("获取测试配置失败: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	reportDir := filepath.Join(testConfig.Reports.Dir, timestamp)

	if err := os.MkdirAll(reportDir, 0755); err != nil {
		return "", fmt.Errorf("创建测试报告目录失败: %w", err)
	}

	// 创建latest软链接
	latestLink := filepath.Join(testConfig.Reports.Dir, "latest")
	_ = os.Remove(latestLink) // 尝试删除可能存在的链接
	if err := os.Symlink(reportDir, latestLink); err != nil {
		log.Printf("创建最新报告链接失败: %v", err)
	}

	return reportDir, nil
}

// autoMigrate 自动迁移数据库结构
func autoMigrate(db *gorm.DB) error {
	// TODO: 添加需要迁移的模型
	// return db.AutoMigrate(&model.User{}, &model.Role{}, &model.Permission{})
	return nil
}

// seedTestData 加载测试数据
func seedTestData(db *gorm.DB) error {
	// 添加测试数据
	// TODO: 添加测试数据

	return nil
}

// CleanupTestDB 清理测试数据库
func CleanupTestDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %w", err)
	}

	return sqlDB.Close()
}
