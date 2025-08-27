package handler

import (
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

type baseHandler struct{}

func (b *baseHandler) isRoleByCode(c *gin.Context, code string) bool {
	roleCode, ok := c.Get("roleCode")
	if !ok {
		return false
	}
	return slices.Contains(roleCode.([]string), code)
}

func (b *baseHandler) getUid(c *gin.Context) uint64 {
	uid, ok := c.Get("id")
	if !ok {
		return 0
	}
	return uid.(uint64)
}

func (b *baseHandler) isErrcode(err error) (*errcode.Error, bool) {
	ec := errcode.ParseError(err)
	if ec.Code() > 0 {
		return ec, true
	}
	return nil, false
}
