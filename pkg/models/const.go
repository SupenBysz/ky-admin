package models

// 状态常量定义
const (
	StatusDisabled uint = iota // 禁用
	StatusEnabled              // 启用
	StatusDeleted              // 已删除
)

// 权限类型常量定义
const (
	PermissionTypeSystem uint = iota + 1 // 系统权限
	PermissionTypeMenu                   // 菜单权限
	PermissionTypeAction                 // 操作权限
	PermissionTypeData                   // 数据权限
)

// 性别常量定义
const (
	GenderUnknown uint = iota // 未知
	GenderMale                // 男
	GenderFemale              // 女
)

// 获取状态名称
func GetStatusName(status uint) string {
	switch status {
	case StatusDisabled:
		return "禁用"
	case StatusEnabled:
		return "启用"
	case StatusDeleted:
		return "已删除"
	default:
		return "未知"
	}
}

// 获取权限类型名称
func GetPermissionTypeName(permType uint) string {
	switch permType {
	case PermissionTypeSystem:
		return "系统权限"
	case PermissionTypeMenu:
		return "菜单权限"
	case PermissionTypeAction:
		return "操作权限"
	case PermissionTypeData:
		return "数据权限"
	default:
		return "未知"
	}
}

// 获取性别名称
func GetGenderName(gender uint) string {
	switch gender {
	case GenderMale:
		return "男"
	case GenderFemale:
		return "女"
	default:
		return "未知"
	}
}
