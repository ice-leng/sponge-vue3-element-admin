package middlewares

import (
	"admin/internal/cache"
	"admin/internal/dao"
	"admin/internal/ecode"
	"admin/internal/model"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/zhufuyi/sponge/pkg/jwt"
	"github.com/zhufuyi/sponge/pkg/utils"
	"time"
)

var iPlatformDao dao.PlatformDao

func VerifyToken(claims *jwt.Claims, tokenTail10 string, c *gin.Context) error {
	if claims.ExpiresAt.Before(time.Now()) {
		return ecode.ErrLogin.Err()
	}

	if iPlatformDao == nil {
		iPlatformDao = dao.NewPlatformDao(
			model.GetDB(),
			cache.NewPlatformCache(model.GetCacheType()),
		)
	}

	platform, err := iPlatformDao.GetByID(context.Background(), utils.StrToUint64(claims.UID))

	if err != nil {
		return err
	}

	if *platform.Status != 1 {
		return ecode.ErrLoginFrozen.Err()
	}

	c.Set("roleId", platform.RoleID)
	c.Set("id", platform.ID)
	return nil
}
