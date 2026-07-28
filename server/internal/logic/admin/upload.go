package admin

import (
	"admin/internal/constant"
	types "admin/internal/types/admin"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-dev-frame/sponge/pkg/errcode"
	"github.com/go-dev-frame/sponge/pkg/gocrypto"
	"github.com/go-dev-frame/sponge/pkg/gofile"
)

type UploadLogic interface {
	Local(ctx context.Context, fileHeader *multipart.FileHeader) (*types.UploadItem, error)
}

type uploadLogic struct {
	configLogic ConfigLogic
}

func NewUploadLogic() UploadLogic {
	return NewUploadLogicByDAO(NewConfigLogic())
}

func NewUploadLogicByDAO(configLogic ConfigLogic) UploadLogic {
	return &uploadLogic{
		configLogic: configLogic,
	}
}

func (l uploadLogic) Local(ctx context.Context, fileHeader *multipart.FileHeader) (*types.UploadItem, error) {
	ext := filepath.Ext(fileHeader.Filename)
	name := strings.TrimSuffix(fileHeader.Filename, ext)
	newFileName := gocrypto.Md5([]byte(name+time.Now().Format("20060102150405"))) + ext

	path := fmt.Sprintf("%s/%s", "uploads", time.Now().Format("2006-01-02"))
	if !gofile.IsExists(path) {
		if err := gofile.CreateDir(path); err != nil {
			return nil, errcode.NewError(10001, err.Error()).Err()
		}
	}
	filePath := path + "/" + newFileName
	f, openError := fileHeader.Open()
	if openError != nil {
		return nil, errcode.NewError(10001, openError.Error()).Err()
	}
	defer f.Close()

	out, createErr := os.Create(filePath)
	if createErr != nil {
		return nil, errcode.NewError(10001, createErr.Error()).Err()
	}
	defer out.Close()

	_, copyErr := io.Copy(out, f)
	if copyErr != nil {
		return nil, errcode.NewError(10001, copyErr.Error()).Err()
	}

	return &types.UploadItem{
		Name: newFileName,
		Path: "/" + filePath,
		Url:  l.configLogic.MakePathByConfig(ctx, "/"+filePath, constant.ConfigKeyImageDomain),
	}, nil
}
