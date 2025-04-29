package repository

import (
	"context"

	"github.com/charmingruby/doris/service/identity/internal/access/core/model"
)

type OTPRepository interface {
	Create(ctx context.Context, otp model.OTP) error
	FindByCorrelationID(ctx context.Context, correlationID string) (model.OTP, error)
}
