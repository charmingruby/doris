package repository

type TransactionManager struct {
	QuotaRepo      QuotaRepository
	QuotaUsageRepo QuotaUsageRepository
}
