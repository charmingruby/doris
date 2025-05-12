package memory

import (
	"context"

	"github.com/charmingruby/doris/service/scribe/internal/quota/core/model"
)

type QuotaUsageRepository struct {
	Items     []model.QuotaUsage
	IsHealthy bool
}

func NewQuotaUsageRepository() *QuotaUsageRepository {
	return &QuotaUsageRepository{
		Items:     []model.QuotaUsage{},
		IsHealthy: true,
	}
}

func (r *QuotaUsageRepository) Create(ctx context.Context, quotaUsage model.QuotaUsage) error {
	if !r.IsHealthy {
		return ErrUnhealthyDatasource
	}

	r.Items = append(r.Items, quotaUsage)

	return nil
}
