package custom_err

import (
	"fmt"
)

// Represents a resource that already exists by a given duplicate field.
type ErrResourceAlreadyExists struct {
	resource string
}

func (e *ErrResourceAlreadyExists) Error() string {
	return fmt.Sprintf("%s already exists", e.resource)
}

func NewErrResourceAlreadyExists(resource string) error {
	return &ErrResourceAlreadyExists{resource: resource}
}

// Represents a resource that was not found.
type ErrResourceNotFound struct {
	resource string
}

func (e *ErrResourceNotFound) Error() string {
	return fmt.Sprintf("%s not found", e.resource)
}

func NewErrResourceNotFound(resource string) error {
	return &ErrResourceNotFound{resource: resource}
}

// Represents an invalid confirmation code.
type ErrInvalidConfirmationCode struct {
	reason string
}

func (e *ErrInvalidConfirmationCode) Error() string {
	return fmt.Sprintf("invalid confirmation code: %s", e.reason)
}

func NewErrInvalidConfirmationCode(reason string) error {
	return &ErrInvalidConfirmationCode{reason: reason}
}

// Represents an api key that is already confirmed.
type ErrAPIKeyAlreadyConfirmed struct{}

func (e *ErrAPIKeyAlreadyConfirmed) Error() string {
	return "api key already confirmed"
}

func NewErrAPIKeyAlreadyConfirmed() error {
	return &ErrAPIKeyAlreadyConfirmed{}
}
