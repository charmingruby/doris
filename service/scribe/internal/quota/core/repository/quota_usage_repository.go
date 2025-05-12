package repository

import (
	"context"

	"github.com/charmingruby/doris/service/scribe/internal/quota/core/model"
)

type QuotaUsageRepository interface {
	FindByCorrelationIDAndQuotaID(ctx context.Context, correlationID, quotaID string) (model.QuotaUsage, error)
	Create(ctx context.Context, quotaUsage model.QuotaUsage) error
	Save(ctx context.Context, quotaUsage model.QuotaUsage) error
}
