package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// menu business-level http error codes.
// the menuNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	menuNO       = 54
	menuName     = "menu"
	menuBaseCode = errcode.HCode(menuNO)

	ErrCreateMenu     = errcode.NewError(menuBaseCode+1, "failed to create "+menuName)
	ErrDeleteByIDMenu = errcode.NewError(menuBaseCode+2, "failed to delete "+menuName)
	ErrUpdateByIDMenu = errcode.NewError(menuBaseCode+3, "failed to update "+menuName)
	ErrGetByIDMenu    = errcode.NewError(menuBaseCode+4, "failed to get "+menuName+" details")
	ErrListMenu       = errcode.NewError(menuBaseCode+5, "failed to list of "+menuName)

	// error codes are globally unique, adding 1 to the previous error code
)
