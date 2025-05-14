package memory

import (
	"context"
	"time"

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

func (r *QuotaUsageRepository) FindByCorrelationIDAndQuotaID(ctx context.Context, correlationID, quotaID string) (model.QuotaUsage, error) {
	if !r.IsHealthy {
		return model.QuotaUsage{}, ErrUnhealthyDatasource
	}

	for _, i := range r.Items {
		if i.CorrelationID == correlationID && i.QuotaID == quotaID {
			return i, nil
		}
	}

	return model.QuotaUsage{}, nil
}

func (r *QuotaUsageRepository) Create(ctx context.Context, quotaUsage model.QuotaUsage) error {
	if !r.IsHealthy {
		return ErrUnhealthyDatasource
	}

	r.Items = append(r.Items, quotaUsage)

	return nil
}

func (r *QuotaUsageRepository) Save(ctx context.Context, quotaUsage model.QuotaUsage) error {
	if !r.IsHealthy {
		return ErrUnhealthyDatasource
	}

	for idx, i := range r.Items {
		if i.ID == quotaUsage.ID {
			r.Items[idx] = quotaUsage
		}
	}

	return nil
}

func (r *QuotaUsageRepository) UpdateAllCurrentUsages(ctx context.Context, now time.Time) error {
	if !r.IsHealthy {
		return ErrUnhealthyDatasource
	}

	for idx := range r.Items {
		r.Items[idx].CurrentUsage = 0
		r.Items[idx].LastResetAt = &now
		r.Items[idx].UpdatedAt = &now
	}

	return nil
}
