package model

import (
	"time"

	"github.com/charmingruby/doris/lib/core/id"
)

type CodexDocumentChunkInput struct {
	CodexDocumentID string
	Embedding       []float64
}

func NewCodexDocumentChunk(in CodexDocumentChunkInput) *CodexDocumentChunk {
	return &CodexDocumentChunk{
		ID:              id.New(),
		CodexDocumentID: in.CodexDocumentID,
		Embedding:       in.Embedding,
		CreatedAt:       time.Now(),
	}
}

type CodexDocumentChunk struct {
	ID              string    `json:"id" db:"id"`
	CodexDocumentID string    `json:"codex_document_id" db:"codex_document_id"`
	Embedding       []float64 `json:"embedding" db:"embedding"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}
