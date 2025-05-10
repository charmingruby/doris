package usecase

import (
	"context"
	"errors"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/core/privilege"
	"github.com/charmingruby/doris/service/scribe/internal/quota/core/model"
)

func (s *Suite) Test_CreateQuotaLimitUsage() {
	validQuota, err := model.NewQuota(model.QuotaInput{
		Tier: privilege.TIER_ROOKIE,
	})
	s.NoError(err)

	validQuotaLimit, err := model.NewQuotaLimit(model.QuotaLimitInput{
		QuotaID:  validQuota.ID,
		Kind:     model.QUOTA_LIMIT_KIND_DOCUMENT,
		MaxValue: 10,
		Unit:     "document",
	})
	s.NoError(err)

	validInput := CreateQuotaLimitUsageInput{
		QuotaLimitID:  validQuotaLimit.ID,
		CorrelationID: "correlation-id-1",
	}

	s.Run("it should be able to create a quota limit usage", func() {
		ctx := context.Background()

		err := s.quotaRepo.Create(ctx, *validQuota)
		s.NoError(err)

		err = s.quotaLimitRepo.Create(ctx, *validQuotaLimit)
		s.NoError(err)

		id, err := s.uc.CreateQuotaLimitUsage(ctx, validInput)
		s.NoError(err)

		storedQuotaLimitUsage := s.quotaLimitUsageRepo.Items[0]

		s.Equal(storedQuotaLimitUsage.ID, id)
		s.Equal(storedQuotaLimitUsage.QuotaLimitID, validInput.QuotaLimitID)
		s.Equal(storedQuotaLimitUsage.CorrelationID, validInput.CorrelationID)
	})

	s.Run("it should not be able to create if quota limit datasource fails", func() {
		ctx := context.Background()
		s.quotaLimitRepo.IsHealthy = false

		id, err := s.uc.CreateQuotaLimitUsage(ctx, validInput)
		s.Empty(id)
		s.Error(err)

		var dsErr *custom_err.ErrDatasourceOperationFailed
		s.True(errors.As(err, &dsErr), "error should be of type ErrDatasourceOperationFailed")
	})

	s.Run("it should not be able to create if quota limit usage datasource fails", func() {
		ctx := context.Background()

		err := s.quotaRepo.Create(ctx, *validQuota)
		s.NoError(err)

		err = s.quotaLimitRepo.Create(ctx, *validQuotaLimit)
		s.NoError(err)

		s.quotaLimitUsageRepo.IsHealthy = false

		id, err := s.uc.CreateQuotaLimitUsage(ctx, validInput)
		s.Empty(id)
		s.Error(err)

		var dsErr *custom_err.ErrDatasourceOperationFailed
		s.True(errors.As(err, &dsErr), "error should be of type ErrDatasourceOperationFailed")
	})

	s.Run("it should not be able to create if quota limit does not exist", func() {
		ctx := context.Background()

		id, err := s.uc.CreateQuotaLimitUsage(ctx, validInput)
		s.Empty(id)
		s.Error(err)

		var notFoundErr *custom_err.ErrResourceNotFound
		s.True(errors.As(err, &notFoundErr), "error should be of type ErrResourceNotFound")
	})

	s.Run("it should not be able to create if usage already exists", func() {
		ctx := context.Background()

		err := s.quotaRepo.Create(ctx, *validQuota)
		s.NoError(err)

		err = s.quotaLimitRepo.Create(ctx, *validQuotaLimit)
		s.NoError(err)

		usage := model.NewQuotaLimitUsage(model.QuotaLimitUsageInput{
			CorrelationID: validInput.CorrelationID,
			QuotaLimitID:  validInput.QuotaLimitID,
		})

		err = s.quotaLimitUsageRepo.Create(ctx, *usage)
		s.NoError(err)

		id, err := s.uc.CreateQuotaLimitUsage(ctx, validInput)
		s.Empty(id)
		s.Error(err)

		var alreadyExistsErr *custom_err.ErrResourceAlreadyExists
		s.True(errors.As(err, &alreadyExistsErr), "error should be of type ErrResourceAlreadyExists")
	})
}
