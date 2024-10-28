package model

import (
	"github.com/zhufuyi/sponge/pkg/ggorm"
)

// Role 角色管理
type Role struct {
	ggorm.Model `gorm:"embedded"` // embed id and time

	Name   string `gorm:"column:name;type:varchar(32);NOT NULL" json:"name"`       // 角色名称
	Code   string `gorm:"column:code;type:varchar(32);NOT NULL" json:"code"`       // 角色编码
	Sort   int    `gorm:"column:sort;type:int(11);default:1;NOT NULL" json:"sort"` // 排序
	Status int    `gorm:"column:status;type:tinyint(4);NOT NULL" json:"status"`    // 状态
}

// TableName table name
func (m *Role) TableName() string {
	return "t_role"
}
