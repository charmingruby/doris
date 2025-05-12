package usecase

import (
	"context"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/service/scribe/internal/quota/core/model"
	"github.com/charmingruby/doris/service/scribe/internal/quota/core/repository"
)

type CreateInitialQuotaUsagesInput struct {
	Tier          string
	CorrelationID string
}

func (uc *UseCase) CreateInitialQuotaUsages(ctx context.Context, in CreateInitialQuotaUsagesInput) error {
	quotas, err := uc.quotaRepo.FindManyByTier(ctx, in.Tier)
	if err != nil {
		return custom_err.NewErrDatasourceOperationFailed("find quotas by tier", err)
	}

	if len(quotas) == 0 {
		return nil
	}

	if err := uc.txManager.Transact(func(tx repository.TransactionManager) error {
		for _, quota := range quotas {
			usage := model.NewQuotaUsage(model.QuotaUsageInput{
				CorrelationID: in.CorrelationID,
				QuotaID:       quota.ID,
			})

			if err := tx.QuotaUsageRepo.Create(ctx, *usage); err != nil {
				return custom_err.NewErrDatasourceOperationFailed("create quota usage", err)
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
