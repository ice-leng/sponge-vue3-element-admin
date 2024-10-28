package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// roleMenu business-level http error codes.
// the roleMenuNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	roleMenuNO       = 50
	roleMenuName     = "roleMenu"
	roleMenuBaseCode = errcode.HCode(roleMenuNO)

	ErrCreateRoleMenu     = errcode.NewError(roleMenuBaseCode+1, "failed to create "+roleMenuName)
	ErrDeleteByIDRoleMenu = errcode.NewError(roleMenuBaseCode+2, "failed to delete "+roleMenuName)
	ErrUpdateByIDRoleMenu = errcode.NewError(roleMenuBaseCode+3, "failed to update "+roleMenuName)
	ErrGetByIDRoleMenu    = errcode.NewError(roleMenuBaseCode+4, "failed to get "+roleMenuName+" details")
	ErrListRoleMenu       = errcode.NewError(roleMenuBaseCode+5, "failed to list of "+roleMenuName)

	// error codes are globally unique, adding 1 to the previous error code
)
