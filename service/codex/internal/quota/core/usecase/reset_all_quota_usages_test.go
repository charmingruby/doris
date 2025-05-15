package usecase

import (
	"context"
	"errors"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/core/privilege"
	"github.com/charmingruby/doris/service/codex/internal/quota/core/model"
	"github.com/charmingruby/doris/service/codex/internal/shared/core/kind"
)

func (s *Suite) Test_ResetAllQuotaUsages() {
	s.Run("it should be able to reset all quota usages", func() {
		ctx := context.Background()
		correlationID := "test-correlation-id"
		tier := privilege.TIER_PRO

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

		usage1 := model.NewQuotaUsage(model.QuotaUsageInput{
			CorrelationID: correlationID,
			QuotaID:       quota1.ID,
		})
		usage1.CurrentUsage = 50

		usage2 := model.NewQuotaUsage(model.QuotaUsageInput{
			CorrelationID: correlationID,
			QuotaID:       quota2.ID,
		})
		usage2.CurrentUsage = 25

		err = s.quotaUsageRepo.Create(ctx, *usage1)
		s.NoError(err)
		err = s.quotaUsageRepo.Create(ctx, *usage2)
		s.NoError(err)

		err = s.uc.ResetAllQuotaUsages(ctx)
		s.NoError(err)

		usage1After, err := s.quotaUsageRepo.FindByCorrelationIDAndQuotaID(ctx, correlationID, quota1.ID)
		s.NoError(err)
		s.Equal(0, usage1After.CurrentUsage)

		usage2After, err := s.quotaUsageRepo.FindByCorrelationIDAndQuotaID(ctx, correlationID, quota2.ID)
		s.NoError(err)
		s.Equal(0, usage2After.CurrentUsage)
	})

	s.Run("it should handle datasource failure correctly", func() {
		ctx := context.Background()
		s.quotaUsageRepo.IsHealthy = false

		err := s.uc.ResetAllQuotaUsages(ctx)
		s.Error(err)

		var dsErr *custom_err.ErrDatasourceOperationFailed
		s.True(errors.As(err, &dsErr), "error should be of type ErrDatasourceOperationFailed")
	})

	s.Run("it should handle empty quota usages correctly", func() {
		ctx := context.Background()

		err := s.uc.ResetAllQuotaUsages(ctx)
		s.NoError(err)
		s.Equal(0, len(s.quotaUsageRepo.Items))
	})
}
