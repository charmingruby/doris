package model

import (
	"time"

	"github.com/charmingruby/doris/lib/core/id"
)

type NewCodexInput struct {
	Name        string
	Description string
}

func NewCodex(in NewCodexInput) *Codex {
	return &Codex{
		ID:          id.New(),
		Name:        in.Name,
		Description: in.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   nil,
	}
}

type Codex struct {
	ID          string     `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description" db:"description"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at" db:"updated_at"`
}
