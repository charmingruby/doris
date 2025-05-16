package repository

import (
	"context"

	"github.com/charmingruby/doris/service/codex/internal/codex/core/model"
)

type CodexRepository interface {
	FindByIDAndCorrelationID(ctx context.Context, id, correlationID string) (model.Codex, error)
	FindByCorrelationIDAndName(ctx context.Context, correlationID, name string) (model.Codex, error)
	Create(ctx context.Context, codex model.Codex) error
}
