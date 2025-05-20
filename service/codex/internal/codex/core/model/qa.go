package model

import (
	"time"

	"github.com/charmingruby/doris/lib/core/id"
)

type QA struct {
	ID        string    `json:"id"`
	CodexID   string    `json:"codex_id"`
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	CreatedAt time.Time `json:"created_at"`
}

type QAInput struct {
	CodexID  string `json:"codex_id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func NewQA(in QAInput) *QA {
	return &QA{
		ID:        id.New(),
		CodexID:   in.CodexID,
		Question:  in.Question,
		Answer:    in.Answer,
		CreatedAt: time.Now(),
	}
}
