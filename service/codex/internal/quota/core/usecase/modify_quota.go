package usecase

import (
	"context"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/service/scribe/internal/quota/core/model"
)

type ModifyQuotaInput struct {
	ID       string
	NewState model.ModifyQuotaInput
}

func (uc *UseCase) ModifyQuota(ctx context.Context, in ModifyQuotaInput) error {
	quota, err := uc.quotaRepo.FindByID(ctx, in.ID)
	if err != nil {
		return custom_err.NewErrDatasourceOperationFailed("find quota by id", err)
	}

	if quota.ID == "" {
		return custom_err.NewErrResourceNotFound("quota")
	}

	if err := quota.Modify(in.NewState); err != nil {
		return err
	}

	if err := uc.quotaRepo.Save(ctx, quota); err != nil {
		return custom_err.NewErrDatasourceOperationFailed("save quota", err)
	}

	return nil
}
