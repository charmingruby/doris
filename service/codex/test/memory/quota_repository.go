package memory

import (
	"context"

	"github.com/charmingruby/doris/service/scribe/internal/quota/core/model"
)

type QuotaRepository struct {
	Items     []model.Quota
	IsHealthy bool
}

func NewQuotaRepository() *QuotaRepository {
	return &QuotaRepository{
		Items:     []model.Quota{},
		IsHealthy: true,
	}
}

func (r *QuotaRepository) FindByID(ctx context.Context, id string) (model.Quota, error) {
	if !r.IsHealthy {
		return model.Quota{}, ErrUnhealthyDatasource
	}

	for _, i := range r.Items {
		if i.ID == id {
			return i, nil
		}
	}

	return model.Quota{}, nil
}

func (r *QuotaRepository) FindManyByTier(ctx context.Context, tier string) ([]model.Quota, error) {
	if !r.IsHealthy {
		return nil, ErrUnhealthyDatasource
	}

	quotas := []model.Quota{}
	for _, i := range r.Items {
		if i.Tier == tier {
			quotas = append(quotas, i)
		}
	}

	return quotas, nil
}

func (r *QuotaRepository) FindByTierAndKind(ctx context.Context, tier, kind string) (model.Quota, error) {
	if !r.IsHealthy {
		return model.Quota{}, ErrUnhealthyDatasource
	}

	for _, i := range r.Items {
		if i.Tier == tier && i.Kind == kind {
			return i, nil
		}
	}

	return model.Quota{}, nil
}

func (r *QuotaRepository) Create(ctx context.Context, quota model.Quota) error {
	if !r.IsHealthy {
		return ErrUnhealthyDatasource
	}

	r.Items = append(r.Items, quota)

	return nil
}

func (r *QuotaRepository) Save(ctx context.Context, quota model.Quota) error {
	if !r.IsHealthy {
		return ErrUnhealthyDatasource
	}

	for idx, i := range r.Items {
		if i.ID == quota.ID {
			r.Items[idx] = quota
		}
	}

	return nil
}
