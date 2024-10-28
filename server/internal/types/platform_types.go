package types

import (
	"time"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreatePlatformRequest request params
type CreatePlatformRequest struct {
	Username       string    `json:"username" binding:""`       // 账号
	Password       string    `json:"password" binding:""`       // 密码
	Avatar         string    `json:"avatar" binding:""`         // 头像
	RoleID         string    `json:"roleID" binding:""`         // 角色
	Status         int       `json:"status" binding:""`         // 状态
	LastTime       time.Time `json:"lastTime" binding:""`       // 上次登录时间
	ClaimTimeLimit int       `json:"claimTimeLimit" binding:""` // 领取时间限制（小时）
}

// UpdatePlatformByIDRequest request params
type UpdatePlatformByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	Username       string    `json:"username" binding:""`       // 账号
	Password       string    `json:"password" binding:""`       // 密码
	Avatar         string    `json:"avatar" binding:""`         // 头像
	RoleID         string    `json:"roleID" binding:""`         // 角色
	Status         int       `json:"status" binding:""`         // 状态
	LastTime       time.Time `json:"lastTime" binding:""`       // 上次登录时间
	ClaimTimeLimit int       `json:"claimTimeLimit" binding:""` // 领取时间限制（小时）
}

// PlatformObjDetail detail
type PlatformObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	CreatedAt      time.Time `json:"createdAt"`      // 创建时间
	UpdatedAt      time.Time `json:"updatedAt"`      // 更新时间
	Username       string    `json:"username"`       // 账号
	Password       string    `json:"password"`       // 密码
	Avatar         string    `json:"avatar"`         // 头像
	RoleID         string    `json:"roleID"`         // 角色
	Status         int       `json:"status"`         // 状态
	LastTime       time.Time `json:"lastTime"`       // 上次登录时间
	ClaimTimeLimit int       `json:"claimTimeLimit"` // 领取时间限制（小时）
}

// CreatePlatformReply only for api docs
type CreatePlatformReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeletePlatformByIDReply only for api docs
type DeletePlatformByIDReply struct {
	Result
}

// UpdatePlatformByIDReply only for api docs
type UpdatePlatformByIDReply struct {
	Result
}

// GetPlatformByIDReply only for api docs
type GetPlatformByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Platform PlatformObjDetail `json:"platform"`
	} `json:"data"` // return data
}

// ListPlatformsRequest request params
type ListPlatformsRequest struct {
	Page     int    `json:"page,omitempty" form:"page" binding:""`         // 分页
	PageSize int    `json:"pageSize,omitempty" form:"pageSize" binding:""` // 分页大小
	Sort     string `json:"sort,omitempty" form:"sort" binding:""`         // 排序

	StartTime string `json:"startTime,omitempty" form:"startTime" binding:""` // 开始时间
	EndTime   string `json:"endTime,omitempty" form:"endTime" binding:""`     // 结束时间
}

// ListPlatformsReply only for api docs
type ListPlatformsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Platforms []PlatformObjDetail `json:"platforms"`
	} `json:"data"` // return data
}
