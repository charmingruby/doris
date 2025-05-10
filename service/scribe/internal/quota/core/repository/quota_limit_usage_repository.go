package repository

import (
	"context"

	"github.com/charmingruby/doris/service/scribe/internal/quota/core/model"
)

type QuotaLimitUsageRepository interface {
	FindByQuotaLimitIDAndCorrelationID(ctx context.Context, quotaLimitID, correlationID string) (model.QuotaLimitUsage, error)
	Create(ctx context.Context, quotaLimit model.QuotaLimitUsage) error
}
