package usecase

import (
	"context"
	"errors"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/core/privilege"
	"github.com/charmingruby/doris/service/codex/internal/quota/core/model"
	"github.com/charmingruby/doris/service/codex/internal/shared/core/kind"
)

func (s *Suite) Test_RecalculateQuotaUsageOnTierChange() {
	s.Run("it should be able to recalculate quota usages when changing tiers", func() {
		ctx := context.Background()
		correlationID := "test-correlation-id"
		oldTier := privilege.TIER_ROOKIE
		newTier := privilege.TIER_PRO

		oldQuota1, err := model.NewQuota(model.QuotaInput{
			Tier:     oldTier,
			Kind:     kind.QUOTA_LIMIT_PROMPT,
			MaxValue: 100,
			Unit:     "request",
		})
		s.NoError(err)

		oldQuota2, err := model.NewQuota(model.QuotaInput{
			Tier:     oldTier,
			Kind:     kind.QUOTA_LIMIT_DOCUMENT,
			MaxValue: 50,
			Unit:     "document",
		})
		s.NoError(err)

		err = s.quotaRepo.Create(ctx, *oldQuota1)
		s.NoError(err)
		err = s.quotaRepo.Create(ctx, *oldQuota2)
		s.NoError(err)

		newQuota1, err := model.NewQuota(model.QuotaInput{
			Tier:     newTier,
			Kind:     kind.QUOTA_LIMIT_PROMPT,
			MaxValue: 200,
			Unit:     "request",
		})
		s.NoError(err)

		newQuota2, err := model.NewQuota(model.QuotaInput{
			Tier:     newTier,
			Kind:     kind.QUOTA_LIMIT_DOCUMENT,
			MaxValue: 100,
			Unit:     "document",
		})
		s.NoError(err)

		err = s.quotaRepo.Create(ctx, *newQuota1)
		s.NoError(err)
		err = s.quotaRepo.Create(ctx, *newQuota2)
		s.NoError(err)

		oldUsage1 := model.NewQuotaUsage(model.QuotaUsageInput{
			CorrelationID: correlationID,
			QuotaID:       oldQuota1.ID,
		})
		oldUsage1.CurrentUsage = 50

		oldUsage2 := model.NewQuotaUsage(model.QuotaUsageInput{
			CorrelationID: correlationID,
			QuotaID:       oldQuota2.ID,
		})
		oldUsage2.CurrentUsage = 25

		err = s.quotaUsageRepo.Create(ctx, *oldUsage1)
		s.NoError(err)
		err = s.quotaUsageRepo.Create(ctx, *oldUsage2)
		s.NoError(err)

		err = s.uc.RecalculateQuotaUsageOnTierChange(ctx, RecalculateQuotaUsageOnTierChangeInput{
			CorrelationID: correlationID,
			OldTier:       oldTier,
			NewTier:       newTier,
		})
		s.NoError(err)

		oldUsage1After, err := s.quotaUsageRepo.FindByCorrelationIDAndQuotaID(ctx, correlationID, oldQuota1.ID)
		s.NoError(err)
		s.False(oldUsage1After.IsActive)
		s.Equal(50, oldUsage1After.CurrentUsage)

		oldUsage2After, err := s.quotaUsageRepo.FindByCorrelationIDAndQuotaID(ctx, correlationID, oldQuota2.ID)
		s.NoError(err)
		s.False(oldUsage2After.IsActive)
		s.Equal(25, oldUsage2After.CurrentUsage)

		newUsage1, err := s.quotaUsageRepo.FindByCorrelationIDAndQuotaID(ctx, correlationID, newQuota1.ID)
		s.NoError(err)
		s.True(newUsage1.IsActive)
		s.Equal(0, newUsage1.CurrentUsage)

		newUsage2, err := s.quotaUsageRepo.FindByCorrelationIDAndQuotaID(ctx, correlationID, newQuota2.ID)
		s.NoError(err)
		s.True(newUsage2.IsActive)
		s.Equal(0, newUsage2.CurrentUsage)
	})

	s.Run("it should handle unmatched quotas correctly", func() {
		ctx := context.Background()
		correlationID := "test-correlation-id"
		oldTier := privilege.TIER_ROOKIE
		newTier := privilege.TIER_PRO

		oldQuota, err := model.NewQuota(model.QuotaInput{
			Tier:     oldTier,
			Kind:     kind.QUOTA_LIMIT_PROMPT,
			MaxValue: 100,
			Unit:     "request",
		})
		s.NoError(err)

		err = s.quotaRepo.Create(ctx, *oldQuota)
		s.NoError(err)

		newQuota, err := model.NewQuota(model.QuotaInput{
			Tier:     newTier,
			Kind:     kind.QUOTA_LIMIT_DOCUMENT,
			MaxValue: 50,
			Unit:     "document",
		})
		s.NoError(err)

		err = s.quotaRepo.Create(ctx, *newQuota)
		s.NoError(err)

		oldUsage := model.NewQuotaUsage(model.QuotaUsageInput{
			CorrelationID: correlationID,
			QuotaID:       oldQuota.ID,
		})
		oldUsage.CurrentUsage = 50

		err = s.quotaUsageRepo.Create(ctx, *oldUsage)
		s.NoError(err)

		err = s.uc.RecalculateQuotaUsageOnTierChange(ctx, RecalculateQuotaUsageOnTierChangeInput{
			CorrelationID: correlationID,
			OldTier:       oldTier,
			NewTier:       newTier,
		})
		s.NoError(err)

		oldUsageAfter, err := s.quotaUsageRepo.FindByCorrelationIDAndQuotaID(ctx, correlationID, oldQuota.ID)
		s.NoError(err)
		s.False(oldUsageAfter.IsActive)
		s.Equal(50, oldUsageAfter.CurrentUsage)

		newUsage, err := s.quotaUsageRepo.FindByCorrelationIDAndQuotaID(ctx, correlationID, newQuota.ID)
		s.NoError(err)
		s.True(newUsage.IsActive)
		s.Equal(0, newUsage.CurrentUsage)
	})

	s.Run("it should be not able to recalculate if quota datasource fails", func() {
		ctx := context.Background()
		correlationID := "test-correlation-id"
		oldTier := privilege.TIER_ROOKIE
		newTier := privilege.TIER_PRO

		s.quotaRepo.IsHealthy = false

		err := s.uc.RecalculateQuotaUsageOnTierChange(ctx, RecalculateQuotaUsageOnTierChangeInput{
			CorrelationID: correlationID,
			OldTier:       oldTier,
			NewTier:       newTier,
		})
		s.Error(err)

		var dsErr *custom_err.ErrDatasourceOperationFailed
		s.True(errors.As(err, &dsErr), "error should be of type ErrDatasourceOperationFailed")
	})

	s.Run("it should be able to recalculate empty tiers correctly", func() {
		ctx := context.Background()
		correlationID := "test-correlation-id"
		oldTier := privilege.TIER_ROOKIE
		newTier := privilege.TIER_PRO

		err := s.uc.RecalculateQuotaUsageOnTierChange(ctx, RecalculateQuotaUsageOnTierChangeInput{
			CorrelationID: correlationID,
			OldTier:       oldTier,
			NewTier:       newTier,
		})
		s.NoError(err)
		s.Equal(0, len(s.quotaUsageRepo.Items))
	})
}
