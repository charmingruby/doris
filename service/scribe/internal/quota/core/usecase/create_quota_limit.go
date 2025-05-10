package usecase

import (
	"context"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/service/scribe/internal/quota/core/model"
)

type CreateQuotaLimitInput struct {
	QuotaID  string
	Kind     string
	MaxValue int
	Unit     string
}

func (uc *UseCase) CreateQuotaLimit(ctx context.Context, in CreateQuotaLimitInput) (string, error) {
	quota, err := uc.quotaRepo.FindByID(ctx, in.QuotaID)

	if err != nil {
		return "", custom_err.NewErrDatasourceOperationFailed("find quota by id", err)
	}

	if quota.ID == "" {
		return "", custom_err.NewErrResourceNotFound("quota")
	}

	existingQuotaLimit, err := uc.quotaLimitRepo.FindByQuotaIDAndKind(ctx, in.QuotaID, in.Kind)

	if err != nil {
		return "", custom_err.NewErrDatasourceOperationFailed("find quota limit by quota id and kind", err)
	}

	if existingQuotaLimit.ID != "" {
		return "", custom_err.NewErrResourceAlreadyExists("quota limit")
	}

	quotaLimit, err := model.NewQuotaLimit(model.QuotaLimitInput{
		QuotaID:  in.QuotaID,
		Kind:     in.Kind,
		MaxValue: in.MaxValue,
		Unit:     in.Unit,
	})
	if err != nil {
		return "", custom_err.NewErrInvalidEntity(err.Error())
	}

	if err := uc.quotaLimitRepo.Create(ctx, *quotaLimit); err != nil {
		return "", custom_err.NewErrDatasourceOperationFailed("create quota limit", err)
	}

	return quotaLimit.ID, nil
}
