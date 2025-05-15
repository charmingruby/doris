package usecase

import (
	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/lib/persistence"
	"github.com/charmingruby/doris/service/codex/internal/quota/core/repository"
)

type UseCase struct {
	logger         *instrumentation.Logger
	quotaRepo      repository.QuotaRepository
	quotaUsageRepo repository.QuotaUsageRepository
	txManager      persistence.TransactionManager[repository.TransactionManager]
}

func New(
	logger *instrumentation.Logger,
	quotaRepo repository.QuotaRepository,
	quotaUsageRepo repository.QuotaUsageRepository,
	txManager persistence.TransactionManager[repository.TransactionManager],
) *UseCase {
	return &UseCase{
		logger:         logger,
		quotaRepo:      quotaRepo,
		quotaUsageRepo: quotaUsageRepo,
		txManager:      txManager,
	}
}
