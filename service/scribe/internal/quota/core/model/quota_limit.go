package model

import (
	"errors"
	"time"

	"github.com/charmingruby/doris/lib/core/id"
)

const (
	QUOTA_LIMIT_KIND_DOCUMENT = "DOCUMENT"
	QUOTA_LIMIT_KIND_REQUEST  = "REQUEST"
)

var (
	ErrInvalidKind = errors.New("invalid kind")

	validKinds = map[string]struct{}{
		QUOTA_LIMIT_KIND_DOCUMENT: {},
		QUOTA_LIMIT_KIND_REQUEST:  {},
	}
)

type QuotaLimitInput struct {
	QuotaID  string
	Kind     string
	MaxValue int
	Unit     string
}

func NewQuotaLimit(in QuotaLimitInput) (*QuotaLimit, error) {
	if _, ok := validKinds[in.Kind]; !ok {
		return nil, ErrInvalidKind
	}

	return &QuotaLimit{
		ID:        id.New(),
		QuotaID:   in.QuotaID,
		Kind:      in.Kind,
		MaxValue:  in.MaxValue,
		Unit:      in.Unit,
		CreatedAt: time.Now(),
		UpdatedAt: nil,
	}, nil
}

type QuotaLimit struct {
	ID        string     `json:"id" db:"id"`
	QuotaID   string     `json:"quota_id" db:"quota_id"`
	Kind      string     `json:"kind" db:"kind"`
	MaxValue  int        `json:"max_value" db:"max_value"`
	Unit      string     `json:"unit" db:"unit"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}
