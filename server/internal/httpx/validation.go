package httpx

import (
	"errors"
	"reflect"

	"github.com/go-playground/validator/v10"
)

func NewValdiator() *validator.Validate {
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		if name := fld.Tag.Get("json"); name != "-" {
			return name
		}
		return ""
	})
	return v
}

func ExtractValidationError(err error) HError {
	hErr := ValidationError
	hErr.Details = collectValidationIssues(err)
	return hErr
}

type fieldViolation struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Param string `json:"param,omitempty"`
}

func collectValidationIssues(err error) []fieldViolation {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return nil
	}

	out := make([]fieldViolation, len(validationErrors))
	for i, vErr := range validationErrors {
		out[i] = fieldViolation{
			Field: vErr.Field(),
			Tag:   vErr.Tag(),
			Param: vErr.Param(),
		}
	}

	return out
}
