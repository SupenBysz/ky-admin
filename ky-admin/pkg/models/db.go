package models

import (
	"errors"
	"sync"

	"gorm.io/gorm"
)

var (
	// DB 全局数据库实例
	DB *gorm.DB
	// 确保DB只初始化一次的锁
	dbOnce sync.Once
	// 全局数据库错误
	dbErr error
	// ErrDBNotRegistered 数据库未注册错误
	ErrDBNotRegistered = errors.New("数据库尚未注册")
)

// RegisterDB 注册全局数据库实例
// 这个函数应该在应用启动时由main函数调用，确保在使用任何模型前完成初始化
func RegisterDB(db *gorm.DB) {
	dbOnce.Do(func() {
		if db == nil {
			dbErr = errors.New("无效的数据库连接")
			return
		}
		DB = db
	})
}

// ResetDB 重置全局数据库实例
// 主要用于测试场景
func ResetDB() {
	DB = nil
	// 重置dbOnce以允许再次调用RegisterDB
	dbOnce = sync.Once{}
}

// GetDB 获取全局数据库实例
func GetDB() (*gorm.DB, error) {
	if DB == nil {
		return nil, ErrDBNotRegistered
	}
	return DB, nil
}

// GormPagination 通用分页结构
type GormPagination struct {
	Page      int   `json:"page" form:"page"`           // 页码
	PageSize  int   `json:"page_size" form:"page_size"` // 每页数量
	Total     int64 `json:"total"`                      // 总数
	TotalPage int   `json:"total_page"`                 // 总页数
}

// DefaultPageSize 默认每页数量
const DefaultPageSize = 10

// SetPagination 设置分页参数
func (p *GormPagination) SetPagination() *gorm.DB {
	// 设置默认值
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = DefaultPageSize
	}

	// 计算偏移量
	offset := (p.Page - 1) * p.PageSize

	// 返回带有分页的查询构建器
	return DB.Offset(offset).Limit(p.PageSize)
}

// CalculateTotalPage 计算总页数
func (p *GormPagination) CalculateTotalPage() {
	if p.PageSize > 0 {
		p.TotalPage = int((p.Total + int64(p.PageSize) - 1) / int64(p.PageSize))
	} else {
		p.TotalPage = 0
	}
}
