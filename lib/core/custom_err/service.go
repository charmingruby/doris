package custom_err

import (
	"fmt"
)

type ErrResourceAlreadyExists struct {
	resource string
}

func (e *ErrResourceAlreadyExists) Error() string {
	return fmt.Sprintf("resource %s already exists", e.resource)
}

func NewErrResourceAlreadyExists(resource string) error {
	return &ErrResourceAlreadyExists{resource: resource}
}
