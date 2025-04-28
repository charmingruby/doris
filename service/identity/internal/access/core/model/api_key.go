package model

import (
	"time"

	"github.com/charmingruby/doris/lib/core/id"
)

const (
	API_KEY_STATUS_NONE      = "NONE"
	API_KEY_STATUS_PENDING   = "PENDING"
	API_KEY_STATUS_ACTIVE    = "ACTIVE"
	API_KEY_STATUS_EXPIRED   = "EXPIRED"
	API_KEY_STATUS_DEFAULTER = "DEFAULTER"
	API_KEY_STATUS_INACTIVE  = "INACTIVE"
)

type APIKeyInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Key       string `json:"key"`
}

func NewAPIKey(in APIKeyInput) *APIKey {
	return &APIKey{
		ID:        id.New(),
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Email:     in.Email,
		Key:       in.Key,
		Status:    API_KEY_STATUS_NONE,
		CreatedAt: time.Now(),
		UpdatedAt: nil,
	}
}

type APIKey struct {
	ID        string     `json:"id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Email     string     `json:"email"`
	Key       string     `json:"key"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
