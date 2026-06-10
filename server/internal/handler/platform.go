package handler

import (
	"admin/internal/constant"
	"admin/internal/database"
	"admin/internal/logic"
	"admin/pkg/gin/handlerfunc"
	"admin/pkg/gin/validator"
	"encoding/base64"
	"errors"
	"strings"

	"github.com/go-dev-frame/sponge/pkg/gocrypto"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"

	"admin/internal/ecode"
	"admin/internal/model"
	"admin/internal/types"

	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/gin/response"
	"github.com/go-dev-frame/sponge/pkg/logger"
)

var _ PlatformHandler = (*platformHandler)(nil)

// PlatformHandler defining the handler interface
type PlatformHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
	Me(c *gin.Context)
	GetProfile(c *gin.Context)
	UpdateProfile(c *gin.Context)
	ChangePassword(c *gin.Context)
	ResetPassword(c *gin.Context)
}

type platformHandler struct {
	logic logic.PlatformLogic
}

// NewPlatformHandler creating the handler interface
func NewPlatformHandler() PlatformHandler {
	return &platformHandler{
		logic: logic.NewPlatformLogic(),
	}
}

// Create a record
// @Summary create platform
// @Description submit information to create platform
// @Tags platform
// @accept json
// @Produce json
// @Param data body types.CreatePlatformRequest true "platform information"
// @Success 200 {object} types.CreatePlatformReply{}
// @Router /api/v1/platform [post]
// @Security BearerAuth
func (h *platformHandler) Create(c *gin.Context) {
	form := &types.CreatePlatformRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		response.Error(c, ecode.InvalidParams.RewriteMsg(validator.GetValidatorErrorMsg(err)))
		return
	}

	ctx := middleware.WrapCtx(c)
	id, err := h.logic.Create(ctx, form)
	if err != nil {
		if ec, ok := handlerfunc.IsErrcode(err); ok {
			response.Error(c, ec)
			return
		}
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": id})
}

// DeleteByID delete a record by id
// @Summary delete platform
// @Description delete platform by id
// @Tags platform
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeletePlatformByIDReply{}
// @Router /api/v1/platform/{id} [delete]
// @Security BearerAuth
func (h *platformHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := handlerfunc.GetIdFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	err := h.logic.DeleteByID(ctx, id)
	if err != nil {
		if ec, ok := handlerfunc.IsErrcode(err); ok {
			response.Error(c, ec)
			return
		}
		logger.Error("DeleteByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// UpdateByID update information by id
// @Summary update platform
// @Description update platform information by id
// @Tags platform
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdatePlatformByIDRequest true "platform information"
// @Success 200 {object} types.UpdatePlatformByIDReply{}
// @Router /api/v1/platform/{id} [put]
// @Security BearerAuth
func (h *platformHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := handlerfunc.GetIdFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdatePlatformByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		response.Error(c, ecode.InvalidParams.RewriteMsg(validator.GetValidatorErrorMsg(err)))
		return
	}
	form.ID = id

	ctx := middleware.WrapCtx(c)
	err = h.logic.UpdateByID(ctx, form)
	if err != nil {
		if ec, ok := handlerfunc.IsErrcode(err); ok {
			response.Error(c, ec)
			return
		}
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get platform detail
// @Description get platform detail by id
// @Tags platform
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetPlatformByIDReply{}
// @Router /api/v1/platform/{id} [get]
// @Security BearerAuth
func (h *platformHandler) GetByID(c *gin.Context) {
	_, id, isAbort := handlerfunc.GetIdFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	data, err := h.logic.GetByID(ctx, id)
	if err != nil {
		if ec, ok := handlerfunc.IsErrcode(err); ok {
			response.Error(c, ec)
			return
		}
		logger.Error("GetByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, data)
}

// List of records by query parameters
// @Summary list of platforms by query parameters
// @Description list of platforms by paging and conditions
// @Tags platform
// @accept json
// @Produce json
// @Param request query types.ListPlatformsRequest true "query parameters"
// @Success 200 {object} types.ListPlatformsReply{}
// @Router /api/v1/platform [get]
// @Security BearerAuth
func (h *platformHandler) List(c *gin.Context) {
	request := &types.ListPlatformsRequest{}
	err := c.ShouldBindQuery(request)
	if err != nil {
		response.Error(c, ecode.InvalidParams.RewriteMsg(validator.GetValidatorErrorMsg(err)))
		return
	}

	ctx := middleware.WrapCtx(c)
	data, total, err := h.logic.List(ctx, request)
	if err != nil {
		if ec, ok := handlerfunc.IsErrcode(err); ok {
			response.Error(c, ec)
			return
		}
		logger.Error("List error", logger.Err(err), logger.Any("request", request), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{
		"list":  data,
		"total": total,
	})
}

// Me of records
// @Summary current information
// @Description current information
// @Tags platform
// @accept json
// @Produce json
// @Success 200 {object} types.MeReply{}
// @Router /api/v1/platform/me [get]
// @Security BearerAuth
func (h *platformHandler) Me(c *gin.Context) {
	ctx := middleware.WrapCtx(c)
	id := c.GetUint64("id")
	data, err := h.logic.Me(ctx, id)
	if err != nil {
		if ec, ok := handlerfunc.IsErrcode(err); ok {
			response.Error(c, ec)
			return
		}
		logger.Error("Me error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}
	response.Success(c, data)
}

func convertPlatform(platform *model.Platform, roleCodes map[uint64]string) (*types.PlatformListPage, error) {
	data := &types.PlatformListPage{}
	err := copier.Copy(data, platform)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	var (
		roleNames []string
	)

	for _, roleId := range platform.RoleID {
		if roleName, ok := roleCodes[roleId]; ok {
			roleNames = append(roleNames, roleName)
		}
	}
	data.RoleNames = roleNames

	data.Mobile = decryptMobile(data.Mobile)
	return data, nil
}

func (h *platformHandler) convertPlatforms(c *gin.Context, fromValues []*model.Platform) ([]*types.PlatformListPage, error) {
	var (
		roleIds  []uint64
		toValues []*types.PlatformListPage
	)
	for _, v := range fromValues {
		roleIds = append(roleIds, v.RoleID...)
	}

	roleCodes := map[uint64]string{}
	roles, _ := h.iRoleDao.GetByIDs(c, roleIds)
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

func convertPassword(password string) string {
	if password == "" {
		return ""
	}
	hash, _ := gocrypto.HashAndSaltPassword(password)
	return hash
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

// GetProfile get me information
// @Summary current information
// @Description current information
// @Tags platform
// @accept json
// @Produce json
// @Success 200 {object} types.ProfileReply{}
// @Router /api/v1/platform/profile [get]
// @Security BearerAuth
func (h *platformHandler) GetProfile(c *gin.Context) {
	ctx := middleware.WrapCtx(c)
	id := c.GetUint64("id")
	platform, err := h.iDao.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			logger.Warn("GetByID not found", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.NotFound)
		} else {
			logger.Error("GetByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
			response.Output(c, ecode.InternalServerError.ToHTTPCode())
		}
		return
	}

	reply := types.ProfileItem{}
	_ = copier.Copy(&reply, platform)
	reply.Avatar = h.iConfigDao.MakePathByConfig(c, platform.Avatar, constant.ConfigKeyImageDomain)
	reply.Mobile = decryptMobile(reply.Mobile)
	var (
		roleCodes []string
	)
	roles, _ := h.iRoleDao.GetByIDs(c, platform.RoleID)
	if len(roles) > 0 {
		for _, role := range roles {
			roleCodes = append(roleCodes, role.Name)
		}
	}
	reply.Roles = strings.Join(roleCodes, ",")
	response.Success(c, reply)
}

// UpdateProfile update information by self
// @Summary update platform
// @Description update platform information by self
// @Tags platform
// @accept json
// @Produce json
// @Param data body types.UpdatePlatformByIDRequest true "platform information"
// @Success 200 {object} types.Result{}
// @Router /api/v1/platform/profile [put]
// @Security BearerAuth
func (h *platformHandler) UpdateProfile(c *gin.Context) {
	form := &types.UpdatePlatformByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = c.GetUint64("id")

	platform := &model.Platform{}
	err = copier.Copy(platform, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDPlatform)
		return
	}

	platform.Mobile = encryptMobile(form.Mobile)
	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, platform)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// ChangePassword change password by self
// @Summary change password by self
// @Description change password by self
// @Tags platform
// @accept json
// @Produce json
// @Param data body types.ChangePasswordRequest true "platform information"
// @Success 200 {object} types.Result{}
// @Router /api/v1/platform/password [put]
// @Security BearerAuth
func (h *platformHandler) ChangePassword(c *gin.Context) {
	request := &types.ChangePasswordRequest{}
	err := c.ShouldBindJSON(request)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form := &model.Platform{}
	form.ID = c.GetUint64("id")

	platform, platformErr := h.iDao.GetByID(c, form.ID)
	if platformErr != nil {
		if errors.Is(platformErr, database.ErrRecordNotFound) {
			logger.Warn("GetByID not found", logger.Err(platformErr), logger.Any("id", form.ID), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.NotFound)
		} else {
			logger.Error("GetByID error", logger.Err(platformErr), logger.Any("id", form.ID), middleware.GCtxRequestIDField(c))
			response.Output(c, ecode.InternalServerError.ToHTTPCode())
		}
		return
	}

	ok := gocrypto.VerifyPassword(request.OldPassword, platform.Password)
	if !ok {
		response.Error(c, ecode.ErrPassword)
		return
	}

	form.Password = convertPassword(request.NewPassword)

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, form)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("request", request), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// ResetPassword reset password by self
// @Summary reset password by self
// @Description reset password by self
// @Tags platform
// @accept json
// @Produce json
// @Param data body types.ResetPasswordRequest true "platform information"
// @Success 200 {object} types.Result{}
// @Router /api/v1/platform/password/reset [put]
// @Security BearerAuth
func (h *platformHandler) ResetPassword(c *gin.Context) {
	request := &types.ResetPasswordRequest{}
	err := c.ShouldBindJSON(request)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	_, err = h.iDao.GetByID(ctx, request.ID)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			logger.Warn("GetByID not found", logger.Err(err), logger.Any("id", request.ID), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.NotFound)
		} else {
			logger.Error("GetByID error", logger.Err(err), logger.Any("id", request.ID), middleware.GCtxRequestIDField(c))
			response.Output(c, ecode.InternalServerError.ToHTTPCode())
		}
		return
	}

	form := &model.Platform{}
	form.ID = request.ID
	form.Password = convertPassword(request.Password)

	err = h.iDao.UpdateByID(ctx, form)
	if err != nil {
		logger.Error("ResetPassword error", logger.Err(err), logger.Any("request", request), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}
