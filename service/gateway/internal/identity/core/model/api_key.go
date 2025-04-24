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

	API_KEY_PLAN_TRIAL = "TRIAL"
	API_KEY_PLAN_PRO   = "PRO"
)

type APIKeyInput struct {
	FirstName               string    `json:"first_name"`
	LastName                string    `json:"last_name"`
	Email                   string    `json:"email"`
	Key                     string    `json:"key"`
	ActivationCode          string    `json:"activation_code"`
	ActivationCodeExpiresAt time.Time `json:"activation_code_expires_at"`
}

func NewAPIKey(in APIKeyInput) *APIKey {
	return &APIKey{
		ID:                      id.New(),
		FirstName:               in.FirstName,
		LastName:                in.LastName,
		Email:                   in.Email,
		Key:                     in.Key,
		Plan:                    API_KEY_PLAN_TRIAL,
		Status:                  API_KEY_STATUS_NONE,
		ActivationCode:          in.ActivationCode,
		ActivationCodeExpiresAt: in.ActivationCodeExpiresAt,
		CreatedAt:               time.Now(),
		UpdatedAt:               nil,
	}
}

type APIKey struct {
	ID                      string     `json:"id"`
	FirstName               string     `json:"first_name"`
	LastName                string     `json:"last_name"`
	Email                   string     `json:"email"`
	Key                     string     `json:"key"`
	Plan                    string     `json:"plan"`
	Status                  string     `json:"status"`
	ActivationCode          string     `json:"activation_code"`
	ActivationCodeExpiresAt time.Time  `json:"activation_code_expires_at"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               *time.Time `json:"updated_at"`
}
