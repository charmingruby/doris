package usecase

import (
	"context"
	"errors"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/core/privilege"
	"github.com/charmingruby/doris/service/scribe/internal/quota/core/model"
)

func (s *Suite) Test_ModifyQuota() {
	s.Run("it should modify quota successfully if quota exists", func() {
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

	s.Run("it should return error if quota not found", func() {
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
		s.True(errors.As(err, &notFoundErr))
	})

	s.Run("it should return error if quotaRepo is unhealthy on find", func() {
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
		s.True(errors.As(err, &dsErr))
	})

	s.Run("it should return error if quotaRepo is unhealthy on save", func() {
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

		s.quotaRepo.IsHealthy = false

		err = s.uc.ModifyQuota(ctx, input)
		s.Error(err)

		var dsErr *custom_err.ErrDatasourceOperationFailed
		s.True(errors.As(err, &dsErr))
	})

	s.Run("it should return error if Modify validation fails", func() {
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
				Tier:   "INVALID_TIER",
				Status: model.QUOTA_STATUS_ENABLED,
			},
		}

		err = s.uc.ModifyQuota(ctx, input)
		s.Error(err)
	})
}
