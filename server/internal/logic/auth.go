package logic

import (
	"admin/internal/cache"
	"admin/internal/dao"
	"admin/internal/database"
	"admin/internal/ecode"
	"admin/internal/model"
	"admin/internal/types"
	"context"
	"fmt"
	"image/color"
	"time"

	"github.com/go-dev-frame/sponge/pkg/gin/middleware/auth"
	"github.com/go-dev-frame/sponge/pkg/gocrypto"
	"github.com/go-dev-frame/sponge/pkg/sgorm"
	"github.com/go-dev-frame/sponge/pkg/utils"
	"github.com/mojocn/base64Captcha"
	"github.com/redis/go-redis/v9"
)

type AuthLogic interface {
	Login(ctx context.Context, request *types.LoginRequest) (*types.LoginItem, error)
	Captcha(ctx context.Context) (*types.CaptchaItem, error)
	Logout(ctx context.Context) error
}

type authLogic struct {
	iDao    dao.PlatformDao
	captcha *base64Captcha.DriverMath
	redis   *redis.Client
}

func NewAuthLogic() AuthLogic {
	return NewAuthLogicByDAO(
		dao.NewPlatformDao(
			database.GetDB(),
			cache.NewPlatformCache(database.GetCacheType()),
		),
		database.GetRedisCli(),
	)
}

func NewAuthLogicByDAO(iDao dao.PlatformDao, redisCli *redis.Client) AuthLogic {
	bgColor := color.RGBA{R: 0, G: 0, B: 0, A: 0}
	driver := base64Captcha.NewDriverMath(60, 240, 0, 0, &bgColor, nil, []string{
		"wqy-microhei.ttc",
	})
	return &authLogic{
		iDao:    iDao,
		captcha: driver,
		redis:   redisCli,
	}
}

func (a authLogic) Login(ctx context.Context, request *types.LoginRequest) (*types.LoginItem, error) {
	code, err := a.redis.Get(ctx, fmt.Sprintf("captcha:%s", request.CaptchaKey)).Result()
	if err != nil || code != request.CaptchaCode {
		return nil, ecode.ErrLoginCaptcha.Err()
	}

	platform, platformErr := a.iDao.GetByUsername(ctx, request.Username)
	if platformErr != nil {
		return nil, ecode.ErrLogin.Err()
	}

	ok := gocrypto.VerifyPassword(request.Password, platform.Password)
	if !ok {
		return nil, ecode.ErrLogin.Err()
	}

	lastTime := time.Now()
	_ = a.iDao.UpdateByID(ctx, &model.Platform{
		Model: sgorm.Model{
			ID: platform.ID,
		},
		LastTime: &lastTime,
	})

	token, tokenErr := auth.GenerateToken(utils.Uint64ToStr(platform.ID))
	if tokenErr != nil {
		return nil, ecode.ErrLogin.Err()
	}

	return &types.LoginItem{
		AccessToken: token,
		Expires:     7200,
		TokenType:   "Bearer",
	}, nil
}

func (a authLogic) Captcha(ctx context.Context) (*types.CaptchaItem, error) {
	id, content, answer := a.captcha.GenerateIdQuestionAnswer()
	item, _ := a.captcha.DrawCaptcha(content)
	result := &types.CaptchaItem{
		CaptchaKey:    id,
		CaptchaBase64: item.EncodeB64string(),
	}
	a.redis.Set(ctx, fmt.Sprintf("captcha:%s", id), answer, 2*time.Minute)
	return result, nil
}

func (a authLogic) Logout(ctx context.Context) error {
	return nil
}
