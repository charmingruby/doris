package repository

import (
	"context"

	"github.com/charmingruby/doris/service/codex/internal/codex/core/model"
)

type CodexDocumentChunkRepository interface {
	Create(ctx context.Context, chunk model.CodexDocumentChunk) error
	FindSimilarChunks(ctx context.Context, codexDocumentID string, embedding []float64, limit int) ([]model.CodexDocumentChunk, error)
}
