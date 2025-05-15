package model

import (
	"time"

	"github.com/charmingruby/doris/lib/core/id"
)

type CodexInput struct {
	Name          string
	Description   string
	CorrelationID string
}

func NewCodex(in CodexInput) *Codex {
	return &Codex{
		ID:            id.New(),
		CorrelationID: in.CorrelationID,
		Name:          in.Name,
		Description:   in.Description,
		CreatedAt:     time.Now(),
		UpdatedAt:     nil,
	}
}

type Codex struct {
	ID            string     `json:"id" db:"id"`
	CorrelationID string     `json:"correlation_id" db:"correlation_id"`
	Name          string     `json:"name" db:"name"`
	Description   string     `json:"description" db:"description"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at" db:"updated_at"`
}
