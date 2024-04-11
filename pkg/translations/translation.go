package translations

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/locales/zh_Hant_TW"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
	"go-gin/pkg/logging"
	"regexp"
)

var Trans ut.Translator

const (
	ErrorArgsDenied = "{0}手机格式不正确"
)

func Setup(lang string) {
	uni := ut.New(en.New(), zh.New(), zh_Hant_TW.New())
	Trans, _ = uni.GetTranslator(lang)
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		switch lang {
		case "zh":
			_ = zhtranslations.RegisterDefaultTranslations(v, Trans)
			break
		case "en":
			_ = entranslations.RegisterDefaultTranslations(v, Trans)
			break
		default:
			_ = zhtranslations.RegisterDefaultTranslations(v, Trans)
			break
		}
	}
	AddValidator()
}

// AddValidator 注册自定义验证器
func AddValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		if err := RegisterValidation(v, "mobile", ErrorArgsDenied, ValidateMobile); err != nil {
			logging.Error("%v", err)
			return
		}
	}
}

func RegisterValidation(v *validator.Validate, tag string, msg string, fn validator.Func) error {
	if err := v.RegisterValidation(tag, fn); err != nil {
		return err
	}
	err := v.RegisterTranslation(tag, Trans, func(ut ut.Translator) error {
		return ut.Add(tag, msg, true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, fe.Field(), fe.Field())
		return t
	})
	if err != nil {
		return err
	}
	return nil
}

// ValidateMobile 校验手机号
func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	ok, _ := regexp.MatchString(`^(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\d{8}$`, mobile)
	if !ok {
		return false
	}
	return true
}
