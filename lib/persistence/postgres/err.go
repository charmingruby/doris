package postgres

import "fmt"

type PersistenceErr struct {
	Message string `json:"message"`
}

func newPersistenceErr(originalErr error, action, datasource string) *PersistenceErr {
	return &PersistenceErr{
		Message: fmt.Sprintf("%s persistence error on `%s`: %s", datasource, action, originalErr.Error()),
	}
}

func (e *PersistenceErr) Error() string {
	return e.Message
}

func NewPreparationErr(queryName string, repository string, err error) *PersistenceErr {
	preparationErr := fmt.Errorf("unable to prepare the query:`%s` on %s repository, original err: %s", queryName, repository, err.Error())
	return newPersistenceErr(preparationErr, "prepare", "postgres")
}

func NewStatementNotPreparedErr(queryName string, repository string) *PersistenceErr {
	preparationErr := fmt.Errorf("query `%s` is not prepared on %s repository", queryName, repository)
	return newPersistenceErr(preparationErr, "statement not prepared", "postgres")
}
