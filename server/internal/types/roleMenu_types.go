package types

import (
	"time"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateRoleMenuRequest request params
type CreateRoleMenuRequest struct {
	RoleID int `json:"roleID" binding:""` // 角色ID
	MenuID int `json:"menuID" binding:""` // 菜单ID
}

// UpdateRoleMenuByIDRequest request params
type UpdateRoleMenuByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	RoleID int `json:"roleID" binding:""` // 角色ID
	MenuID int `json:"menuID" binding:""` // 菜单ID
}

// RoleMenuObjDetail detail
type RoleMenuObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	CreatedAt time.Time `json:"createdAt"` // 创建时间
	UpdatedAt time.Time `json:"updatedAt"` // 更新时间
	RoleID    int       `json:"roleID"`    // 角色ID
	MenuID    int       `json:"menuID"`    // 菜单ID
}

// CreateRoleMenuReply only for api docs
type CreateRoleMenuReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteRoleMenuByIDReply only for api docs
type DeleteRoleMenuByIDReply struct {
	Result
}

// UpdateRoleMenuByIDReply only for api docs
type UpdateRoleMenuByIDReply struct {
	Result
}

// GetRoleMenuByIDReply only for api docs
type GetRoleMenuByIDReply struct {
	Code int               `json:"code"` // return code
	Msg  string            `json:"msg"`  // return information description
	Data RoleMenuObjDetail `json:"data"` // return data
}

// ListRoleMenusRequest request params
type ListRoleMenusRequest struct {
	Page     int    `json:"page,omitempty" form:"page" binding:""`         // 分页
	PageSize int    `json:"pageSize,omitempty" form:"pageSize" binding:""` // 分页大小
	Sort     string `json:"sort,omitempty" form:"sort" binding:""`         // 排序

	RoleId    *uint64 `json:"roleId,omitempty" form:"roleId" binding:""`       // 排序
	StartTime string  `json:"startTime,omitempty" form:"startTime" binding:""` // 开始时间
	EndTime   string  `json:"endTime,omitempty" form:"endTime" binding:""`     // 结束时间
}

// ListRoleMenusReply only for api docs
type ListRoleMenusReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		RoleMenus []RoleMenuObjDetail `json:"roleMenus"`
	} `json:"data"` // return data
}
