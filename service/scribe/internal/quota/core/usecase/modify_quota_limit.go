package usecase

import (
	"context"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/service/scribe/internal/quota/core/model"
)

type ModifyQuotaLimitInput struct {
	ID       string
	NewState model.ModifyQuotaLimitInput
}

func (uc *UseCase) ModifyQuotaLimit(ctx context.Context, in ModifyQuotaLimitInput) error {
	quotaLimit, err := uc.quotaLimitRepo.FindByID(ctx, in.ID)
	if err != nil {
		return custom_err.NewErrDatasourceOperationFailed("find quota limit by id", err)
	}
	if quotaLimit.ID == "" {

		return custom_err.NewErrResourceNotFound("quota limit")
	}

	if err := quotaLimit.Modify(in.NewState); err != nil {
		return err
	}

	if err := uc.quotaLimitRepo.Save(ctx, quotaLimit); err != nil {
		return custom_err.NewErrDatasourceOperationFailed("save quota limit", err)
	}

	return nil
}
