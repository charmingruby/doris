package repository

import (
	"context"

	"github.com/charmingruby/doris/service/codex/internal/quota/core/model"
)

type QuotaSnapshotRepository interface {
	FindByCorrelationIDAndKind(ctx context.Context, correlationID, kind string) (model.QuotaSnapshot, error)
	UpdateCurrentUsage(ctx context.Context, correlationID, kind string, usage int) error
}
