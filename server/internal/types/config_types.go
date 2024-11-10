package types

import (
	"time"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateConfigRequest request params
type CreateConfigRequest struct {
	Name        string `json:"name" binding:""`        // 配置名称
	Description string `json:"description" binding:""` // 描述
	Key         string `json:"key" binding:""`         // 配置键
	Value       string `json:"value" binding:""`       // 配置值
}

// UpdateConfigByIDRequest request params
type UpdateConfigByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	Name        string `json:"name" binding:""`        // 配置名称
	Description string `json:"description" binding:""` // 描述
	Key         string `json:"key" binding:""`         // 配置键
	Value       string `json:"value" binding:""`       // 配置值
}

// ConfigObjDetail detail
type ConfigObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	CreatedAt   time.Time `json:"createdAt"`   // 创建时间
	UpdatedAt   time.Time `json:"updatedAt"`   // 更新时间
	Name        string    `json:"name"`        // 配置名称
	Description string    `json:"description"` // 描述
	Key         string    `json:"key"`         // 配置键
	Value       string    `json:"value"`       // 配置值
}

// CreateConfigReply only for api docs
type CreateConfigReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteConfigByIDReply only for api docs
type DeleteConfigByIDReply struct {
	Result
}

// UpdateConfigByIDReply only for api docs
type UpdateConfigByIDReply struct {
	Result
}

// GetConfigByIDReply only for api docs
type GetConfigByIDReply struct {
	Code int             `json:"code"` // return code
	Msg  string          `json:"msg"`  // return information description
	Data ConfigObjDetail `json:"data"` // return data
}

// ListConfigsRequest request params
type ListConfigsRequest struct {
	Page     int    `json:"page,omitempty" form:"page" binding:""`         // 分页
	PageSize int    `json:"pageSize,omitempty" form:"pageSize" binding:""` // 分页大小
	Sort     string `json:"sort,omitempty" form:"sort" binding:""`         // 排序

	StartTime string `json:"startTime,omitempty" form:"startTime" binding:""` // 开始时间
	EndTime   string `json:"endTime,omitempty" form:"endTime" binding:""`     // 结束时间
}

// ListConfigsReply only for api docs
type ListConfigsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Configs []ConfigObjDetail `json:"configs"`
	} `json:"data"` // return data
}
