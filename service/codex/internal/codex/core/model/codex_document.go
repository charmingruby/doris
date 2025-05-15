package model

import (
	"time"

	"github.com/charmingruby/doris/lib/core/id"
)

const (
	CodexDocumentStatusPending    = "PENDING"
	CodexDocumentStatusFailed     = "FAILED"
	CodexDocumentStatusProcessing = "PROCESSING"
	CodexDocumentStatusReady      = "READY"
)

type CodexDocumentInput struct {
	CodexID  string
	Title    string
	ImageURL string
}

func NewCodexDocument(in CodexDocumentInput) *CodexDocument {
	return &CodexDocument{
		ID:        id.New(),
		CodexID:   in.CodexID,
		Title:     in.Title,
		ImageURL:  in.ImageURL,
		Status:    CodexDocumentStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: nil,
	}
}

type CodexDocument struct {
	ID        string     `json:"id" db:"id"`
	CodexID   string     `json:"codex_id" db:"codex_id"`
	Title     string     `json:"title" db:"title"`
	ImageURL  string     `json:"image_url" db:"image_url"`
	Status    string     `json:"status" db:"status"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}
