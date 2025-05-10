package model

import (
	"time"

	"github.com/charmingruby/doris/lib/core/id"
	"github.com/charmingruby/doris/lib/core/privilege"
)

const (
	QUOTA_STATUS_DRAFT    = "DRAFT"
	QUOTA_STATUS_DISABLED = "DISABLED"
	QUOTA_STATUS_ENABLED  = "ENABLED"
)

type QuotaInput struct {
	Tier string
}

func NewQuota(in QuotaInput) (*Quota, error) {
	if err := privilege.IsTierValid(in.Tier); err != nil {
		return nil, err
	}

	return &Quota{
		ID:        id.New(),
		Tier:      in.Tier,
		Status:    QUOTA_STATUS_DRAFT,
		CreatedAt: time.Now(),
		UpdatedAt: nil,
	}, nil
}

type Quota struct {
	ID        string     `json:"id" db:"id"`
	Tier      string     `json:"tier" db:"tier"`
	Status    string     `json:"status" db:"status"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}
