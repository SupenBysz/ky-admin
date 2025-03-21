package models

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// MigrateManager 数据库迁移管理器
type MigrateManager struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewMigrateManager 创建迁移管理器
func NewMigrateManager(db *gorm.DB, logger *zap.Logger) *MigrateManager {
	return &MigrateManager{
		db:     db,
		logger: logger,
	}
}

// MigrateModels 执行数据库迁移（自动创建表）
func (m *MigrateManager) MigrateModels() error {
	if m.db == nil {
		return fmt.Errorf("数据库未初始化")
	}

	m.logger.Info("开始执行数据库迁移")
	startTime := time.Now()

	// 定义需要迁移的模型
	models := []interface{}{
		&User{},
		&Role{},
		&Permission{},
	}

	// 执行迁移
	err := m.db.AutoMigrate(models...)
	if err != nil {
		m.logger.Error("数据库迁移失败", zap.Error(err))
		return err
	}

	m.logger.Info("数据库迁移完成", zap.Duration("耗时", time.Since(startTime)))
	return nil
}

// CreateInitData 创建初始数据
func (m *MigrateManager) CreateInitData() error {
	if m.db == nil {
		return fmt.Errorf("数据库未初始化")
	}

	m.logger.Info("开始初始化基础数据")

	// 使用事务确保数据一致性
	err := m.db.Transaction(func(tx *gorm.DB) error {
		// 创建管理员角色
		adminRole := &Role{
			StatusModel: StatusModel{
				Status: uint8(StatusEnabled),
			},
			Name:        "超级管理员",
			Code:        "admin",
			Description: "系统超级管理员，拥有所有权限",
			Sort:        0,
			IsSystem:    true,
		}

		// 检查角色是否已存在
		var existingRole Role
		result := tx.Where("code = ?", adminRole.Code).First(&existingRole)
		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			return result.Error
		}

		// 不存在则创建
		if result.Error == gorm.ErrRecordNotFound {
			if err := tx.Create(adminRole).Error; err != nil {
				m.logger.Error("创建管理员角色失败", zap.Error(err))
				return err
			}
			m.logger.Info("成功创建管理员角色", zap.String("role", adminRole.Name))
		} else {
			adminRole = &existingRole
			m.logger.Info("管理员角色已存在", zap.String("role", adminRole.Name))
		}

		// 创建管理员用户
		adminUser := &User{
			StatusModel: StatusModel{
				Status: uint8(StatusEnabled),
			},
			Username: "admin",
			Password: "admin123", // 将在BeforeSave中加密
			Nickname: "系统管理员",
			Email:    "admin@example.com",
			Remark:   "初始管理员账号",
		}

		// 检查用户是否已存在
		var existingUser User
		result = tx.Where("username = ?", adminUser.Username).First(&existingUser)
		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			return result.Error
		}

		// 不存在则创建
		if result.Error == gorm.ErrRecordNotFound {
			if err := tx.Create(adminUser).Error; err != nil {
				m.logger.Error("创建管理员用户失败", zap.Error(err))
				return err
			}
			m.logger.Info("成功创建管理员用户", zap.String("user", adminUser.Username))

			// 为管理员用户分配管理员角色
			if err := tx.Model(adminUser).Association("Roles").Append(adminRole); err != nil {
				m.logger.Error("为管理员用户分配角色失败", zap.Error(err))
				return err
			}
		} else {
			m.logger.Info("管理员用户已存在", zap.String("user", adminUser.Username))
		}

		return nil
	})

	if err != nil {
		m.logger.Error("初始化基础数据失败", zap.Error(err))
		return err
	}

	m.logger.Info("初始化基础数据完成")
	return nil
}

// BackupDatabaseInfo 数据库备份信息
type BackupDatabaseInfo struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

// BackupDatabase 备份数据库
func (m *MigrateManager) BackupDatabase(backupDir string, dbInfo BackupDatabaseInfo) (string, error) {
	// 确保备份目录存在
	if backupDir == "" {
		backupDir = "./backups"
	}

	if err := os.MkdirAll(backupDir, 0755); err != nil {
		m.logger.Error("创建备份目录失败", zap.Error(err))
		return "", err
	}

	// 生成备份文件名
	timestamp := time.Now().Format("20060102_150405")
	backupFile := filepath.Join(backupDir, fmt.Sprintf("%s_%s.sql", dbInfo.Database, timestamp))

	// 此处仅记录备份信息，实际备份操作将由调用者执行
	m.logger.Info("准备备份数据库",
		zap.String("database", dbInfo.Database),
		zap.String("file", backupFile),
	)

	return backupFile, nil
}

// MigrationRecord 迁移记录
type MigrationRecord struct {
	ID          uint      `gorm:"primarykey"`
	Version     string    `gorm:"size:50;not null;uniqueIndex"`
	Description string    `gorm:"size:200;not null"`
	ExecutedAt  time.Time `gorm:"not null"`
}

// TableName 指定表名
func (MigrationRecord) TableName() string {
	return "migrations"
}

// SeedData 种子数据
func SeedData(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// 检查是否已有数据
		var count int64
		if err := tx.Model(&User{}).Count(&count).Error; err != nil {
			return err
		}

		// 已有数据，跳过初始化
		if count > 0 {
			return nil
		}

		// 创建超级管理员角色
		adminRole := &Role{
			Name:        "超级管理员",
			Code:        "admin",
			Description: "系统超级管理员",
			StatusModel: StatusModel{
				Status: uint8(StatusEnabled),
			},
		}

		if err := tx.Create(adminRole).Error; err != nil {
			return err
		}

		// 创建基础权限
		permissions := []Permission{
			{
				Name:        "系统管理",
				Code:        "system:manage",
				Type:        PermissionTypeSystem,
				Description: "系统管理权限",
				StatusModel: StatusModel{
					Status: uint8(StatusEnabled),
				},
			},
			{
				Name:        "用户管理",
				Code:        "user:manage",
				Type:        PermissionTypeMenu,
				Description: "用户管理权限",
				StatusModel: StatusModel{
					Status: uint8(StatusEnabled),
				},
			},
			{
				Name:        "角色管理",
				Code:        "role:manage",
				Type:        PermissionTypeMenu,
				Description: "角色管理权限",
				StatusModel: StatusModel{
					Status: uint8(StatusEnabled),
				},
			},
			{
				Name:        "权限管理",
				Code:        "permission:manage",
				Type:        PermissionTypeMenu,
				Description: "权限管理权限",
				StatusModel: StatusModel{
					Status: uint8(StatusEnabled),
				},
			},
		}

		if err := tx.Create(&permissions).Error; err != nil {
			return err
		}

		// 为超级管理员角色分配所有权限
		var allPermissionIDs []uint
		for _, p := range permissions {
			allPermissionIDs = append(allPermissionIDs, p.ID)
		}

		if err := adminRole.AssignPermissions(tx, allPermissionIDs); err != nil {
			return err
		}

		// 创建初始管理员用户
		adminUser := &User{
			Username: "admin",
			Email:    "admin@kysion.com",
			Nickname: "系统管理员",
			StatusModel: StatusModel{
				Status: uint8(StatusEnabled),
			},
			IsAdmin: true,
		}

		// 设置默认密码 admin123
		if err := adminUser.SetPassword("admin123"); err != nil {
			return err
		}

		if err := tx.Create(adminUser).Error; err != nil {
			return err
		}

		// 为用户分配角色
		if err := adminUser.AssignRole(tx, adminRole.ID); err != nil {
			return err
		}

		return nil
	})
}
