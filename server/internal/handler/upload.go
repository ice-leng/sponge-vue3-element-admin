package handler

import (
	"admin/internal/cache"
	"admin/internal/dao"
	"admin/internal/ecode"
	"admin/internal/model"
	"admin/internal/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zhufuyi/sponge/pkg/errcode"
	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/gocrypto"
	"github.com/zhufuyi/sponge/pkg/gofile"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type UploadHandler interface {
	Local(c *gin.Context)
}

type uploadHandler struct {
	iConfigDao dao.ConfigDao
}

func NewUploadHandler() UploadHandler {
	return &uploadHandler{
		iConfigDao: dao.NewConfigDao(
			model.GetDB(),
			cache.NewConfigCache(model.GetCacheType()),
		),
	}
}

// Local upload local file
// @Summary upload local file
// @Description upload local file
// @Tags upload
// @accept json
// @Produce json
// @Param file formData file true "file"
// @Success 200 {object} types.UploadLocalReply{}
// @Router /api/v1/upload/local [post]
// @Security BearerAuth
func (h *uploadHandler) Local(c *gin.Context) {
	_, file, err := c.Request.FormFile("file")
	if err != nil {
		response.Error(c, ecode.InvalidParams)
		return
	}
	ext := filepath.Ext(file.Filename)
	name := strings.TrimSuffix(file.Filename, ext)
	newFileName := gocrypto.Md5([]byte(name+time.Now().Format("20060102150405"))) + ext

	path := fmt.Sprintf("%s/%s", "uploads", time.Now().Format("2006-01-02"))
	if !gofile.IsExists(path) {
		if err := gofile.CreateDir(path); err != nil {
			response.Error(c, errcode.NewError(10001, err.Error()))
			return
		}
	}
	filePath := path + newFileName
	f, openError := file.Open() // 读取文件
	if openError != nil {
		response.Error(c, errcode.NewError(10001, openError.Error()))
		return
	}
	defer f.Close() // 创建文件 defer 关闭

	out, createErr := os.Create(filePath)
	if createErr != nil {
		response.Error(c, errcode.NewError(10001, createErr.Error()))
		return
	}
	defer out.Close() // 创建文件 defer 关闭

	_, copyErr := io.Copy(out, f) // 传输（拷贝）文件
	if copyErr != nil {
		response.Error(c, errcode.NewError(10001, copyErr.Error()))
		return
	}

	// 上传成功
	response.Success(c, types.UploadItem{
		Name: newFileName,
		Path: "/" + filePath,
		Url:  h.iConfigDao.MakePathByConfig(c, "/"+filePath, "imageDomain"),
	})
}
