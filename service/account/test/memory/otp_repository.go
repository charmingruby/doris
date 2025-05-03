package memory

import (
	"context"

	"github.com/charmingruby/doris/service/account/internal/access/core/model"
)

type OTPRepository struct {
	Items     []model.OTP
	IsHealthy bool
}

func NewOTPRepository() *OTPRepository {
	return &OTPRepository{
		Items:     []model.OTP{},
		IsHealthy: true,
	}
}

func (r *OTPRepository) Create(ctx context.Context, otp model.OTP) error {
	if !r.IsHealthy {
		return ErrUnhealthyDatasource
	}

	r.Items = append(r.Items, otp)

	return nil
}

func (r *OTPRepository) FindMostRecentByCorrelationID(ctx context.Context, correlationID string) (model.OTP, error) {
	if !r.IsHealthy {
		return model.OTP{}, ErrUnhealthyDatasource
	}

	var mostRecent model.OTP
	found := false

	for _, i := range r.Items {
		if i.CorrelationID == correlationID {
			if !found || i.CreatedAt.After(mostRecent.CreatedAt) {
				mostRecent = i
				found = true
			}
		}
	}

	if !found {
		return model.OTP{}, nil
	}

	return mostRecent, nil
}
