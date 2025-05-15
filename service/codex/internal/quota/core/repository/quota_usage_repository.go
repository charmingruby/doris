package repository

import (
	"context"
	"time"

	"github.com/charmingruby/doris/service/codex/internal/quota/core/model"
)

type QuotaUsageRepository interface {
	FindByCorrelationIDAndQuotaID(ctx context.Context, correlationID, quotaID string) (model.QuotaUsage, error)
	Create(ctx context.Context, quotaUsage model.QuotaUsage) error
	UpdateAllCurrentUsages(ctx context.Context, now time.Time) error
	Save(ctx context.Context, quotaUsage model.QuotaUsage) error
}
