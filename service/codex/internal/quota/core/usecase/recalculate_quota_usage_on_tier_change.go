package usecase

import (
	"context"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/service/codex/internal/quota/core/model"
	"github.com/charmingruby/doris/service/codex/internal/quota/core/repository"
)

type RecalculateQuotaUsageOnTierChangeInput struct {
	CorrelationID string
	OldTier       string
	NewTier       string
}

type quotaMigration struct {
	correlationID string
	quotaID       string
	quotaUsageID  string
	currentUsage  int
}

func (uc *UseCase) RecalculateQuotaUsageOnTierChange(ctx context.Context, in RecalculateQuotaUsageOnTierChangeInput) error {
	oldQuotas, err := uc.quotaRepo.FindManyByTier(ctx, in.OldTier)
	if err != nil {
		return custom_err.NewErrDatasourceOperationFailed("find quotas by old tier", err)
	}

	newQuotas, err := uc.quotaRepo.FindManyByTier(ctx, in.NewTier)
	if err != nil {
		return custom_err.NewErrDatasourceOperationFailed("find quotas by new tier", err)
	}

	quotasToDeactivate := []quotaMigration{}
	for _, oldQuota := range oldQuotas {
		usage, err := uc.quotaUsageRepo.FindByCorrelationIDAndQuotaID(ctx, in.CorrelationID, oldQuota.ID)
		if err != nil {
			return custom_err.NewErrDatasourceOperationFailed("find quota usage", err)
		}

		if usage.ID != "" {
			quotasToDeactivate = append(quotasToDeactivate, quotaMigration{
				correlationID: in.CorrelationID,
				quotaID:       oldQuota.ID,
				quotaUsageID:  usage.ID,
				currentUsage:  usage.CurrentUsage,
			})
		}
	}

	if err := uc.txManager.Transact(func(tx repository.TransactionManager) error {
		for _, migration := range quotasToDeactivate {
			usage := model.QuotaUsage{
				ID:            migration.quotaUsageID,
				CorrelationID: migration.correlationID,
				QuotaID:       migration.quotaID,
				CurrentUsage:  migration.currentUsage,
				IsActive:      false,
			}

			if err := tx.QuotaUsageRepo.Save(ctx, usage); err != nil {
				return custom_err.NewErrDatasourceOperationFailed("save quota usage", err)
			}
		}

		for _, newQuota := range newQuotas {
			usage := model.NewQuotaUsage(model.QuotaUsageInput{
				CorrelationID: in.CorrelationID,
				QuotaID:       newQuota.ID,
			})

			if err := tx.QuotaUsageRepo.Create(ctx, *usage); err != nil {
				return custom_err.NewErrDatasourceOperationFailed("create new quota usage", err)
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
