package usecase

import (
	"context"
	"errors"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/core/privilege"
	"github.com/charmingruby/doris/service/scribe/internal/quota/core/model"
	"github.com/charmingruby/doris/service/scribe/internal/shared/core/kind"
)

func (s *Suite) Test_ModifyQuota() {
	s.Run("it should be able to modify a quota", func() {
		ctx := context.Background()

		quota, err := model.NewQuota(model.QuotaInput{
			Tier:     privilege.TIER_ROOKIE,
			Kind:     kind.QUOTA_LIMIT_REQUEST,
			MaxValue: 100,
			Unit:     "request",
		})
		s.NoError(err)

		err = s.quotaRepo.Create(ctx, *quota)
		s.NoError(err)

		input := ModifyQuotaInput{
			ID: quota.ID,
			NewState: model.ModifyQuotaInput{
				Tier:     privilege.TIER_PRO,
				Status:   model.QUOTA_STATUS_ENABLED,
				MaxValue: 200,
				Unit:     "request",
			},
		}

		err = s.uc.ModifyQuota(ctx, input)
		s.NoError(err)

		storedQuota := s.quotaRepo.Items[0]
		s.Equal(input.NewState.Tier, storedQuota.Tier)
		s.Equal(input.NewState.Status, storedQuota.Status)
		s.Equal(input.NewState.MaxValue, storedQuota.MaxValue)
		s.Equal(input.NewState.Unit, storedQuota.Unit)
	})

	s.Run("it should not be able to modify if quota does not exist", func() {
		ctx := context.Background()

		input := ModifyQuotaInput{
			ID: "non-existent-id",
			NewState: model.ModifyQuotaInput{
				Tier:     privilege.TIER_PRO,
				Status:   model.QUOTA_STATUS_ENABLED,
				MaxValue: 200,
				Unit:     "request",
			},
		}

		err := s.uc.ModifyQuota(ctx, input)
		s.Error(err)

		var notFoundErr *custom_err.ErrResourceNotFound
		s.True(errors.As(err, &notFoundErr), "error should be of type ErrResourceNotFound")
	})

	s.Run("it should not be able to modify if datasource fails", func() {
		ctx := context.Background()

		s.quotaRepo.IsHealthy = false

		input := ModifyQuotaInput{
			ID: "some-id",
			NewState: model.ModifyQuotaInput{
				Tier:     privilege.TIER_PRO,
				Status:   model.QUOTA_STATUS_ENABLED,
				MaxValue: 200,
				Unit:     "request",
			},
		}

		err := s.uc.ModifyQuota(ctx, input)
		s.Error(err)

		var dsErr *custom_err.ErrDatasourceOperationFailed
		s.True(errors.As(err, &dsErr), "error should be of type ErrDatasourceOperationFailed")
	})

	s.Run("it should not be able to modify if nothing has changed", func() {
		ctx := context.Background()

		quota, err := model.NewQuota(model.QuotaInput{
			Tier:     privilege.TIER_PRO,
			Kind:     kind.QUOTA_LIMIT_REQUEST,
			MaxValue: 100,
			Unit:     "request",
		})
		s.NoError(err)

		quota.Status = model.QUOTA_STATUS_ENABLED

		err = s.quotaRepo.Create(ctx, *quota)
		s.NoError(err)

		input := ModifyQuotaInput{
			ID: quota.ID,
			NewState: model.ModifyQuotaInput{
				Tier:     quota.Tier,
				Status:   quota.Status,
				MaxValue: quota.MaxValue,
				Unit:     quota.Unit,
			},
		}

		err = s.uc.ModifyQuota(ctx, input)
		s.Error(err)

		var nothingToChangeErr *custom_err.ErrNothingToChange
		s.True(errors.As(err, &nothingToChangeErr), "error should be of type ErrNothingToChange")
	})

	s.Run("it should not be able to modify if tier is invalid", func() {
		ctx := context.Background()

		quota, err := model.NewQuota(model.QuotaInput{
			Tier:     privilege.TIER_ROOKIE,
			Kind:     kind.QUOTA_LIMIT_REQUEST,
			MaxValue: 100,
			Unit:     "request",
		})
		s.NoError(err)

		err = s.quotaRepo.Create(ctx, *quota)
		s.NoError(err)

		input := ModifyQuotaInput{
			ID: quota.ID,
			NewState: model.ModifyQuotaInput{
				Tier:     privilege.TIER_ROOKIE + "-invalid",
				Status:   model.QUOTA_STATUS_ENABLED,
				MaxValue: 200,
				Unit:     "request",
			},
		}

		err = s.uc.ModifyQuota(ctx, input)
		s.Error(err)

		var invalidEntityErr *custom_err.ErrInvalidEntity
		s.True(errors.As(err, &invalidEntityErr), "error should be of type ErrInvalidEntity")
	})
}
