package repository

import (
	"context"

	"github.com/charmingruby/doris/service/codex/internal/codex/core/model"
)

type CodexDocumentRepository interface {
	FindByID(ctx context.Context, id string) (model.CodexDocument, error)
	Create(ctx context.Context, codexDocument model.CodexDocument) error
}
