package usecase

import (
	"context"
	"errors"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/core/privilege"
	"github.com/charmingruby/doris/service/scribe/internal/quota/core/model"
)

func (s *Suite) Test_CreateQuotaLimit() {
	validQuota, err := model.NewQuota(model.QuotaInput{
		Tier: privilege.TIER_ROOKIE,
	})
	s.NoError(err)

	validInput := CreateQuotaLimtInput{
		QuotaID:  validQuota.ID,
		Kind:     model.QUOTA_LIMIT_KIND_DOCUMENT,
		MaxValue: 10,
		Unit:     "document",
	}

	s.Run("it should be able to create a quota limit", func() {
		ctx := context.Background()

		err := s.quotaRepo.Create(ctx, *validQuota)
		s.NoError(err)

		id, err := s.uc.CreateQuotaLimit(ctx, validInput)
		s.NoError(err)

		storedQuotaLimit := s.quotaLimitRepo.Items[0]

		s.Equal(storedQuotaLimit.ID, id)
		s.Equal(storedQuotaLimit.QuotaID, validInput.QuotaID)
		s.Equal(storedQuotaLimit.Kind, validInput.Kind)
		s.Equal(storedQuotaLimit.MaxValue, validInput.MaxValue)
		s.Equal(storedQuotaLimit.Unit, validInput.Unit)
	})

	s.Run("it should be not able to create a quota limit if quota datasource fails", func() {
		ctx := context.Background()

		err := s.quotaRepo.Create(ctx, *validQuota)
		s.NoError(err)

		s.quotaRepo.IsHealthy = false

		id, err := s.uc.CreateQuotaLimit(ctx, validInput)
		s.Empty(id)
		s.Error(err)

		var dsErr *custom_err.ErrDatasourceOperationFailed
		s.True(errors.As(err, &dsErr), "error should be of type ErrDatasourceOperationFailed")
	})

	s.Run("it should be not able to create a quota limit if quota limit datasource fails", func() {
		s.quotaLimitRepo.IsHealthy = false

		ctx := context.Background()

		err := s.quotaRepo.Create(ctx, *validQuota)
		s.NoError(err)

		id, err := s.uc.CreateQuotaLimit(ctx, validInput)
		s.Empty(id)
		s.Error(err)

		var dsErr *custom_err.ErrDatasourceOperationFailed
		s.True(errors.As(err, &dsErr), "error should be of type ErrDatasourceOperationFailed")
	})

	s.Run("it should be not able to create a quota limit if quota does not exists", func() {
		ctx := context.Background()

		id, err := s.uc.CreateQuotaLimit(ctx, validInput)
		s.Empty(id)
		s.Error(err)

		var resourceNotFoundErr *custom_err.ErrResourceNotFound
		s.True(errors.As(err, &resourceNotFoundErr), "error should be of type ErrResourceNotFound")
	})

	s.Run("it should be not able to create a quota limit if already exists", func() {
		ctx := context.Background()

		err := s.quotaRepo.Create(ctx, *validQuota)
		s.NoError(err)

		quotaLimit, err := model.NewQuotaLimit(model.QuotaLimitInput{
			QuotaID:  validInput.QuotaID,
			Kind:     validInput.Kind,
			MaxValue: validInput.MaxValue,
			Unit:     validInput.Unit,
		})
		s.NoError(err)

		err = s.quotaLimitRepo.Create(ctx, *quotaLimit)
		s.NoError(err)

		id, err := s.uc.CreateQuotaLimit(ctx, validInput)
		s.Empty(id)
		s.Error(err)

		var resourceAlreadyExistsErr *custom_err.ErrResourceAlreadyExists
		s.True(errors.As(err, &resourceAlreadyExistsErr), "error should be of type ErrResourceAlreadyExists")
	})

	s.Run("it should be not able to create a quota limit if kind is invalid", func() {
		ctx := context.Background()

		err := s.quotaRepo.Create(ctx, *validQuota)
		s.NoError(err)

		invalidInput := validInput
		invalidInput.Kind = model.QUOTA_LIMIT_KIND_DOCUMENT + "-invalid"

		id, err := s.uc.CreateQuotaLimit(ctx, invalidInput)
		s.Empty(id)
		s.Error(err)

		var invalidEntityErr *custom_err.ErrInvalidEntity
		s.True(errors.As(err, &invalidEntityErr), "error should be of type ErrInvalidEntity")
	})
}
