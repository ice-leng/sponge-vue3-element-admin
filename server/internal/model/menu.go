package model

import (
	"github.com/zhufuyi/sponge/pkg/ggorm"
)

// Menu 菜单管理
type Menu struct {
	ggorm.Model `gorm:"embedded"` // embed id and time

	ParentID   int    `gorm:"column:parent_id;type:int(11);default:0;NOT NULL" json:"parentID"`        // 父级
	Name       string `gorm:"column:name;type:varchar(32);NOT NULL" json:"name"`                       // 菜单名称
	Type       string `gorm:"column:type;type:enum('CATALOG','MENU','BUTTON','EXTLINK')" json:"type"`  // 菜单类型(CATALOG-菜单；MENU-目录；BUTTON-按钮；EXTLINK-外链)
	Path       string `gorm:"column:path;type:varchar(255);NOT NULL" json:"path"`                      // 路由路径
	Component  string `gorm:"column:component;type:varchar(255);NOT NULL" json:"component"`            // 组件路径(vue页面完整路径，省略.vue后缀)
	Perm       string `gorm:"column:perm;type:varchar(255)" json:"perm"`                               // 权限标识
	Sort       int    `gorm:"column:sort;type:int(11);default:1;NOT NULL" json:"sort"`                 // 排序
	Visible    int    `gorm:"column:visible;type:tinyint(4);NOT NULL" json:"visible"`                  // 显示状态
	Icon       string `gorm:"column:icon;type:varchar(255);NOT NULL" json:"icon"`                      // 菜单图标
	Redirect   string `gorm:"column:redirect;type:varchar(255);NOT NULL" json:"redirect"`              // 跳转路径
	AlwaysShow int    `gorm:"column:always_show;type:tinyint(4);default:0;NOT NULL" json:"alwaysShow"` // 始终显示
	KeepAlive  int    `gorm:"column:keep_alive;type:tinyint(4);default:1;NOT NULL" json:"keepAlive"`   // 始终显示
	Params     string `gorm:"column:params;type:json" json:"params"`                                   // 路由参数
}

// TableName table name
func (m *Menu) TableName() string {
	return "t_menu"
}

type MenuMeta struct {
	Title      string  `json:"title"`      // 标题
	Icon       string  `json:"icon"`       // 图标
	Hidden     bool    `json:"hidden"`     // 隐藏
	AlwaysShow bool    `json:"alwaysShow"` // 始终显示
	Params     *string `json:"params"`     // 路由参数
}

type MenuItem struct {
	Path      string     `json:"path"`      // 路由路径
	Name      string     `json:"name"`      // 路由名称
	Component string     `json:"component"` // 组件路径
	Redirect  string     `json:"redirect"`  // 跳转路径
	Meta      MenuMeta   `json:"meta"`      // meta
	Children  []Children `json:"children"`  // 子级
}

type ChildrenMeta struct {
	Title      string  `json:"title"`      // 标题
	Icon       string  `json:"icon"`       // 图标
	Hidden     bool    `json:"hidden"`     // 隐藏
	AlwaysShow bool    `json:"alwaysShow"` // 始终显示
	Params     *string `json:"params"`     // 路由参数
	KeepAlive  bool    `json:"keepAlive"`  // 始终显示
}

type Children struct {
	Path      string       `json:"path"`      // 路由路径
	Name      string       `json:"name"`      // 路由名称
	Component string       `json:"component"` // 组件路径
	Meta      ChildrenMeta `json:"meta"`      // meta
}
