package model

import (
	"errors"
	"time"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/core/id"
	"github.com/charmingruby/doris/lib/core/privilege"
)

const (
	QUOTA_STATUS_DRAFT    = "DRAFT"
	QUOTA_STATUS_DISABLED = "DISABLED"
	QUOTA_STATUS_ENABLED  = "ENABLED"
)

var (
	ErrInvalidQuotaStatus = errors.New("invalid quota status")

	validQuotaStatus = map[string]struct{}{
		QUOTA_STATUS_DRAFT:    {},
		QUOTA_STATUS_DISABLED: {},
		QUOTA_STATUS_ENABLED:  {},
	}
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

type ModifyQuotaInput struct {
	Tier   string
	Status string
}

func (q *Quota) Modify(in ModifyQuotaInput) error {
	hasChange := in.Tier != q.Tier || in.Status != q.Status

	if !hasChange {
		return custom_err.NewErrNothingToChange()
	}

	if err := privilege.IsTierValid(in.Tier); err != nil {
		return custom_err.NewErrInvalidEntity(err.Error())
	}

	if _, exists := validQuotaStatus[in.Status]; !exists {
		return custom_err.NewErrInvalidEntity(ErrInvalidQuotaStatus.Error())
	}

	q.Tier = in.Tier
	q.Status = in.Status

	now := time.Now()
	q.UpdatedAt = &now

	return nil
}
