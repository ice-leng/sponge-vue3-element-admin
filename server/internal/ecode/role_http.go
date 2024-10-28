package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// role business-level http error codes.
// the roleNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	roleNO       = 84
	roleName     = "role"
	roleBaseCode = errcode.HCode(roleNO)

	ErrCreateRole     = errcode.NewError(roleBaseCode+1, "failed to create "+roleName)
	ErrDeleteByIDRole = errcode.NewError(roleBaseCode+2, "failed to delete "+roleName)
	ErrUpdateByIDRole = errcode.NewError(roleBaseCode+3, "failed to update "+roleName)
	ErrGetByIDRole    = errcode.NewError(roleBaseCode+4, "failed to get "+roleName+" details")
	ErrListRole       = errcode.NewError(roleBaseCode+5, "failed to list of "+roleName)

	// error codes are globally unique, adding 1 to the previous error code
)
