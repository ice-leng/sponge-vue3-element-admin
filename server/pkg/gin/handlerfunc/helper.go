package handlerfunc

import (
	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/errcode"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/utils"
)

// GetCurrentUid get the current user id from context
func GetCurrentUid(c *gin.Context) uint64 {
	uid, ok := c.Get("id")
	if !ok {
		return 0
	}
	return uid.(uint64)
}

// IsErrcode check if the error is errcode.Error
func IsErrcode(err error) (*errcode.Error, bool) {
	ec := errcode.ParseError(err)
	if ec.Code() > 0 {
		return ec, true
	}
	return nil, false
}

// GetIdFromPath get the id from path parameter
func GetIdFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}
