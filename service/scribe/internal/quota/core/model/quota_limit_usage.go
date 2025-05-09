package model

import "time"

type QuotaLimitUsage struct {
	ID            string     `json:"id" db:"id"`
	CorrelationID string     `json:"correlation_id" db:"correlation_id"`
	QuotaLimitID  string     `json:"quota_limit_id" db:"quota_limit_id"`
	CurrentUsage  int        `json:"current_usage" db:"current_usage"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	LastResetAt   *time.Time `json:"last_reset_at" db:"last_reset_at"`
	UpdatedAt     *time.Time `json:"updated_at" db:"updated_at"`
}
