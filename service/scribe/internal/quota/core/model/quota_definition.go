package model

import "time"

type QuotaDefinition struct {
	ID        string     `json:"id" db:"id"`
	QuotaID   string     `json:"quota_id" db:"quota_id"`
	Tier      string     `json:"tier" db:"tier"`
	Status    string     `json:"status" db:"status"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}
