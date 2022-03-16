package validator

import (
	enloc "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entrans "github.com/go-playground/validator/v10/translations/en"
	"reflect"
)

var v *validate

type validate struct {
	validate *validator.Validate
	trans    ut.Translator
}

type ValidationError struct {
	Source string `json:"source"`
	Detail string `json:"detail"`
}

func init() {
	en := enloc.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")

	v = &validate{validator.New(), trans}

	_ = entrans.RegisterDefaultTranslations(v.validate, v.trans)
}

// ValidateStruct Receives and struct and check if is valid from given rules
func ValidateStruct(obj interface{}) []ValidationError {
	value := reflect.ValueOf(obj)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}

	if valueType == reflect.Struct {
		if err := v.validate.Struct(obj); err != nil {
			validationErrs := err.(validator.ValidationErrors)

			return formatError(validationErrs)
		}
	}

	return nil
}

// ValidateVar Validates a field following the given rules (tag)
func ValidateVar(field interface{}, tag string) []ValidationError {
	if err := v.validate.Var(field, tag); err != nil {
		validationErrs := err.(validator.ValidationErrors)

		return formatError(validationErrs)
	}

	return nil
}

// formatError Formats the error in a more friendly way
func formatError(validationErrs validator.ValidationErrors) (errs []ValidationError) {
	for _, validationErr := range validationErrs {
		vErr := ValidationError{
			Source: validationErr.Field(),
			Detail: validationErr.Translate(v.trans),
		}
		errs = append(errs, vErr)
	}

	return errs
}
