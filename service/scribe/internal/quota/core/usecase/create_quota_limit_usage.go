package usecase

import (
	"context"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/service/scribe/internal/quota/core/model"
)

type CreateQuotaLimitUsageInput struct {
	QuotaLimitID  string
	CorrelationID string
}

func (uc *UseCase) CreateQuotaLimitUsage(ctx context.Context, in CreateQuotaLimitUsageInput) (string, error) {
	quotaLimit, err := uc.quotaLimitRepo.FindByID(ctx, in.QuotaLimitID)

	if err != nil {
		return "", custom_err.NewErrDatasourceOperationFailed("find quota limit by id", err)
	}

	if quotaLimit.ID == "" {
		return "", custom_err.NewErrResourceNotFound("quota limit")
	}

	existingQuotaLimitUsage, err := uc.quotaLimitUsageRepo.FindByQuotaLimitIDAndCorrelationID(ctx, in.QuotaLimitID, in.CorrelationID)

	if err != nil {
		return "", custom_err.NewErrDatasourceOperationFailed("find quota limit usage by quota limit id and correlation id", err)
	}

	if existingQuotaLimitUsage.ID != "" {
		return "", custom_err.NewErrResourceAlreadyExists("quota limit usage")
	}

	quotaLimitUsage := model.NewQuotaLimitUsage(model.QuotaLimitUsageInput{
		CorrelationID: in.CorrelationID,
		QuotaLimitID:  in.QuotaLimitID,
	})

	if err := uc.quotaLimitUsageRepo.Create(ctx, *quotaLimitUsage); err != nil {
		return "", custom_err.NewErrDatasourceOperationFailed("create quota limit usage", err)
	}

	return quotaLimitUsage.ID, nil
}
