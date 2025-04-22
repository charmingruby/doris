package custom_err

import (
	"fmt"
)

type ErrDatasourceOperationFailed struct {
	operation   string
	originalErr error
}

func (e *ErrDatasourceOperationFailed) Error() string {
	return fmt.Sprintf("datasource operation failed: `%s`, %s", e.operation, e.originalErr.Error())
}

func NewErrDatasourceOperationFailed(operation string, originalErr error) error {
	return &ErrDatasourceOperationFailed{operation: operation, originalErr: originalErr}
}
