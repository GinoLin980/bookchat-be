package customvalidator

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct{ V *validator.Validate }

func (cv *CustomValidator) Validate(i any) error {
	return cv.V.Struct(i)
}

func InitValidator() *validator.Validate {
	v := validator.New()

	// 2. 註冊 TagNameFunc，讓驗證器能夠讀取 json 標籤 (例如 "book_author")
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.Split(fld.Tag.Get("json"), ",")[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return v
}
