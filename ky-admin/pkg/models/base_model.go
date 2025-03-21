package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型，包含通用字段
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`                 // 主键ID
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`           // 创建时间
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`           // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`    // 删除时间（软删除）
	CreatedBy uint           `gorm:"not null;default:0" json:"created_by"` // 创建者ID
	UpdatedBy uint           `gorm:"not null;default:0" json:"updated_by"` // 更新者ID
}

// BeforeCreate GORM钩子，创建记录前自动设置时间
func (m *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	if m.UpdatedAt.IsZero() {
		m.UpdatedAt = time.Now()
	}

	// 从上下文中获取当前用户ID
	if userID, ok := tx.Statement.Context.Value("current_user_id").(uint); ok && userID > 0 {
		m.CreatedBy = userID
		m.UpdatedBy = userID
	}

	return nil
}

// BeforeUpdate GORM钩子，更新记录前自动更新时间
func (m *BaseModel) BeforeUpdate(tx *gorm.DB) error {
	m.UpdatedAt = time.Now()

	// 从上下文中获取当前用户ID
	if userID, ok := tx.Statement.Context.Value("current_user_id").(uint); ok && userID > 0 {
		m.UpdatedBy = userID
	}

	return nil
}

// BaseModelWithoutDelete 不包含软删除字段的基础模型
type BaseModelWithoutDelete struct {
	ID        uint      `gorm:"primarykey" json:"id"`                 // 主键ID
	CreatedAt time.Time `gorm:"not null" json:"created_at"`           // 创建时间
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`           // 更新时间
	CreatedBy uint      `gorm:"not null;default:0" json:"created_by"` // 创建者ID
	UpdatedBy uint      `gorm:"not null;default:0" json:"updated_by"` // 更新者ID
}

// BeforeCreate GORM钩子，创建记录前自动设置时间
func (m *BaseModelWithoutDelete) BeforeCreate(tx *gorm.DB) error {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	if m.UpdatedAt.IsZero() {
		m.UpdatedAt = time.Now()
	}

	// 从上下文中获取当前用户ID
	if userID, ok := tx.Statement.Context.Value("current_user_id").(uint); ok && userID > 0 {
		m.CreatedBy = userID
		m.UpdatedBy = userID
	}

	return nil
}

// BeforeUpdate GORM钩子，更新记录前自动更新时间
func (m *BaseModelWithoutDelete) BeforeUpdate(tx *gorm.DB) error {
	m.UpdatedAt = time.Now()

	// 从上下文中获取当前用户ID
	if userID, ok := tx.Statement.Context.Value("current_user_id").(uint); ok && userID > 0 {
		m.UpdatedBy = userID
	}

	return nil
}

// StatusModel 具有状态字段的基础模型
type StatusModel struct {
	BaseModel
	Status uint8 `gorm:"not null;default:1;comment:状态 1:正常 0:禁用" json:"status"` // 状态
}

// EnableDisable 启用/禁用接口
type EnableDisable interface {
	Enable() error  // 启用
	Disable() error // 禁用
}

// Pagination 分页参数
type Pagination struct {
	Page     int `json:"page" form:"page"`           // 页码
	PageSize int `json:"page_size" form:"page_size"` // 每页数量
}

// GetOffset 获取偏移量
func (p *Pagination) GetOffset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.GetLimit()
}

// GetLimit 获取限制数量
func (p *Pagination) GetLimit() int {
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	return p.PageSize
}

// PageResult 分页结果
type PageResult struct {
	List     interface{} `json:"list"`      // 数据列表
	Total    int64       `json:"total"`     // 总数
	Page     int         `json:"page"`      // 当前页码
	PageSize int         `json:"page_size"` // 每页数量
	Pages    int         `json:"pages"`     // 总页数
}

// NewPageResult 创建分页结果
func NewPageResult(list interface{}, total int64, page, pageSize int) *PageResult {
	pages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		pages++
	}
	return &PageResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Pages:    pages,
	}
}
