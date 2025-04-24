package validation

import (
	"errors"
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

func (v *Validator) Validate(s any) []error {
	err := v.validator.Struct(s)

	if err == nil {
		return nil
	}

	return v.UnwrapValidationErr(err)
}

func (v *Validator) UnwrapValidationErr(err error) []error {
	var wrappedErrs []error

	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		wrappedErrs = append(wrappedErrs, err)
		return wrappedErrs
	}

	for _, vErr := range validationErrs {
		fmtErr := fmt.Sprintf("%s does not satisfy %s", vErr.Field(), vErr.Tag())

		wrappedErrs = append(wrappedErrs, errors.New(fmtErr))
	}

	return wrappedErrs
}
