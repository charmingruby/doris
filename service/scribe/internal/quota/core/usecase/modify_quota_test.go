package usecase

import (
	"context"
	"errors"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/core/privilege"
	"github.com/charmingruby/doris/service/scribe/internal/quota/core/model"
)

func (s *Suite) Test_ModifyQuota() {
	s.Run("it should be able to modify a quota", func() {
		ctx := context.Background()

		quota, err := model.NewQuota(model.QuotaInput{
			Tier: privilege.TIER_ROOKIE,
		})
		s.NoError(err)

		err = s.quotaRepo.Create(ctx, *quota)
		s.NoError(err)

		input := ModifyQuotaInput{
			ID: quota.ID,
			NewState: model.ModifyQuotaInput{
				Tier:   privilege.TIER_PRO,
				Status: model.QUOTA_STATUS_ENABLED,
			},
		}

		err = s.uc.ModifyQuota(ctx, input)
		s.NoError(err)

		storedQuota := s.quotaRepo.Items[0]
		s.Equal(input.NewState.Tier, storedQuota.Tier)
		s.Equal(input.NewState.Status, storedQuota.Status)
	})

	s.Run("it should not be able to modify if quota does not exist", func() {
		ctx := context.Background()

		input := ModifyQuotaInput{
			ID: "non-existent-id",
			NewState: model.ModifyQuotaInput{
				Tier:   privilege.TIER_PRO,
				Status: model.QUOTA_STATUS_ENABLED,
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
				Tier:   privilege.TIER_PRO,
				Status: model.QUOTA_STATUS_ENABLED,
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
			Tier: privilege.TIER_PRO,
		})
		s.NoError(err)

		quota.Status = model.QUOTA_STATUS_ENABLED

		err = s.quotaRepo.Create(ctx, *quota)
		s.NoError(err)

		input := ModifyQuotaInput{
			ID: quota.ID,
			NewState: model.ModifyQuotaInput{
				Tier:   quota.Tier,
				Status: quota.Status,
			},
		}

		err = s.uc.ModifyQuota(ctx, input)
		s.Error(err)

		var nothingToChangeErr *custom_err.ErrNothingToChange
		s.True(errors.As(err, &nothingToChangeErr), "error should be of type ErrNothingToChange")
	})
}
