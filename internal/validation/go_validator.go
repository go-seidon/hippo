package validation

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type CustomValidator = func(fl validator.FieldLevel) bool
type CustomTagName = func(field reflect.StructField) string

type goValidator struct {
	client     *validator.Validate
	translator ut.Translator
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
		Message: errs[0].Translate(v.translator),
	}
}

func NewGoValidator() *goValidator {
	en := en.New()
	uniTrans := ut.New(en, en)
	translator, _ := uniTrans.GetTranslator("en")

	client := validator.New()

	en_translations.RegisterDefaultTranslations(client, translator)
	client.RegisterTagNameFunc(LabelTagName())

	return &goValidator{
		client:     client,
		translator: translator,
	}
}

func LabelTagName() CustomTagName {
	return func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("label"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	}
}
