package custom_err

import "fmt"

// ErrExternalService is an error to handle external service errors.
type ErrExternalService struct {
	originalErr error
}

func (e *ErrExternalService) Error() string {
	return fmt.Sprintf("external service error: %s", e.originalErr.Error())
}

func NewErrExternalService(originalErr error) error {
	return &ErrExternalService{originalErr: originalErr}
}
