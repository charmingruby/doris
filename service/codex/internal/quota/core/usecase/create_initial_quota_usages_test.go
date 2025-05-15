package usecase

import (
	"context"
	"errors"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/core/privilege"
	"github.com/charmingruby/doris/service/codex/internal/quota/core/model"
	"github.com/charmingruby/doris/service/codex/internal/shared/core/kind"
)

func (s *Suite) Test_CreateInitialQuotaUsages() {
	s.Run("it should be able to create initial quota usages", func() {
		ctx := context.Background()
		tier := privilege.TIER_ROOKIE
		correlationID := "test-correlation-id"

		quota, err := model.NewQuota(model.QuotaInput{
			Tier:     tier,
			Kind:     kind.QUOTA_LIMIT_REQUEST,
			MaxValue: 100,
			Unit:     "request",
		})
		s.NoError(err)

		err = s.quotaRepo.Create(ctx, *quota)
		s.NoError(err)

		err = s.uc.CreateInitialQuotaUsages(ctx, CreateInitialQuotaUsagesInput{
			Tier:          tier,
			CorrelationID: correlationID,
		})
		s.NoError(err)

		s.Equal(1, len(s.quotaUsageRepo.Items))

		usage := s.quotaUsageRepo.Items[0]
		s.Equal(quota.ID, usage.QuotaID)
		s.Equal(correlationID, usage.CorrelationID)
	})

	s.Run("it should be not able to create initial quotes if there is no quotas for the tier", func() {
		ctx := context.Background()
		tier := privilege.TIER_ROOKIE
		correlationID := "test-correlation-id"

		err := s.uc.CreateInitialQuotaUsages(ctx, CreateInitialQuotaUsagesInput{
			Tier:          tier,
			CorrelationID: correlationID,
		})

		s.NoError(err)
		s.Equal(0, len(s.quotaUsageRepo.Items))
	})

	s.Run("it should be not able to create initial quotas if datasource fails", func() {
		ctx := context.Background()
		tier := privilege.TIER_ROOKIE
		correlationID := "test-correlation-id"

		s.quotaRepo.IsHealthy = false

		err := s.uc.CreateInitialQuotaUsages(ctx, CreateInitialQuotaUsagesInput{
			Tier:          tier,
			CorrelationID: correlationID,
		})
		s.Error(err)

		var dsErr *custom_err.ErrDatasourceOperationFailed
		s.True(errors.As(err, &dsErr), "error should be of type ErrDatasourceOperationFailed")
	})

	s.Run("it should be able to create multiple quota usages if multiple quotas exist", func() {
		ctx := context.Background()
		tier := privilege.TIER_ROOKIE
		correlationID := "test-correlation-id"

		quota1, err := model.NewQuota(model.QuotaInput{
			Tier:     tier,
			Kind:     kind.QUOTA_LIMIT_REQUEST,
			MaxValue: 100,
			Unit:     "request",
		})
		s.NoError(err)

		quota2, err := model.NewQuota(model.QuotaInput{
			Tier:     tier,
			Kind:     kind.QUOTA_LIMIT_DOCUMENT,
			MaxValue: 50,
			Unit:     "document",
		})
		s.NoError(err)

		err = s.quotaRepo.Create(ctx, *quota1)
		s.NoError(err)
		err = s.quotaRepo.Create(ctx, *quota2)
		s.NoError(err)

		err = s.uc.CreateInitialQuotaUsages(ctx, CreateInitialQuotaUsagesInput{
			Tier:          tier,
			CorrelationID: correlationID,
		})
		s.NoError(err)

		s.Equal(2, len(s.quotaUsageRepo.Items))

		s.Equal(quota1.ID, s.quotaUsageRepo.Items[0].QuotaID)
		s.Equal(correlationID, s.quotaUsageRepo.Items[0].CorrelationID)

		s.Equal(quota2.ID, s.quotaUsageRepo.Items[1].QuotaID)
		s.Equal(correlationID, s.quotaUsageRepo.Items[1].CorrelationID)
	})
}
