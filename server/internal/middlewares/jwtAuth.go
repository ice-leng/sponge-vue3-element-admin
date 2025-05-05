package middlewares

import (
	"admin/internal/cache"
	"admin/internal/dao"
	"admin/internal/database"
	"admin/internal/ecode"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/jwt"
	"github.com/go-dev-frame/sponge/pkg/utils"
	"time"
)

var iPlatformDao dao.PlatformDao

const JwtSignKey = "UxeY8GUv4CH8fH7hCQM9CA2"

func VerifyToken(claims *jwt.Claims, c *gin.Context) error {
	if claims.ExpiresAt.Before(time.Now().Add(time.Minute * 10)) {
		token, err := claims.NewToken(time.Hour*2, jwt.HS384, []byte(JwtSignKey))
		if err != nil {
			return err
		}
		c.Header("X-Renewed-Token", token)
	}

	if iPlatformDao == nil {
		iPlatformDao = dao.NewPlatformDao(
			database.GetDB(),
			cache.NewPlatformCache(database.GetCacheType()),
		)
	}

	platform, err := iPlatformDao.GetByID(context.Background(), utils.StrToUint64(claims.UID))

	if err != nil {
		return err
	}

	if *platform.Status != 1 {
		return ecode.ErrLoginFrozen.Err()
	}

	c.Set("id", platform.ID)
	c.Set("roleId", platform.RoleID)
	return nil
}
