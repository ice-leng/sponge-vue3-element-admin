package logic

import (
	"admin/internal/cache"
	"admin/internal/constant"
	"admin/internal/dao"
	"admin/internal/database"
	"admin/internal/ecode"
	"admin/internal/model"
	"admin/internal/types"
	"context"
	"encoding/base64"
	"errors"
	"strings"

	"github.com/go-dev-frame/sponge/pkg/copier"
	"github.com/go-dev-frame/sponge/pkg/gocrypto"
	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

type PlatformLogic interface {
	Create(ctx context.Context, request *types.CreatePlatformRequest) (uint64, error)
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, request *types.UpdatePlatformByIDRequest) error
	GetByID(ctx context.Context, id uint64) (*types.PlatformObjDetail, error)
	List(ctx context.Context, request *types.ListPlatformsRequest) ([]*types.PlatformObjDetail, int64, error)
	Me(ctx context.Context, id uint64) (*types.MeItem, error)
	Profile(ctx context.Context, id uint64) (*types.ProfileItem, error)
	ChangePassword(ctx context.Context, request *types.ChangePasswordRequest) error
	ResetPassword(ctx context.Context, request *types.ResetPasswordRequest) error
}

type platformLogic struct {
	iDao       dao.PlatformDao
	iRoleDao   dao.RoleDao
	iConfigDao dao.ConfigDao
}

func NewPlatformLogic() PlatformLogic {
	return NewPlatformLogicByDAO(
		dao.NewPlatformDao(
			database.GetDB(), // db driver is mysql
			cache.NewPlatformCache(database.GetCacheType()),
		),
		dao.NewRoleDao(
			database.GetDB(),
			cache.NewRoleCache(database.GetCacheType()),
		),
		dao.NewConfigDao(
			database.GetDB(),
			cache.NewConfigCache(database.GetCacheType()),
		),
	)
}

func NewPlatformLogicByDAO(iDao dao.PlatformDao, iRoleDao dao.RoleDao, iConfigDao dao.ConfigDao) PlatformLogic {
	return &platformLogic{
		iDao:       iDao,
		iRoleDao:   iRoleDao,
		iConfigDao: iConfigDao,
	}
}

func (p platformLogic) Create(ctx context.Context, request *types.CreatePlatformRequest) (uint64, error) {
	table := &model.Platform{}
	err := copier.Copy(table, request)
	if err != nil {
		return 0, ecode.ErrCreatePlatform.Err()
	}
	table.Mobile = encryptMobile(table.Mobile)
	table.Password = convertPassword(table.Password)
	// Note: if copier.Copy cannot assign a value to a field, add it here
	err = p.iDao.Create(ctx, table)
	return table.ID, err
}

func (p platformLogic) DeleteByID(ctx context.Context, id uint64) error {
	return p.iDao.DeleteByID(ctx, id)
}

func (p platformLogic) UpdateByID(ctx context.Context, request *types.UpdatePlatformByIDRequest) error {
	table := &model.Platform{}
	err := copier.Copy(table, request)
	if err != nil {
		return ecode.ErrUpdateByIDPlatform.Err()
	}
	table.Mobile = encryptMobile(table.Mobile)
	table.Password = convertPassword(table.Password)

	// Note: if copier.Copy cannot assign a value to a field, add it here
	return p.iDao.UpdateByID(ctx, table)
}

func (p platformLogic) GetByID(ctx context.Context, id uint64) (*types.PlatformObjDetail, error) {
	result, err := p.iDao.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			return nil, ecode.ErrGetByIDPlatform.Err()
		}
		return nil, err
	}
	result.Avatar = p.iConfigDao.MakePathByConfig(ctx, result.Avatar, constant.ConfigKeyImageDomain)

	roleCodes := make(map[uint64]string)
	roles, _ := p.iRoleDao.GetByIDs(ctx, result.RoleID)
	if len(roles) > 0 {
		for _, role := range roles {
			roleCodes[role.ID] = role.Name
		}
	}
	data, err := convertPlatform(result, roleCodes)
	if err != nil {
		return nil, ecode.ErrGetByIDPlatform.Err()
	}

	return data, nil
}

func (p platformLogic) List(ctx context.Context, request *types.ListPlatformsRequest) ([]*types.PlatformObjDetail, int64, error) {
	params := &query.Params{
		Page:    request.Page - 1,
		Limit:   request.PageSize,
		Sort:    request.Sort,
		Columns: []query.Column{},
	}
	if request.StartTime != "" && request.EndTime != "" {
		params.Columns = append(params.Columns, query.Column{
			Name:  "created_at",
			Exp:   ">=",
			Value: request.StartTime,
		})
		params.Columns = append(params.Columns, query.Column{
			Name:  "created_at",
			Exp:   "<",
			Value: request.EndTime + " 23:59:59",
		})
	}
	if request.Mobile != "" {
		params.Columns = append(params.Columns, query.Column{
			Name:  "mobile",
			Exp:   "=",
			Value: request.Mobile,
		})
	}
	if request.Status != nil {
		params.Columns = append(params.Columns, query.Column{
			Name:  "status",
			Exp:   "=",
			Value: *request.Status,
		})
	}
	if request.Keyword != "" {
		params.Columns = append(params.Columns, query.Column{
			Name:  "keyword",
			Exp:   "like",
			Value: "%" + request.Keyword + "%",
		})
	}

	result, total, err := p.iDao.GetByColumns(ctx, params)
	if err != nil {
		return nil, 0, err
	}
	data, err := p.convertPlatforms(ctx, result)
	if err != nil {
		return nil, 0, ecode.ErrListPlatform.Err()
	}
	return data, total, nil
}

func (p platformLogic) Me(ctx context.Context, id uint64) (*types.MeItem, error) {
	reply := &types.MeItem{}
	platform, err := p.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	_ = copier.Copy(&reply, platform)

	perms, _ := p.iRoleDao.GetPermissionsByIds(ctx, platform.RoleID)
	reply.Perms = perms
	return reply, nil
}

func (p platformLogic) Profile(ctx context.Context, id uint64) (*types.ProfileItem, error) {
	reply := &types.ProfileItem{}
	platform, err := p.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	_ = copier.Copy(&reply, platform)
	reply.Roles = strings.Join(platform.RoleNames, ",")

	return reply, nil
}

func (p platformLogic) ChangePassword(ctx context.Context, request *types.ChangePasswordRequest) error {
	platform, err := p.iDao.GetByID(ctx, request.ID)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			return ecode.ErrGetByIDPlatform.Err()
		}
		return err
	}
	ok := gocrypto.VerifyPassword(request.OldPassword, platform.Password)
	if !ok {
		return ecode.ErrPassword.Err()
	}

	form := &types.UpdatePlatformByIDRequest{}
	form.ID = request.ID
	form.Password = request.NewPassword

	return p.UpdateByID(ctx, form)
}

func (p platformLogic) ResetPassword(ctx context.Context, request *types.ResetPasswordRequest) error {
	if _, err := p.iDao.GetByID(ctx, request.ID); err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			return ecode.ErrGetByIDPlatform.Err()
		}
		return err
	}
	form := &types.UpdatePlatformByIDRequest{}
	form.ID = request.ID
	form.Password = request.Password

	return p.UpdateByID(ctx, form)
}

func encryptMobile(mobile string) string {
	if mobile == "" {
		return ""
	}
	hash, _ := gocrypto.AesEncrypt([]byte(mobile))
	return base64.StdEncoding.EncodeToString(hash)
}

func decryptMobile(mobile string) string {
	if mobile == "" {
		return ""
	}
	hash, _ := base64.StdEncoding.DecodeString(mobile)
	str, _ := gocrypto.AesDecrypt(hash)
	return string(str)
}

func convertPassword(password string) string {
	if password == "" {
		return ""
	}
	hash, _ := gocrypto.HashAndSaltPassword(password)
	return hash
}

func convertPlatform(platform *model.Platform, roleCodes map[uint64]string) (*types.PlatformObjDetail, error) {
	data := &types.PlatformObjDetail{}
	err := copier.Copy(data, platform)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here
	data.Mobile = decryptMobile(data.Mobile)
	data.RoleNames = make([]string, 0)
	for _, roleId := range platform.RoleID {
		if roleName, ok := roleCodes[roleId]; ok {
			data.RoleNames = append(data.RoleNames, roleName)
		}
	}

	return data, nil
}

func (p platformLogic) convertPlatforms(ctx context.Context, fromValues []*model.Platform) ([]*types.PlatformObjDetail, error) {
	var (
		roleIds  []uint64
		toValues []*types.PlatformObjDetail
	)
	for _, v := range fromValues {
		roleIds = append(roleIds, v.RoleID...)
	}

	roleCodes := map[uint64]string{}
	roles, _ := p.iRoleDao.GetByIDs(ctx, roleIds)
	if len(roles) > 0 {
		for _, role := range roles {
			roleCodes[role.ID] = role.Name
		}
	}

	for _, v := range fromValues {
		data, err := convertPlatform(v, roleCodes)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
