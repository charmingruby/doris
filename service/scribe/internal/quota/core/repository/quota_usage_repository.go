package repository

import (
	"context"

	"github.com/charmingruby/doris/service/scribe/internal/quota/core/model"
)

type QuotaUsageRepository interface {
	Create(ctx context.Context, quotaUsage model.QuotaUsage) error
}
