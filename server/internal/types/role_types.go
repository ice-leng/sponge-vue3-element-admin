package types

import (
	"time"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateRoleRequest request params
type CreateRoleRequest struct {
	Name   string `json:"name" binding:""`   // 角色名称
	Code   string `json:"code" binding:""`   // 角色编码
	Sort   int    `json:"sort" binding:""`   // 排序
	Status int    `json:"status" binding:""` // 状态
}

// UpdateRoleByIDRequest request params
type UpdateRoleByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	Name   string `json:"name" binding:""`   // 角色名称
	Code   string `json:"code" binding:""`   // 角色编码
	Sort   int    `json:"sort" binding:""`   // 排序
	Status int    `json:"status" binding:""` // 状态
}

// RoleObjDetail detail
type RoleObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	CreatedAt time.Time `json:"createdAt"` // 创建时间
	UpdatedAt time.Time `json:"updatedAt"` // 更新时间
	Name      string    `json:"name"`      // 角色名称
	Code      string    `json:"code"`      // 角色编码
	Sort      int       `json:"sort"`      // 排序
	Status    int       `json:"status"`    // 状态
}

// CreateRoleReply only for api docs
type CreateRoleReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteRoleByIDReply only for api docs
type DeleteRoleByIDReply struct {
	Result
}

// UpdateRoleByIDReply only for api docs
type UpdateRoleByIDReply struct {
	Result
}

// GetRoleByIDReply only for api docs
type GetRoleByIDReply struct {
	Code int           `json:"code"` // return code
	Msg  string        `json:"msg"`  // return information description
	Data RoleObjDetail `json:"data"` // return data
}

// ListRolesRequest request params
type ListRolesRequest struct {
	Page     int    `json:"page,omitempty" form:"page" binding:""`         // 分页
	PageSize int    `json:"pageSize,omitempty" form:"pageSize" binding:""` // 分页大小
	Sort     string `json:"sort,omitempty" form:"sort" binding:""`         // 排序
	Status   *int   `json:"status" form:"status" binding:""`               // 状态

	StartTime string `json:"startTime,omitempty" form:"startTime" binding:""` // 开始时间
	EndTime   string `json:"endTime,omitempty" form:"endTime" binding:""`     // 结束时间
}

// ListRolesReply only for api docs
type ListRolesReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		List  []RoleObjDetail `json:"roles"`
		Total int             `json:"total"`
	} `json:"data"` // return data
}
