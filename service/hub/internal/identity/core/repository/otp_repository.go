package repository

import (
	"context"

	"github.com/charmingruby/doris/service/hub/internal/identity/core/model"
)

type OTPRepository interface {
	Create(ctx context.Context, otp model.OTP) error
	FindByCorrelationID(ctx context.Context, correlationID string) (model.OTP, error)
}
