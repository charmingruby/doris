package custom_err

import "fmt"

// Represents an invalid entity.
type ErrInvalidEntity struct {
	entity string
}

func (e *ErrInvalidEntity) Error() string {
	return fmt.Sprintf("invalid entity: %s", e.entity)
}

func NewErrInvalidEntity(entity string) error {
	return &ErrInvalidEntity{entity: entity}
}
