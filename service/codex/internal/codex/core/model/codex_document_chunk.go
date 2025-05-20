package model

import (
	"time"

	"github.com/charmingruby/doris/lib/core/id"
)

type CodexDocumentChunkInput struct {
	CodexDocumentID string
	Embedding       []float64
	Content         string
}

func NewCodexDocumentChunk(in CodexDocumentChunkInput) *CodexDocumentChunk {
	return &CodexDocumentChunk{
		ID:              id.New(),
		CodexDocumentID: in.CodexDocumentID,
		Embedding:       in.Embedding,
		Content:         in.Content,
		CreatedAt:       time.Now(),
	}
}

type CodexDocumentChunk struct {
	ID              string    `json:"id" db:"id"`
	CodexDocumentID string    `json:"codex_document_id" db:"codex_document_id"`
	Embedding       []float64 `json:"embedding" db:"embedding"`
	Content         string    `json:"content" db:"content"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}
