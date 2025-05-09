package model

import "time"

type Quota struct {
	ID        string     `json:"id" db:"id"`
	Tier      string     `json:"tier" db:"tier"`
	Status    string     `json:"status" db:"status"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}
