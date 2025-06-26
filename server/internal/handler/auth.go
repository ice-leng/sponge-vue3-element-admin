package handler

import (
	"admin/internal/cache"
	"admin/internal/dao"
	"admin/internal/database"
	"admin/internal/ecode"
	"admin/internal/model"
	"admin/internal/types"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware/auth"
	"github.com/go-dev-frame/sponge/pkg/gin/response"
	"github.com/go-dev-frame/sponge/pkg/gocrypto"
	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/sgorm"
	"github.com/go-dev-frame/sponge/pkg/utils"
	"github.com/mojocn/base64Captcha"
	"github.com/redis/go-redis/v9"
	"image/color"
	"time"
)

type AuthHandler interface {
	Login(c *gin.Context)
	Captcha(c *gin.Context)
	Logout(c *gin.Context)
}

type authHandler struct {
	iDao    dao.PlatformDao
	captcha *base64Captcha.DriverMath
	redis   *redis.Client
}

func NewAuthHandler() AuthHandler {
	bgColor := color.RGBA{R: 0, G: 0, B: 0, A: 0}
	driver := base64Captcha.NewDriverMath(60, 240, 0, 0, &bgColor, nil, []string{
		"wqy-microhei.ttc",
	})
	return &authHandler{
		iDao: dao.NewPlatformDao(
			database.GetDB(),
			cache.NewPlatformCache(database.GetCacheType()),
		),
		captcha: driver,
		redis:   database.GetRedisCli(),
	}
}

// Login
// @Summary with username and password
// @Description with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param data body types.LoginRequest true "login information"
// @Success 200 {object} types.LoginReply{}
// @Router /api/v1/auth/login [post]
func (a authHandler) Login(c *gin.Context) {
	request := &types.LoginRequest{}
	err := c.ShouldBindJSON(request)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	code := a.redis.Get(context.Background(), fmt.Sprintf("captcha:%s", request.CaptchaKey)).Val()
	if code != request.CaptchaCode {
		response.Error(c, ecode.ErrLoginCaptcha)
		return
	}

	platform, platformErr := a.iDao.GetByUsername(c, request.Username)
	if platformErr != nil {
		response.Error(c, ecode.ErrLogin)
		return
	}

	ok := gocrypto.VerifyPassword(request.Password, platform.Password)
	if !ok {
		response.Error(c, ecode.ErrLogin)
		return
	}

	lastTime := time.Now()
	_ = a.iDao.UpdateByID(c, &model.Platform{
		Model: sgorm.Model{
			ID: platform.ID,
		},
		LastTime: &lastTime,
	})

	token, tokenErr := auth.GenerateToken(utils.Uint64ToStr(platform.ID))
	if tokenErr != nil {
		response.Error(c, ecode.ErrLogin)
		return
	}

	response.Success(c, a.loginReply(token))
}

func (a authHandler) loginReply(token string) types.LoginItem {
	return types.LoginItem{
		AccessToken: token,
		Expires:     7200,
		TokenType:   "Bearer",
	}
}

// Logout of logout
// @Summary logout
// @Description logout
// @Tags auth
// @accept json
// @Produce json
// @Success 200 {object} types.Result{}
// @Router /api/v1/auth/logout [delete]
// @Security BearerAuth
func (a authHandler) Logout(c *gin.Context) {
	response.Success(c)
}

// Captcha get a captcha
// @Summary get a captcha
// @Description get a captcha
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} types.CaptchaReply{}
// @Router /api/v1/auth/captcha [get]
func (a authHandler) Captcha(c *gin.Context) {
	id, content, answer := a.captcha.GenerateIdQuestionAnswer()
	item, _ := a.captcha.DrawCaptcha(content)
	result := types.CaptchaItem{
		CaptchaKey:    id,
		CaptchaBase64: item.EncodeB64string(),
	}
	a.redis.Set(context.Background(), fmt.Sprintf("captcha:%s", id), answer, 2*time.Minute)
	response.Success(c, result)
}
