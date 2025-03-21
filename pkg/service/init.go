package service

import (
	"sync"

	"github.com/SupenBysz/ky-admin/pkg/database"
)

var (
	dbService     *DBService
	dbServiceOnce sync.Once

	authService     *AuthService
	authServiceOnce sync.Once
)

// GetDBService 获取数据库服务实例（单例模式）
func GetDBService() *DBService {
	dbServiceOnce.Do(func() {
		dbService = NewDBService(database.GetDB())
	})
	return dbService
}

// GetAuthService 获取认证服务实例（单例模式）
func GetAuthService() *AuthService {
	authServiceOnce.Do(func() {
		authService = NewAuthService(GetDBService())
	})
	return authService
}

// InitServices 初始化所有服务
func InitServices() {
	// 初始化数据库服务
	GetDBService()

	// 初始化认证服务
	GetAuthService()

	// 未来可以在这里初始化其他服务
}
