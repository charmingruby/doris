package usecase

import (
	"context"
	"errors"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/core/privilege"
	"github.com/charmingruby/doris/service/scribe/internal/quota/core/model"
)

func (s *Suite) Test_CreateQuota() {
	s.Run("it should be able to create a quota", func() {
		tier := privilege.TIER_ROOKIE

		id, err := s.uc.CreateQuota(context.Background(), CreateQuotaInput{
			Tier: tier,
		})
		s.NoError(err)

		storedQuota := s.quotaRepo.Items[0]

		s.Equal(storedQuota.ID, id)
		s.Equal(storedQuota.Status, model.QUOTA_STATUS_DRAFT)
		s.Equal(storedQuota.Tier, tier)
	})

	s.Run("it should be not able to create a quota if datasource fails", func() {
		s.quotaRepo.IsHealthy = false

		id, err := s.uc.CreateQuota(context.Background(), CreateQuotaInput{
			Tier: privilege.TIER_ROOKIE,
		})
		s.Empty(id)
		s.Error(err)

		var dsErr *custom_err.ErrDatasourceOperationFailed
		s.True(errors.As(err, &dsErr), "error should be of type ErrDatasourceOperationFailed")
	})

	s.Run("it should be not able to create a quota if tier quota already exists", func() {
		ctx := context.Background()

		tier := privilege.TIER_ROOKIE

		quota, err := model.NewQuota(&model.QuotaInput{
			Tier: tier,
		})
		s.NoError(err)

		err = s.quotaRepo.Create(ctx, *quota)
		s.NoError(err)

		id, err := s.uc.CreateQuota(ctx, CreateQuotaInput{
			Tier: tier,
		})
		s.Empty(id)
		s.Error(err)

		var resourceAlreadyExistsErr *custom_err.ErrResourceAlreadyExists
		s.True(errors.As(err, &resourceAlreadyExistsErr), "error should be of type ErrResourceAlreadyExists")
	})

	s.Run("it should be not able to create a quota if tier is invalid", func() {
		id, err := s.uc.CreateQuota(context.Background(), CreateQuotaInput{
			Tier: privilege.TIER_ROOKIE + "-invalid",
		})
		s.Empty(id)
		s.Error(err)

		var invalidEntityErr *custom_err.ErrInvalidEntity
		s.True(errors.As(err, &invalidEntityErr), "error should be of type ErrInvalidEntity")
	})
}
