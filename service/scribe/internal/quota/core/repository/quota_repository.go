package repository

import (
	"context"

	"github.com/charmingruby/doris/service/scribe/internal/quota/core/model"
)

type QuotaRepository interface {
	FindByID(ctx context.Context, id string) (model.Quota, error)
	FindManyByTier(ctx context.Context, tier string) ([]model.Quota, error)
	FindByTierAndKind(ctx context.Context, tier, kind string) (model.Quota, error)
	Create(ctx context.Context, quota model.Quota) error
	Save(ctx context.Context, quota model.Quota) error
}
