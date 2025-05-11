package usecase

import (
	"context"
	"errors"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/core/privilege"
	"github.com/charmingruby/doris/service/scribe/internal/quota/core/model"
)

func (s *Suite) Test_ModifyQuotaLimit() {
	validQuota, err := model.NewQuota(model.QuotaInput{
		Tier: privilege.TIER_ROOKIE,
	})
	s.NoError(err)

	validQuotaLimit, err := model.NewQuotaLimit(model.QuotaLimitInput{
		QuotaID:  validQuota.ID,
		Kind:     model.QUOTA_LIMIT_KIND_REQUEST,
		MaxValue: 10,
		Unit:     "request",
	})
	s.NoError(err)

	s.Run("it should be able to modify a quota limit", func() {
		ctx := context.Background()

		err := s.quotaRepo.Create(ctx, *validQuota)
		s.NoError(err)

		err = s.quotaLimitRepo.Create(ctx, *validQuotaLimit)
		s.NoError(err)

		input := ModifyQuotaLimitInput{
			ID: validQuotaLimit.ID,
			NewState: model.ModifyQuotaLimitInput{
				Kind:     model.QUOTA_LIMIT_KIND_DOCUMENT,
				MaxValue: 15,
				Unit:     "document",
				IsActive: false,
			},
		}

		err = s.uc.ModifyQuotaLimit(ctx, input)
		s.NoError(err)

		storedQuotaLimit := s.quotaLimitRepo.Items[0]
		s.Equal(input.NewState.Kind, storedQuotaLimit.Kind)
		s.Equal(input.NewState.MaxValue, storedQuotaLimit.MaxValue)
		s.Equal(input.NewState.Unit, storedQuotaLimit.Unit)
		s.Equal(input.NewState.IsActive, storedQuotaLimit.IsActive)
	})

	s.Run("it should not be able to modify if quota limit does not exist", func() {
		ctx := context.Background()

		input := ModifyQuotaLimitInput{
			ID: "non-existent-id",
			NewState: model.ModifyQuotaLimitInput{
				Kind:     model.QUOTA_LIMIT_KIND_DOCUMENT,
				MaxValue: 15,
				Unit:     "document",
				IsActive: false,
			},
		}

		err := s.uc.ModifyQuotaLimit(ctx, input)
		s.Error(err)

		var notFoundErr *custom_err.ErrResourceNotFound
		s.True(errors.As(err, &notFoundErr), "error should be of type ErrResourceNotFound")
	})

	s.Run("it should not be able to modify if datasource fails", func() {
		ctx := context.Background()

		err := s.quotaRepo.Create(ctx, *validQuota)
		s.NoError(err)

		err = s.quotaLimitRepo.Create(ctx, *validQuotaLimit)
		s.NoError(err)

		s.quotaLimitRepo.IsHealthy = false

		input := ModifyQuotaLimitInput{
			ID: validQuotaLimit.ID,
			NewState: model.ModifyQuotaLimitInput{
				Kind:     model.QUOTA_LIMIT_KIND_DOCUMENT,
				MaxValue: 15,
				Unit:     "document",
				IsActive: false,
			},
		}

		err = s.uc.ModifyQuotaLimit(ctx, input)
		s.Error(err)

		var dsErr *custom_err.ErrDatasourceOperationFailed
		s.True(errors.As(err, &dsErr), "error should be of type ErrDatasourceOperationFailed")
	})

	s.Run("it should not be able to modify if nothing has changed", func() {
		ctx := context.Background()

		err := s.quotaRepo.Create(ctx, *validQuota)
		s.NoError(err)

		err = s.quotaLimitRepo.Create(ctx, *validQuotaLimit)
		s.NoError(err)

		input := ModifyQuotaLimitInput{
			ID: validQuotaLimit.ID,
			NewState: model.ModifyQuotaLimitInput{
				Kind:     validQuotaLimit.Kind,
				MaxValue: validQuotaLimit.MaxValue,
				Unit:     validQuotaLimit.Unit,
				IsActive: validQuotaLimit.IsActive,
			},
		}

		err = s.uc.ModifyQuotaLimit(ctx, input)
		s.Error(err)

		var nothingToChangeErr *custom_err.ErrNothingToChange
		s.True(errors.As(err, &nothingToChangeErr), "error should be of type ErrNothingToChange")
	})
}
