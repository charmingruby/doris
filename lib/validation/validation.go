package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

func (v *Validator) Validate(s any) []string {
	err := v.validator.Struct(s)

	if err == nil {
		return nil
	}

	return v.UnwrapValidationErr(err)
}

func (v *Validator) UnwrapValidationErr(err error) []string {
	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return []string{err.Error()}
	}

	reasonsWrapper := make([]string, 0, len(validationErrs))

	for _, vErr := range validationErrs {
		reason := fmt.Sprintf("field `%s` does not satisfy %s rule", vErr.Field(), vErr.Tag())

		reasonsWrapper = append(reasonsWrapper, reason)
	}

	return reasonsWrapper
}
