package casbin

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// Setup 初始化Casbin
func Setup(db *gorm.DB, modelPath string) (*casbin.Enforcer, error) {
	// 初始化适配器
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	// 加载模型
	m, err := model.NewModelFromFile(modelPath)
	if err != nil {
		return nil, err
	}

	// 创建enforcer
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}

	// 加载策略
	err = enforcer.LoadPolicy()
	if err != nil {
		return nil, err
	}

	return enforcer, nil
}

// InitRBACPolicy 初始化基本的RBAC策略
func InitRBACPolicy(enforcer *casbin.Enforcer) error {
	// 清除现有策略
	enforcer.ClearPolicy()

	// 添加角色
	// 超级管理员
	_, err := enforcer.AddPolicy("admin", "*", "*")
	if err != nil {
		return err
	}

	// 管理员角色 - 有用户管理权限
	_, err = enforcer.AddPolicy("manager", "/api/users/*", "*")
	if err != nil {
		return err
	}
	_, err = enforcer.AddPolicy("manager", "/api/roles", "GET")
	if err != nil {
		return err
	}

	// 普通用户角色 - 只有基本权限
	_, err = enforcer.AddPolicy("user", "/api/user", "GET")
	if err != nil {
		return err
	}
	_, err = enforcer.AddPolicy("user", "/api/user", "PUT")
	if err != nil {
		return err
	}

	// 新建初始用户的角色关系
	_, err = enforcer.AddGroupingPolicy("1", "admin") // 用户ID 1 -> admin角色
	if err != nil {
		return err
	}
	_, err = enforcer.AddGroupingPolicy("2", "user") // 用户ID 2 -> user角色
	if err != nil {
		return err
	}

	// 保存策略到数据库
	return enforcer.SavePolicy()
}
