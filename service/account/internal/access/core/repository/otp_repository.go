package repository

import (
	"context"

	"github.com/charmingruby/doris/service/account/internal/access/core/model"
)

type OTPRepository interface {
	Create(ctx context.Context, otp model.OTP) error
	FindMostRecentByCorrelationID(ctx context.Context, correlationID string) (model.OTP, error)
}
