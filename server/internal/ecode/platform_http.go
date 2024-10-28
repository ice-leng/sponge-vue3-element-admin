package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// platform business-level http error codes.
// the platformNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	platformNO       = 35
	platformName     = "platform"
	platformBaseCode = errcode.HCode(platformNO)

	ErrCreatePlatform     = errcode.NewError(platformBaseCode+1, "failed to create "+platformName)
	ErrDeleteByIDPlatform = errcode.NewError(platformBaseCode+2, "failed to delete "+platformName)
	ErrUpdateByIDPlatform = errcode.NewError(platformBaseCode+3, "failed to update "+platformName)
	ErrGetByIDPlatform    = errcode.NewError(platformBaseCode+4, "failed to get "+platformName+" details")
	ErrListPlatform       = errcode.NewError(platformBaseCode+5, "failed to list of "+platformName)

	// error codes are globally unique, adding 1 to the previous error code
)
