package usecase

import (
	"context"
	"time"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/service/codex/internal/quota/core/repository"
)

func (uc *UseCase) ResetAllQuotaUsages(ctx context.Context) error {
	if err := uc.txManager.Transact(func(tx repository.TransactionManager) error {
		now := time.Now()

		if err := tx.QuotaUsageRepo.UpdateAllCurrentUsages(ctx, now); err != nil {
			return custom_err.NewErrDatasourceOperationFailed("reset all quota usages", err)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
