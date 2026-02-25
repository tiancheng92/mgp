package mgp

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/tiancheng92/mgp/errors"
	"github.com/tiancheng92/mgp/errors/default_error_code"
)

var (
	translator ut.Translator
	validate   *validator.Validate
)

func init() {
	if _, ok := binding.Validator.Engine().(*validator.Validate); !ok {
		panic("binding validate engine failed")
	}

	validate = binding.Validator.Engine().(*validator.Validate)

	zhT := zh.New()
	enT := en.New()
	uni := ut.New(enT, zhT)

	var found bool
	translator, found = uni.GetTranslator("zh")
	if !found {
		panic("translator not found")
	}

	if err := zhTranslations.RegisterDefaultTranslations(validate, translator); err != nil {
		panic(fmt.Sprintf("register default translations failed: %+v", err))
	}
}

func registerTranslator(tag string, msg string) validator.RegisterTranslationsFunc {
	return func(trans ut.Translator) error {
		return trans.Add(tag, msg, true)
	}
}

func RegisterValidate(tag, kind, msg string, validatorFunc func(validator.FieldLevel) bool, translateFunc func(ut.Translator, validator.FieldError) string) {
	if validatorFunc != nil {
		if err := validate.RegisterValidation(tag, validatorFunc); err != nil {
			panic(fmt.Sprintf("register validation failed: %+v", err))
		}
	}

	if translateFunc == nil {
		translateFunc = defaultTranslateFunc
	}

	if err := validate.RegisterTranslation(tag, translator, registerTranslator(fmt.Sprintf("%s-%s", tag, kind), msg), translateFunc); err != nil {
		panic(fmt.Sprintf("register validation translation failed: %+v", err))
	}
}

func defaultTranslateFunc(translator ut.Translator, fieldError validator.FieldError) string {
	var kindStr string
	kind := fieldError.Kind()
	if kind == reflect.Ptr {
		kind = fieldError.Type().Elem().Kind()
	}

	switch kind {
	case reflect.String:
		kindStr = "string"
	case reflect.Slice, reflect.Map, reflect.Array:
		kindStr = "object"
	default:
		kindStr = "number"
	}
	msg, err := translator.T(fmt.Sprintf("%s-%s", fieldError.Tag(), kindStr), fieldError.Field(), fieldError.Param())
	if err != nil {
		panic(fmt.Sprintf("register validation failed: %+v", err))
	}
	return msg
}

func HandleValidationErr(err error) error {
	if validationErr, ok := errors.AsType[validator.ValidationErrors](err); ok {
		errList := make([]string, 0, len(validationErr))
		for _, v := range validationErr.Translate(translator) {
			errList = append(errList, v)
		}
		return errors.WithCode(default_error_code.ErrClientParam, strings.Join(errList, "; "))
	}

	return errors.WithCode(default_error_code.ErrClientParam, err)
}
