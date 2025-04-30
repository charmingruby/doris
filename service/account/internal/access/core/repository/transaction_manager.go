package repository

type TransactionManager struct {
	OTPRepo    OTPRepository
	APIKeyRepo APIKeyRepository
}
