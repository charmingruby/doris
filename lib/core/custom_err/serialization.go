package custom_err

import "fmt"

type ErrSerializationFailed struct {
	originalErr error
}

func (e *ErrSerializationFailed) Error() string {
	return fmt.Sprintf("serialization failed: %s", e.originalErr.Error())
}

func NewErrSerializationFailed(originalErr error) error {
	return &ErrSerializationFailed{originalErr: originalErr}
}
