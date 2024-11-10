package model

import (
	"github.com/zhufuyi/sponge/pkg/ggorm"
)

// RoleMenu 角色菜单关联
type RoleMenu struct {
	ggorm.Model `gorm:"embedded"` // embed id and time

	RoleID uint64 `gorm:"column:role_id;type:int(11);default:0;NOT NULL" json:"roleID"` // 角色ID
	MenuID uint64 `gorm:"column:menu_id;type:int(11);default:0;NOT NULL" json:"menuID"` // 菜单ID
}

// TableName table name
func (m *RoleMenu) TableName() string {
	return "t_role_menu"
}
