package repository

import (
	"context"

	"github.com/charmingruby/doris/service/scribe/internal/quota/core/model"
)

type QuotaLimitRepository interface {
	FindByQuotaIDAndKind(ctx context.Context, quotaID, kind string) (model.QuotaLimit, error)
	Create(ctx context.Context, quotaLimit model.QuotaLimit) error
}
