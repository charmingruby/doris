package provider

import (
	"context"

	"github.com/charmingruby/doris/service/scribe/internal/quota/core/repository"
)

type QuotaUsageManagmentProvider struct {
	quotaSnapshotRepository repository.QuotaSnapshotRepository
}

func NewQuotaUsageManagmentProvider(quotaSnapshotRepository repository.QuotaSnapshotRepository) *QuotaUsageManagmentProvider {
	return &QuotaUsageManagmentProvider{
		quotaSnapshotRepository: quotaSnapshotRepository,
	}
}

func (p *QuotaUsageManagmentProvider) CheckQuotaAvailability(ctx context.Context, correlationID, kind string, usage int) (bool, error) {
	snapshot, err := p.quotaSnapshotRepository.FindByCorrelationIDAndKind(ctx, correlationID, kind)
	if err != nil {
		return false, err
	}

	if snapshot.CurrentUsage+usage > snapshot.MaxValue {
		return false, nil
	}

	return true, nil
}

func (p *QuotaUsageManagmentProvider) ConsumeQuota(ctx context.Context, correlationID, kind string, usage int) error {
	return p.quotaSnapshotRepository.UpdateCurrentUsage(ctx, correlationID, kind, usage)
}
