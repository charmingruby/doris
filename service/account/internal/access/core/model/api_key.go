package model

import (
	"time"

	"github.com/charmingruby/doris/lib/core/id"
	"github.com/charmingruby/doris/lib/core/privilege"
)

const (
	API_KEY_STATUS_PENDING   = "PENDING"
	API_KEY_STATUS_ACTIVE    = "ACTIVE"
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
		Status:    API_KEY_STATUS_PENDING,
		Tier:      privilege.API_KEY_TIER_ROOKIE,
		CreatedAt: time.Now(),
	}
}

type APIKey struct {
	ID        string    `json:"id" db:"id"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	Email     string    `json:"email" db:"email"`
	Key       string    `json:"key" db:"key"`
	Tier      string    `json:"tier" db:"tier"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (a *APIKey) Validate() error {
	return privilege.ValidateAPIKeyTier(a.Tier)
}
