package repository

import (
	"context"

	"github.com/charmingruby/doris/service/scribe/internal/quota/core/model"
)

type QuotaRepository interface {
	FindByTier(ctx context.Context, tier string) (model.Quota, error)
	Create(ctx context.Context, quota model.Quota) error
}
