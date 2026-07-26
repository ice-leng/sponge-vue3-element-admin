package handlerfunc

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/errcode"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/utils"
)

// GetCurrentUid get the current user id from context
func GetCurrentUid(c *gin.Context, key ...string) uint64 {
	if len(key) == 0 {
		key = append(key, "id")
	}
	uid, ok := c.Get(key[0])
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
		logger.Warn("StrToUint64E error", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

// GetIdsFromPath get multiple ids from path parameter (comma separated)
func GetIdsFromPath(c *gin.Context) (string, []uint64, bool) {
	idStr := c.Param("id")
	if idStr == "" {
		return "", nil, true
	}
	ids, isAbort := GetArrUint64FromString(idStr)
	return idStr, ids, isAbort
}

// GetArrUint64FromString string change []uint64
func GetArrUint64FromString(str string) ([]uint64, bool) {
	parts := strings.Split(str, ",")
	ids := make([]uint64, 0, len(parts))
	for _, v := range parts {
		id, err := utils.StrToUint64E(strings.TrimSpace(v))
		if err != nil || id == 0 {
			logger.Warn("StrToUint64E error", logger.Err(err), logger.String("str", v))
			return nil, true
		}
		ids = append(ids, id)
	}
	return ids, false
}

// GetArrStringFromString string change []string
func GetArrStringFromString(str string) []string {
	parts := strings.Split(str, ",")
	result := make([]string, 0, len(parts))
	for _, v := range parts {
		result = append(result, strings.TrimSpace(v))
	}
	return result
}
