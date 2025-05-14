package model

import (
	"errors"
	"time"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/core/id"
	"github.com/charmingruby/doris/lib/core/privilege"
	"github.com/charmingruby/doris/service/scribe/internal/shared/core/kind"
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
	Tier     string
	Kind     string
	MaxValue int
	Unit     string
}

func NewQuota(in QuotaInput) (*Quota, error) {
	if err := privilege.IsTierValid(in.Tier); err != nil {
		return nil, err
	}

	if err := kind.IsValid(in.Kind); err != nil {
		return nil, err
	}

	return &Quota{
		ID:        id.New(),
		Tier:      in.Tier,
		Kind:      in.Kind,
		MaxValue:  in.MaxValue,
		Unit:      in.Unit,
		Status:    QUOTA_STATUS_DRAFT,
		CreatedAt: time.Now(),
		UpdatedAt: nil,
	}, nil
}

type Quota struct {
	ID        string     `json:"id" db:"id"`
	Tier      string     `json:"tier" db:"tier"`
	Kind      string     `json:"kind" db:"kind"`
	MaxValue  int        `json:"max_value" db:"max_value"`
	Unit      string     `json:"unit" db:"unit"`
	Status    string     `json:"status" db:"status"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

type ModifyQuotaInput struct {
	Tier     string
	Kind     string
	MaxValue int
	Unit     string
	Status   string
}

func (q *Quota) Modify(in ModifyQuotaInput) error {
	hasChange := false

	if in.Tier != "" && in.Tier != q.Tier {
		if err := privilege.IsTierValid(in.Tier); err != nil {
			return custom_err.NewErrInvalidEntity(err.Error())
		}

		hasChange = true
	}

	if in.Status != "" && in.Status != q.Status {
		if _, ok := validQuotaStatus[in.Status]; !ok {
			return custom_err.NewErrInvalidEntity(ErrInvalidQuotaStatus.Error())
		}

		hasChange = true
	}

	if in.Kind != "" && in.Kind != q.Kind {
		if err := kind.IsValid(in.Kind); err != nil {
			return custom_err.NewErrInvalidEntity(err.Error())
		}

		hasChange = true
	}

	if in.MaxValue != 0 && in.MaxValue != q.MaxValue {
		hasChange = true
	}

	if in.Unit != "" && in.Unit != q.Unit {
		hasChange = true
	}

	if !hasChange {
		return custom_err.NewErrNothingToChange()
	}

	if in.Kind != "" && in.Kind != q.Kind {
		q.Kind = in.Kind
	}

	if in.MaxValue != 0 && in.MaxValue != q.MaxValue {
		q.MaxValue = in.MaxValue
	}

	if in.Unit != "" && in.Unit != q.Unit {
		q.Unit = in.Unit
	}

	if in.Tier != "" {
		q.Tier = in.Tier
	}

	if in.Status != "" {
		q.Status = in.Status
	}

	now := time.Now()
	q.UpdatedAt = &now

	return nil
}
