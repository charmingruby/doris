package usecase

import (
	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/service/scribe/internal/quota/core/repository"
)

type UseCase struct {
	logger              *instrumentation.Logger
	quotaRepo           repository.QuotaRepository
	quotaLimitRepo      repository.QuotaLimitRepository
	quotaLimitUsageRepo repository.QuotaLimitUsageRepository
}

func New(
	logger *instrumentation.Logger,
	quotaRepo repository.QuotaRepository,
	quotaLimitRepo repository.QuotaLimitRepository,
	quotaLimitUsageRepo repository.QuotaLimitUsageRepository,
) *UseCase {
	return &UseCase{
		logger:              logger,
		quotaRepo:           quotaRepo,
		quotaLimitRepo:      quotaLimitRepo,
		quotaLimitUsageRepo: quotaLimitUsageRepo,
	}
}
