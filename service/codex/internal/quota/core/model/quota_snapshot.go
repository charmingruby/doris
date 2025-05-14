package model

import "time"

type QuotaSnapshot struct {
	CorrelationID string     `db:"correlation_id"`
	Kind          string     `db:"kind"`
	CurrentUsage  int        `db:"current_usage"`
	MaxValue      int        `db:"max_value"`
	Unit          string     `db:"unit"`
	Tier          string     `db:"tier"`
	Status        string     `db:"status"`
	LastResetAt   *time.Time `db:"last_reset_at"`
}
