package validation

import (
	"github.com/go-playground/validator/v10"
)

type goValidator struct {
	client *validator.Validate
}

// @note: returning first invalid error
func (v *goValidator) Validate(i interface{}) error {
	err := v.client.Struct(i)
	if err == nil {
		return nil
	}

	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return &ValidationError{
			Message: err.Error(),
		}
	}

	if len(errs) == 0 {
		return &ValidationError{
			Message: err.Error(),
		}
	}

	return &ValidationError{
		Message: errs[0].Error(),
	}
}

func NewGoValidator() *goValidator {
	return &goValidator{
		client: validator.New(),
	}
}
