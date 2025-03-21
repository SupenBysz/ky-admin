// Package dto 包含数据传输对象的定义
package dto

import "time"

// PermissionCreateDTO 权限创建请求数据传输对象
type PermissionCreateDTO struct {
	Name        string `json:"name" binding:"required" example:"用户管理"`        // 权限名称
	Code        string `json:"code" binding:"required" example:"user:manage"` // 权限编码
	Description string `json:"description" example:"用户管理相关权限"`                // 权限描述
	Type        int    `json:"type" example:"1"`                              // 权限类型 1:菜单 2:按钮 3:接口
	ParentID    uint   `json:"parent_id" example:"0"`                         // 父级ID
	Path        string `json:"path" example:"/users"`                         // 路径
	Method      string `json:"method" example:"GET"`                          // HTTP方法
	Status      int    `json:"status" example:"1"`                            // 状态 (使用int保持一致)
	Sort        int    `json:"sort" example:"1"`                              // 排序
	Icon        string `json:"icon" example:"user"`                           // 图标
}

// PermissionUpdateDTO 权限更新请求数据传输对象
type PermissionUpdateDTO struct {
	Name        string `json:"name" example:"用户管理"`            // 权限名称
	Description string `json:"description" example:"用户管理相关权限"` // 权限描述
	ParentID    *uint  `json:"parent_id" example:"0"`          // 父级ID
	Path        string `json:"path" example:"/users"`          // 路径
	Method      string `json:"method" example:"GET"`           // HTTP方法
	Status      *int   `json:"status" example:"1"`             // 状态 (使用int保持一致)
	Sort        *int   `json:"sort" example:"1"`               // 排序
	Icon        string `json:"icon" example:"user"`            // 图标
}

// PermissionResponseDTO 权限响应数据传输对象
type PermissionResponseDTO struct {
	ID          uint                    `json:"id"`                 // 权限ID
	Name        string                  `json:"name"`               // 权限名称
	Code        string                  `json:"code"`               // 权限编码
	Description string                  `json:"description"`        // 权限描述
	Type        int                     `json:"type"`               // 权限类型 1:菜单 2:按钮 3:接口
	ParentID    uint                    `json:"parent_id"`          // 父级ID
	Path        string                  `json:"path"`               // 路径
	Method      string                  `json:"method"`             // HTTP方法
	Status      int                     `json:"status"`             // 状态 (使用int保持一致)
	Sort        int                     `json:"sort"`               // 排序
	Icon        string                  `json:"icon"`               // 图标
	CreatedAt   time.Time               `json:"created_at"`         // 创建时间
	UpdatedAt   time.Time               `json:"updated_at"`         // 更新时间
	Children    []PermissionResponseDTO `json:"children,omitempty"` // 子权限
}

// PermissionTreeResponseDTO 权限树响应数据传输对象
type PermissionTreeResponseDTO struct {
	List []PermissionResponseDTO `json:"list"` // 权限列表
}

// PermissionPageQueryDTO 权限分页查询参数
type PermissionPageQueryDTO struct {
	Name   string `form:"name" json:"name"`           // 权限名称
	Code   string `form:"code" json:"code"`           // 权限编码
	Type   *int   `form:"type" json:"type"`           // 权限类型
	Status *int   `form:"status" json:"status"`       // 状态 (使用int保持一致)
	Page   int    `form:"page" json:"page"`           // 页码
	Size   int    `form:"page_size" json:"page_size"` // 每页条数
}

// PermissionPageResponseDTO 权限分页查询响应
type PermissionPageResponseDTO struct {
	Total int64                   `json:"total"` // 总记录数
	List  []PermissionResponseDTO `json:"list"`  // 权限列表
	Page  int                     `json:"page"`  // 当前页码
	Size  int                     `json:"size"`  // 每页大小
}
