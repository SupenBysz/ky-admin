package database

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	// 全局数据库管理器实例
	globalDBManager *DBManager
	// 确保只初始化一次的锁
	globalDBManagerOnce sync.Once
)

// RegisterGlobalDBManager 注册全局数据库管理器实例
func RegisterGlobalDBManager(manager *DBManager) {
	if manager == nil {
		return
	}

	globalDBManagerOnce.Do(func() {
		globalDBManager = manager
	})
}

// GetGlobalDBManager 获取全局数据库管理器实例
func GetGlobalDBManager() (*DBManager, error) {
	if globalDBManager == nil {
		return nil, fmt.Errorf("全局数据库管理器未初始化")
	}
	return globalDBManager, nil
}

// GetDB 获取全局数据库连接（兼容旧代码）
func GetDB() *gorm.DB {
	if globalDBManager == nil {
		return nil
	}
	return globalDBManager.GetDB()
}

// GlobalRunInTransaction 在全局数据库连接上执行事务（兼容旧代码）
func GlobalRunInTransaction(fn TxFunc) error {
	db := GetDB()
	if db == nil {
		return fmt.Errorf("全局数据库管理器未初始化")
	}
	return RunInTransactionWithDB(db, fn)
}

// GlobalRunInTransactionWithContext 在全局数据库连接上执行带上下文的事务（兼容旧代码）
func GlobalRunInTransactionWithContext(ctx interface{}, fn TxFunc) error {
	db := GetDB()
	if db == nil {
		return fmt.Errorf("全局数据库管理器未初始化")
	}

	if ctxObj, ok := ctx.(context.Context); ok {
		return RunInTransactionWithContextAndDB(ctxObj, db, fn)
	}

	return RunInTransactionWithDB(db, fn)
}

// GetMigrationManager 创建迁移管理器（兼容旧代码）
func GetMigrationManager() (*DBMigrationManager, error) {
	manager, err := GetGlobalDBManager()
	if err != nil {
		return nil, err
	}

	db := manager.GetDB()
	if db == nil {
		return nil, fmt.Errorf("数据库连接未初始化")
	}

	return NewDBMigrationManager(db), nil
}

// DBMigrationManager 数据库迁移管理器
type DBMigrationManager struct {
	db *gorm.DB
}

// NewDBMigrationManager 创建新的迁移管理器
func NewDBMigrationManager(db *gorm.DB) *DBMigrationManager {
	return &DBMigrationManager{
		db: db,
	}
}

// RunMigration 执行自动迁移
func (m *DBMigrationManager) RunMigration(models ...interface{}) error {
	if m.db == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return m.db.AutoMigrate(models...)
}

// RunSeed 执行种子数据初始化
func (m *DBMigrationManager) RunSeed(seedFn func(tx *gorm.DB) error) error {
	if m.db == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return m.db.Transaction(seedFn)
}

// BackupDBInfo 备份数据库信息
type BackupDBInfo struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

// GetBackupDBInfo 获取数据库备份信息
func (m *DBManager) GetBackupDBInfo() BackupDBInfo {
	return BackupDBInfo{
		Host:     m.config.GetStringWithDefault("database.host", "localhost"),
		Port:     m.config.GetIntWithDefault("database.port", 3306),
		Username: m.config.GetStringWithDefault("database.username", "root"),
		Password: m.config.GetStringWithDefault("database.password", ""),
		Database: m.config.GetStringWithDefault("database.name", "ky_admin"),
	}
}

// WithLogger 日志接口兼容层
func WithLogger(logger *zap.Logger) interface{} {
	return logger
}

// CheckTableExists 检查表是否存在
func (m *DBMigrationManager) CheckTableExists(tableName string) (bool, error) {
	if m.db == nil {
		return false, fmt.Errorf("数据库未初始化")
	}

	var count int64
	err := m.db.Raw("SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?", tableName).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetTableColumns 获取表的列名
func (m *DBMigrationManager) GetTableColumns(tableName string) ([]string, error) {
	if m.db == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	var columns []string
	err := m.db.Raw("SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?", tableName).Scan(&columns).Error
	if err != nil {
		return nil, err
	}

	return columns, nil
}

// ExecuteSQL 执行SQL语句
func (m *DBMigrationManager) ExecuteSQL(sql string, args ...interface{}) error {
	if m.db == nil {
		return fmt.Errorf("数据库未初始化")
	}

	return m.db.Exec(sql, args...).Error
}
