package model

import "time"

type QuotaLimit struct {
	ID        string     `json:"id" db:"id"`
	QuotaID   string     `json:"quota_id" db:"quota_id"`
	Kind      string     `json:"kind" db:"kind"`
	MaxValue  int        `json:"max_value" db:"max_value"`
	Unit      string     `json:"unit" db:"unit"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}
