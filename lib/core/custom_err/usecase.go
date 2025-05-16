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

// Represents an invalid OTP code.
type ErrInvalidOTPCode struct {
	reason string
}

func (e *ErrInvalidOTPCode) Error() string {
	return fmt.Sprintf("invalid OTP code: %s", e.reason)
}

func NewErrInvalidOTPCode(reason string) error {
	return &ErrInvalidOTPCode{reason: reason}
}

// Represents an api key that is already confirmed.
type ErrAPIKeyAlreadyActivated struct{}

func (e *ErrAPIKeyAlreadyActivated) Error() string {
	return "api key already activated"
}

func NewErrAPIKeyAlreadyActivated() error {
	return &ErrAPIKeyAlreadyActivated{}
}

// Represents an action with insufficient permission.
type ErrInsufficientPermission struct{}

func (e *ErrInsufficientPermission) Error() string {
	return "insufficient permission"
}

func NewErrInsufficientPermission() error {
	return &ErrInsufficientPermission{}
}

// Represents an intent to modify something that does not change nothing.
type ErrNothingToChange struct{}

func (e *ErrNothingToChange) Error() string {
	return "nothing to change"
}

func NewErrNothingToChange() error {
	return &ErrNothingToChange{}
}

// Represents an intent to modify something that does not change nothing.
type ErrOTPGenerationCooldown struct{}

func (e *ErrOTPGenerationCooldown) Error() string {
	return "unable to generate new otp code yet"
}

func NewErrOTPGenerationCooldown() error {
	return &ErrOTPGenerationCooldown{}
}

// Represents a quota that is exceeded.
type ErrQuotaExceeded struct{}

func (e *ErrQuotaExceeded) Error() string {
	return "quota exceeded"
}

func NewErrQuotaExceeded() error {
	return &ErrQuotaExceeded{}
}
