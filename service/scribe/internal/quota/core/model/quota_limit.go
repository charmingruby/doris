package model

import (
	"errors"
	"time"

	"github.com/charmingruby/doris/lib/core/custom_err"
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
		IsActive:  true,
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
	IsActive  bool       `json:"is_active" db:"is_active"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

type ModifyQuotaLimitInput struct {
	Kind     string
	MaxValue int
	Unit     string
	IsActive bool
}

func (q *QuotaLimit) Modify(in ModifyQuotaLimitInput) error {
	hasChange := in.Kind != q.Kind ||
		in.MaxValue != q.MaxValue ||
		in.Unit != q.Unit ||
		in.IsActive != q.IsActive

	if !hasChange {
		return custom_err.NewErrNothingToChange()
	}

	q.Kind = in.Kind
	q.MaxValue = in.MaxValue
	q.Unit = in.Unit
	q.IsActive = in.IsActive

	now := time.Now()
	q.UpdatedAt = &now

	return nil
}
