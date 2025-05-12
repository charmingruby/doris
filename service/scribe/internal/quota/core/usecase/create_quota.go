package usecase

import (
	"context"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/service/scribe/internal/quota/core/model"
)

type CreateQuotaInput struct {
	Tier     string
	Kind     string
	MaxValue int
	Unit     string
}

func (uc *UseCase) CreateQuota(ctx context.Context, in CreateQuotaInput) (string, error) {
	quota, err := uc.quotaRepo.FindByTierAndKind(ctx, in.Tier, in.Kind)
	if err != nil {
		return "", custom_err.NewErrDatasourceOperationFailed("find quota by tier and kind", err)
	}

	if quota.ID != "" {
		return "", custom_err.NewErrResourceAlreadyExists("quota")
	}

	q, err := model.NewQuota(model.QuotaInput{
		Tier:     in.Tier,
		Kind:     in.Kind,
		MaxValue: in.MaxValue,
		Unit:     in.Unit,
	})
	if err != nil {
		return "", custom_err.NewErrInvalidEntity(err.Error())
	}

	if err := uc.quotaRepo.Create(ctx, *q); err != nil {
		return "", custom_err.NewErrDatasourceOperationFailed("create quota", err)
	}

	return q.ID, nil
}
