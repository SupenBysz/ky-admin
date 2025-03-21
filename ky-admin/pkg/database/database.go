package database

import (
	"fmt"
	"time"

	"github.com/SupenBysz/ky-admin/pkg/config"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// DBManager 数据库管理器
type DBManager struct {
	db        *gorm.DB
	config    *config.Config
	logger    *zap.Logger
	connected bool
}

// NewDBManager 创建数据库管理器
func NewDBManager(config *config.Config) *DBManager {
	return &DBManager{
		config: config,
	}
}

// SetLogger 设置日志记录器
func (m *DBManager) SetLogger(logger *zap.Logger) {
	m.logger = logger
}

// Connect 连接数据库
func (m *DBManager) Connect() error {
	if m.connected && m.db != nil {
		return nil
	}

	// 确保日志记录器已设置
	if m.logger == nil {
		m.logger = zap.L()
	}

	// 构建DSN
	dsn := m.buildDSN()

	// 配置GORM
	gormConfig := m.buildGormConfig()

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		m.logger.Error("数据库连接失败", zap.Error(err))
		return err
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		m.logger.Error("获取底层sqlDB失败", zap.Error(err))
		return err
	}

	// 设置最大连接数
	maxConnections := m.config.GetIntWithDefault("database.max_connections", 100)
	sqlDB.SetMaxOpenConns(maxConnections)

	// 设置最大空闲连接数
	maxIdleConnections := m.config.GetIntWithDefault("database.max_idle_connections", 10)
	sqlDB.SetMaxIdleConns(maxIdleConnections)

	// 设置连接最大存活时间
	connMaxLifetime := m.config.GetDurationWithDefault("database.conn_max_lifetime", 1*time.Hour)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	// 设置空闲连接最大存活时间
	connMaxIdleTime := m.config.GetDurationWithDefault("database.conn_max_idle_time", 30*time.Minute)
	sqlDB.SetConnMaxIdleTime(connMaxIdleTime)

	// 验证连接
	if err := sqlDB.Ping(); err != nil {
		m.logger.Error("数据库Ping失败", zap.Error(err))
		return err
	}

	m.db = db
	m.connected = true
	m.logger.Info("数据库连接成功")

	return nil
}

// Close 关闭数据库连接
func (m *DBManager) Close() error {
	if !m.connected || m.db == nil {
		return nil
	}

	sqlDB, err := m.db.DB()
	if err != nil {
		m.logger.Error("获取底层sqlDB失败", zap.Error(err))
		return err
	}

	if err := sqlDB.Close(); err != nil {
		m.logger.Error("关闭数据库连接失败", zap.Error(err))
		return err
	}

	m.connected = false
	m.db = nil
	m.logger.Info("数据库连接已关闭")

	return nil
}

// GetDB 获取GORM DB实例
func (m *DBManager) GetDB() *gorm.DB {
	return m.db
}

// CheckHealth 检查数据库健康状态
func (m *DBManager) CheckHealth() error {
	if !m.connected || m.db == nil {
		return fmt.Errorf("数据库未连接")
	}

	sqlDB, err := m.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}

// buildDSN 构建数据库连接字符串
func (m *DBManager) buildDSN() string {
	host := m.config.GetStringWithDefault("database.host", "localhost")
	port := m.config.GetIntWithDefault("database.port", 3306)
	username := m.config.GetStringWithDefault("database.username", "root")
	password := m.config.GetStringWithDefault("database.password", "")
	dbname := m.config.GetStringWithDefault("database.name", "ky_admin")
	charset := m.config.GetStringWithDefault("database.charset", "utf8mb4")
	parseTime := m.config.GetBoolWithDefault("database.parse_time", true)
	loc := m.config.GetStringWithDefault("database.loc", "Local")

	// 构建DSN
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		username, password, host, port, dbname, charset,
		parseTime, loc,
	)

	// 添加可选参数
	timeout := m.config.GetIntWithDefault("database.timeout", 10)
	if timeout > 0 {
		dsn = fmt.Sprintf("%s&timeout=%ds", dsn, timeout)
	}

	return dsn
}

// buildGormConfig 构建GORM配置
func (m *DBManager) buildGormConfig() *gorm.Config {
	// 配置GORM
	gormConfig := &gorm.Config{
		// 禁用默认事务
		SkipDefaultTransaction: m.config.GetBoolWithDefault("database.skip_default_transaction", true),
		// 命名策略
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   m.config.GetStringWithDefault("database.table_prefix", ""),
			SingularTable: m.config.GetBoolWithDefault("database.singular_table", false),
		},
		// 禁用外键约束
		DisableForeignKeyConstraintWhenMigrating: m.config.GetBoolWithDefault("database.disable_foreign_key_constraint", true),
	}

	// 配置日志级别
	logLevel := logger.Info
	switch m.config.GetStringWithDefault("database.log_level", "info") {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	}

	// 配置GORM日志
	gormConfig.Logger = logger.New(
		// 配置日志输出
		NewGormLogWriter(m.logger),
		logger.Config{
			SlowThreshold:             m.config.GetDurationWithDefault("database.slow_threshold", 200*time.Millisecond),
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: m.config.GetBoolWithDefault("database.ignore_record_not_found", true),
			Colorful:                  false,
		},
	)

	return gormConfig
}

// NewGormLogWriter 创建GORM日志写入器
type GormLogWriter struct {
	logger *zap.Logger
}

// NewGormLogWriter 创建GORM日志写入器
func NewGormLogWriter(logger *zap.Logger) *GormLogWriter {
	return &GormLogWriter{logger: logger}
}

// Printf 实现GORM日志写入器接口
func (w *GormLogWriter) Printf(format string, args ...interface{}) {
	w.logger.Debug(fmt.Sprintf(format, args...))
}
