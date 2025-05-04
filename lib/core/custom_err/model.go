package custom_err

import "fmt"

// Represents an invalid entity.
type ErrInvalidEntity struct {
	msg string
}

func (e *ErrInvalidEntity) Error() string {
	return fmt.Sprintf("invalid entity: %s", e.msg)
}

func NewErrInvalidEntity(msg string) error {
	return &ErrInvalidEntity{msg: msg}
}
