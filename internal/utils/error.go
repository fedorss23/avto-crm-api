package utils

import (
	"github.com/go-playground/validator/v10"
	"errors"
	"fmt"
)

func ParseValidationErrors(err error) map[string]string {
	errs := make(map[string]string)

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fe := range ve {
			errs[fe.Field()] = validationMessage(fe)
		}
		return errs
	}

	errs["_"] = err.Error()
	return errs
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "field required"
	case "email":
		return "incorrect email"
	case "min":
		return fmt.Sprintf("mininal length/value: %s", fe.Param())
	case "max":
		return fmt.Sprintf("maximal length/value: %s", fe.Param())
	default:
		return fmt.Sprintf("invalid value (%s)", fe.Tag())
	}
}