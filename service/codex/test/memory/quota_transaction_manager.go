package memory

import (
	"github.com/charmingruby/doris/service/scribe/internal/quota/core/model"
	"github.com/charmingruby/doris/service/scribe/internal/quota/core/repository"
)

type TransactionManager struct {
	quotaRepo      *QuotaRepository
	quotaUsageRepo *QuotaUsageRepository
}

func NewQuotaTransactionManager(quotaRepo *QuotaRepository, quotaUsageRepo *QuotaUsageRepository) *TransactionManager {
	return &TransactionManager{
		quotaRepo:      quotaRepo,
		quotaUsageRepo: quotaUsageRepo,
	}
}

func (t *TransactionManager) Transact(fn func(repos repository.TransactionManager) error) error {
	quotaUsageSnapshot := make([]model.QuotaUsage, len(t.quotaUsageRepo.Items))
	copy(quotaUsageSnapshot, t.quotaUsageRepo.Items)

	quotaSnapshot := make([]model.Quota, len(t.quotaRepo.Items))
	copy(quotaSnapshot, t.quotaRepo.Items)

	err := fn(repository.TransactionManager{
		QuotaUsageRepo: t.quotaUsageRepo,
		QuotaRepo:      t.quotaRepo,
	})

	if err != nil {
		t.quotaUsageRepo.Items = quotaUsageSnapshot
		t.quotaRepo.Items = quotaSnapshot

		return err
	}

	return nil
}
