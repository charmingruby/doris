package memory

import (
	"context"

	"github.com/charmingruby/doris/service/scribe/internal/quota/core/model"
)

type QuotaLimitRepository struct {
	Items     []model.QuotaLimit
	IsHealthy bool
}

func NewQuotaLimitRepository() *QuotaLimitRepository {
	return &QuotaLimitRepository{
		Items:     []model.QuotaLimit{},
		IsHealthy: true,
	}
}
func (r *QuotaLimitRepository) FindByQuotaIDAndKind(ctx context.Context, quotaID, kind string) (model.QuotaLimit, error) {
	if !r.IsHealthy {
		return model.QuotaLimit{}, ErrUnhealthyDatasource
	}

	for _, i := range r.Items {
		if i.QuotaID == quotaID && i.Kind == kind {
			return i, nil
		}
	}

	return model.QuotaLimit{}, nil
}

func (r *QuotaLimitRepository) Create(ctx context.Context, quotaLimit model.QuotaLimit) error {
	if !r.IsHealthy {
		return ErrUnhealthyDatasource
	}

	r.Items = append(r.Items, quotaLimit)

	return nil
}
