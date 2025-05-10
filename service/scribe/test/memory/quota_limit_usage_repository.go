package memory

import (
	"context"

	"github.com/charmingruby/doris/service/scribe/internal/quota/core/model"
)

type QuotaLimitUsageRepository struct {
	Items     []model.QuotaLimitUsage
	IsHealthy bool
}

func NewQuotaLimitUsageRepository() *QuotaLimitUsageRepository {
	return &QuotaLimitUsageRepository{
		Items:     []model.QuotaLimitUsage{},
		IsHealthy: true,
	}
}

func (r *QuotaLimitUsageRepository) FindByQuotaLimitIDAndCorrelationID(ctx context.Context, quotaLimitID, correlationID string) (model.QuotaLimitUsage, error) {
	if !r.IsHealthy {
		return model.QuotaLimitUsage{}, ErrUnhealthyDatasource
	}

	for _, i := range r.Items {
		if i.QuotaLimitID == quotaLimitID && i.CorrelationID == correlationID {
			return i, nil
		}
	}

	return model.QuotaLimitUsage{}, nil
}

func (r *QuotaLimitUsageRepository) Create(ctx context.Context, QuotaLimitUsage model.QuotaLimitUsage) error {
	if !r.IsHealthy {
		return ErrUnhealthyDatasource
	}

	r.Items = append(r.Items, QuotaLimitUsage)

	return nil
}
