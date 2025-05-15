package repository

import (
	"context"

	"github.com/charmingruby/doris/service/codex/internal/codex/core/model"
)

type CodexDocumentRepository interface {
	Create(ctx context.Context, codexDocument model.CodexDocument) error
}
