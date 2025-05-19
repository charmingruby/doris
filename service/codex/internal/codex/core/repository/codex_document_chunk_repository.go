package repository

import (
	"context"

	"github.com/charmingruby/doris/service/codex/internal/codex/core/model"
)

type CodexDocumentChunkRepository interface {
	Create(ctx context.Context, chunk model.CodexDocumentChunk) error
}
