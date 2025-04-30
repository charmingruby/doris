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

func (r *OTPRepository) FindByCorrelationID(ctx context.Context, correlationID string) (model.OTP, error) {
	for _, i := range r.Items {
		if i.CorrelationID == correlationID {
			return i, nil
		}
	}

	if !r.IsHealthy {
		return model.OTP{}, ErrUnhealthyDatasource
	}

	return model.OTP{}, nil
}
