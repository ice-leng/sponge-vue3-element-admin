package types

import (
	"time"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateMenuRequest request params
type CreateMenuRequest struct {
	ParentID   int    `json:"parentId" binding:""`   // 父级
	Name       string `json:"name" binding:""`       // 菜单名称
	Type       string `json:"type" binding:""`       // 菜单类型(CATALOG-菜单；MENU-目录；BUTTON-按钮；EXTLINK-外链)
	Path       string `json:"routePath" binding:""`  // 路由路径
	Component  string `json:"component" binding:""`  // 组件路径(vue页面完整路径，省略.vue后缀)
	Perm       string `json:"perm" binding:""`       // 权限标识
	Sort       int    `json:"sort" binding:""`       // 排序
	Visible    int    `json:"visible" binding:""`    // 显示状态
	Icon       string `json:"icon" binding:""`       // 菜单图标
	Redirect   string `json:"redirect" binding:""`   // 跳转路径
	AlwaysShow int    `json:"alwaysShow" binding:""` // 始终显示
	KeepAlive  int    `json:"keepAlive" binding:""`  // 始终显示
	Params     string `json:"params" binding:""`     // 路由参数
}

// UpdateMenuByIDRequest request params
type UpdateMenuByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	ParentID   int    `json:"parentId" binding:""`   // 父级
	Name       string `json:"name" binding:""`       // 菜单名称
	Type       string `json:"type" binding:""`       // 菜单类型(CATALOG-菜单；MENU-目录；BUTTON-按钮；EXTLINK-外链)
	Path       string `json:"routePath" binding:""`  // 路由路径
	Component  string `json:"component" binding:""`  // 组件路径(vue页面完整路径，省略.vue后缀)
	Perm       string `json:"perm" binding:""`       // 权限标识
	Sort       int    `json:"sort" binding:""`       // 排序
	Visible    int    `json:"visible" binding:""`    // 显示状态
	Icon       string `json:"icon" binding:""`       // 菜单图标
	Redirect   string `json:"redirect" binding:""`   // 跳转路径
	AlwaysShow int    `json:"alwaysShow" binding:""` // 始终显示
	KeepAlive  int    `json:"keepAlive" binding:""`  // 始终显示
	Params     string `json:"params" binding:""`     // 路由参数
}

// MenuObjDetail detail
type MenuObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	CreatedAt  time.Time `json:"createdAt"`  // 创建时间
	UpdatedAt  time.Time `json:"updatedAt"`  // 更新时间
	ParentID   int       `json:"parentId"`   // 父级
	Name       string    `json:"name"`       //
	RouteName  string    `json:"routeName"`  // 路由名称
	Type       string    `json:"type"`       // 菜单类型(CATALOG-菜单；MENU-目录；BUTTON-按钮；EXTLINK-外链)
	Path       string    `json:"routePath"`  // 路由路径
	Component  string    `json:"component"`  // 组件路径(vue页面完整路径，省略.vue后缀)
	Perm       string    `json:"perm"`       // 权限标识
	Sort       int       `json:"sort"`       // 排序
	Visible    int       `json:"visible"`    // 显示状态
	Icon       string    `json:"icon"`       // 菜单图标
	Redirect   string    `json:"redirect"`   // 跳转路径
	AlwaysShow int       `json:"alwaysShow"` // 始终显示
	KeepAlive  int       `json:"keepAlive"`  // 始终显示
	Params     string    `json:"params"`     // 路由参数
}

// MenuObjPage page
type MenuObjPage struct {
	ID uint64 `json:"id"` // convert to uint64 id

	ParentID  int            `json:"parentId"`  // 父级
	Name      string         `json:"name"`      // 菜单名称
	Type      string         `json:"type"`      // 菜单类型(CATALOG-菜单；MENU-目录；BUTTON-按钮；EXTLINK-外链)
	Path      string         `json:"routePath"` // 路由路径
	Component string         `json:"component"` // 组件路径(vue页面完整路径，省略.vue后缀)
	Perm      string         `json:"perm"`      // 权限标识
	Sort      int            `json:"sort"`      // 排序
	Visible   int            `json:"visible"`   // 显示状态
	Icon      string         `json:"icon"`      // 菜单图标
	Redirect  string         `json:"redirect"`  // 跳转路径
	RouteName string         `json:"routeName"` // 路由名称
	Children  []*MenuObjPage `json:"children"`
}

// CreateMenuReply only for api docs
type CreateMenuReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteMenuByIDReply only for api docs
type DeleteMenuByIDReply struct {
	Result
}

// UpdateMenuByIDReply only for api docs
type UpdateMenuByIDReply struct {
	Result
}

// GetMenuByIDReply only for api docs
type GetMenuByIDReply struct {
	Code int           `json:"code"` // return code
	Msg  string        `json:"msg"`  // return information description
	Data MenuObjDetail `json:"data"` // return data
}

// ListMenusRequest request params
type ListMenusRequest struct {
	Page     int    `json:"page,omitempty" form:"page" binding:""`         // 分页
	PageSize int    `json:"pageSize,omitempty" form:"pageSize" binding:""` // 分页大小
	Sort     string `json:"sort,omitempty" form:"sort" binding:""`         // 排序

	ParentID *uint64 `json:"parentId,omitempty" form:"parentId" binding:""`
	Keywords string  `json:"keywords,omitempty" form:"keywords" binding:""`

	StartTime string `json:"startTime,omitempty" form:"startTime" binding:""` // 开始时间
	EndTime   string `json:"endTime,omitempty" form:"endTime" binding:""`     // 结束时间
}

// OptionMenusRequest request params
type OptionMenusRequest struct {
	OnlyParent bool `json:"onlyParent,omitempty" form:"onlyParent" binding:""`
}

// ListMenusReply only for api docs
type ListMenusReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		List  []MenuObjPage `json:"list"`
		Total int           `json:"total"`
	} `json:"data"` // return data
}

type RouteReply struct {
}
