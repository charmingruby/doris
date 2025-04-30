package persistence

type TransactionManager[T any] interface {
	Transact(func(repos T) error) error
}
