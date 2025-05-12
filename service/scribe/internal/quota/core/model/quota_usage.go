package model

import (
	"time"

	"github.com/charmingruby/doris/lib/core/id"
)

type QuotaUsageInput struct {
	CorrelationID string
	QuotaID       string
}

func NewQuotaUsage(in QuotaUsageInput) *QuotaUsage {
	return &QuotaUsage{
		ID:            id.New(),
		CorrelationID: in.CorrelationID,
		QuotaID:       in.QuotaID,
		CurrentUsage:  0,
		CreatedAt:     time.Now(),
		LastResetAt:   nil,
		UpdatedAt:     nil,
	}
}

type QuotaUsage struct {
	ID            string     `json:"id" db:"id"`
	CorrelationID string     `json:"correlation_id" db:"correlation_id"`
	QuotaID       string     `json:"quota_id" db:"quota_id"`
	CurrentUsage  int        `json:"current_usage" db:"current_usage"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	LastResetAt   *time.Time `json:"last_reset_at" db:"last_reset_at"`
	UpdatedAt     *time.Time `json:"updated_at" db:"updated_at"`
}
