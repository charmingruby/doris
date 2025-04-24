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
	var errs []error

	err := v.validator.Struct(s)

	if err == nil {
		return nil
	}

	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		errs = append(errs, err)
		return errs
	}

	for _, vErr := range validationErrs {
		fmtErr := fmt.Sprintf("%s does not satisfy %s", vErr.Field(), vErr.Tag())
		errs = append(errs, errors.New(fmtErr))
	}

	return errs
}
