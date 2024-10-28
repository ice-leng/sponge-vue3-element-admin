package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// config business-level http error codes.
// the configNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	configNO       = 38
	configName     = "config"
	configBaseCode = errcode.HCode(configNO)

	ErrCreateConfig     = errcode.NewError(configBaseCode+1, "failed to create "+configName)
	ErrDeleteByIDConfig = errcode.NewError(configBaseCode+2, "failed to delete "+configName)
	ErrUpdateByIDConfig = errcode.NewError(configBaseCode+3, "failed to update "+configName)
	ErrGetByIDConfig    = errcode.NewError(configBaseCode+4, "failed to get "+configName+" details")
	ErrListConfig       = errcode.NewError(configBaseCode+5, "failed to list of "+configName)

	// error codes are globally unique, adding 1 to the previous error code
)
