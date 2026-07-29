package admin

import (
	"admin/internal/ecode"
	logic "admin/internal/logic/admin"
	"admin/pkg/gin/handlerfunc"

	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/gin/response"
	"github.com/go-dev-frame/sponge/pkg/logger"
)

var _ UploadHandler = (*uploadHandler)(nil)

// UploadHandler defining the handler interface
type UploadHandler interface {
	Local(c *gin.Context)
}

type uploadHandler struct {
	logic logic.UploadLogic
}

// NewUploadHandler creating the handler interface
func NewUploadHandler() UploadHandler {
	return &uploadHandler{
		logic: logic.NewUploadLogic(),
	}
}

// Local upload local file
// @Summary upload local file
// @Description upload local file
// @Tags admin/upload
// @accept json
// @Produce json
// @Param file formData file true "file"
// @Success 200 {object} admin.UploadLocalReply{}
// @Router /admin/v1/upload/local [post]
// @Security BearerAuth
func (h *uploadHandler) Local(c *gin.Context) {
	_, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	data, err := h.logic.Local(ctx, fileHeader)
	if err != nil {
		if ec, ok := handlerfunc.IsErrcode(err); ok {
			response.Error(c, ec)
			return
		}
		logger.Error("Local upload error", logger.Err(err), logger.Any("filename", fileHeader.Filename), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, data)
}
