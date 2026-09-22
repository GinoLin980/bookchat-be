package customvalidator

import "github.com/go-playground/validator/v10"

type CustomValidator struct{ V *validator.Validate }

func (cv *CustomValidator) Validate(i any) error {
	return cv.V.Struct(i)
}
