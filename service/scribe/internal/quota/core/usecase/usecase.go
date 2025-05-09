package usecase

import (
	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/service/scribe/internal/quota/core/repository"
)

type UseCase struct {
	logger    *instrumentation.Logger
	quotaRepo repository.QuotaRepository
}

func New(logger *instrumentation.Logger, quotaRepo repository.QuotaRepository) *UseCase {
	return &UseCase{
		logger:    logger,
		quotaRepo: quotaRepo,
	}
}
