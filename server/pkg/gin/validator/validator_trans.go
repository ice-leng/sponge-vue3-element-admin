package validator

import (
	"errors"
	"reflect"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-dev-frame/sponge/pkg/gin/validator"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	validatorV10 "github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	Trans ut.Translator
)

// InitTrans initialize the translator
func InitTrans() {
	v := validator.Init()
	zh := zh.New()
	uni := ut.New(zh, zh)
	trans, _ := uni.GetTranslator("zh")
	Trans = trans
	_ = zhTranslations.RegisterDefaultTranslations(v.Validate, trans)

	v.Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		label := fld.Tag.Get("label")
		if label == "" {
			return fld.Name
		}
		return label
	})

	binding.Validator = v
}

// GetValidatorErrorMsg get the validator error message
func GetValidatorErrorMsg(err error) string {
	var verrs validatorV10.ValidationErrors
	if errors.As(err, &verrs) {
		return verrs[0].Translate(Trans)
	}
	return err.Error()
}
