package types

import (
	"time"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreatePlatformRequest request params
type CreatePlatformRequest struct {
	Username string   `json:"username" binding:""` // 账号
	Password string   `json:"password" binding:""` // 密码
	Avatar   string   `json:"avatar" binding:""`   // 头像
	RoleID   []uint64 `json:"roleId" binding:""`   // 角色
	Status   int      `json:"status" binding:""`   // 状态
}

// UpdatePlatformByIDRequest request params
type UpdatePlatformByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	Username string   `json:"username" binding:""` // 账号
	Password string   `json:"password" binding:""` // 密码
	Avatar   string   `json:"avatar" binding:""`   // 头像
	RoleID   []uint64 `json:"roleId" binding:""`   // 角色
	Status   *int     `json:"status" binding:""`   // 状态
}

type LoginRequest struct {
	Username    string `json:"username" binding:""`    // 账号
	Password    string `json:"password" binding:""`    // 密码
	CaptchaKey  string `json:"captchaKey" binding:""`  // 验证码key
	CaptchaCode string `json:"captchaCode" binding:""` // 验证码code
}

type ChangePasswordRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	OldPassword     string `json:"oldPassword" binding:"required"`                         // 密码
	NewPassword     string `json:"newPassword" binding:"required"`                         // 密码
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=NewPassword"` // 密码
}

type ResetPasswordRequest struct {
	ID       uint64 `json:"id" binding:""`               // uint64 id
	Password string `json:"password" binding:"required"` // 密码
}

// PlatformObjDetail detail
type PlatformObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	CreatedAt LocalDateTime `json:"createdAt"` // 创建时间
	UpdatedAt LocalDateTime `json:"updatedAt"` // 更新时间
	Username  string        `json:"username"`  // 账号
	Avatar    string        `json:"avatar"`    // 头像
	RoleID    []uint64      `json:"roleId"`    // 角色
	Status    int           `json:"status"`    // 状态
	LastTime  LocalDateTime `json:"lastTime"`  // 上次登录时间
}

// PlatformListPage list
type PlatformListPage struct {
	PlatformObjDetail
	RoleNames []string `json:"roleNames"`
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
	Code int               `json:"code"` // return code
	Msg  string            `json:"msg"`  // return information description
	Data PlatformObjDetail `json:"data"` // return data
}

// ListPlatformsRequest request params
type ListPlatformsRequest struct {
	Page     int    `json:"page,omitempty" form:"page" binding:""`         // 分页
	PageSize int    `json:"pageSize,omitempty" form:"pageSize" binding:""` // 分页大小
	Sort     string `json:"sort,omitempty" form:"sort" binding:""`         // 排序

	StartTime string `json:"startTime,omitempty" form:"startTime" binding:""` // 开始时间
	EndTime   string `json:"endTime,omitempty" form:"endTime" binding:""`     // 结束时间
	Username  string `json:"username,omitempty" form:"username" binding:""`   // 账号
	Status    *int   `json:"status,omitempty" form:"status" binding:""`       // 状态
}

// ListPlatformsReply only for api docs
type ListPlatformsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		List  []PlatformObjDetail `json:"list"`
		Total int                 `json:"total"`
	} `json:"data"` // return data
}

type LoginReply struct {
	Code int       `json:"code"` // return code
	Msg  string    `json:"msg"`  // return information description
	Data LoginItem `json:"data"` // return data
}

type LoginItem struct {
	AccessToken string `json:"accessToken"` // access token
	Expires     int    `json:"expires"`     // expire time
	TokenType   string `json:"tokenType"`   // token type
}

type CaptchaReply struct {
	Code int         `json:"code"` // return code
	Msg  string      `json:"msg"`  // return information description
	Data CaptchaItem `json:"data"` // return data
}

type CaptchaItem struct {
	CaptchaKey    string `json:"captchaKey"`
	CaptchaBase64 string `json:"captchaBase64"`
}

type MeReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data MeItem `json:"data"` // return data
}

type MeItem struct {
	ID       uint64   `json:"id"`       // convert to uint64 id
	Username string   `json:"username"` // 账号
	Avatar   string   `json:"avatar"`   // 头像
	Roles    []string `json:"roles"`    // 角色组
	Perms    []string `json:"perms"`    // 权限组
}

type ProfileReply struct {
	Code int         `json:"code"` // return code
	Msg  string      `json:"msg"`  // return information description
	Data ProfileItem `json:"data"` // return data
}

type ProfileItem struct {
	ID        uint64        `json:"id"`        // convert to uint64 id
	Username  string        `json:"username"`  // 账号
	Avatar    string        `json:"avatar"`    // 头像
	Roles     string        `json:"roleNames"` // 角色组
	CreatedAt LocalDateTime `json:"createdAt"` // 创建时间
}
