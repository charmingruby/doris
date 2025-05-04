package service

import (
	"context"
	"time"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/security"
)

type VerifySignInIntentInput struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

func (s *Service) VerifySignInIntent(ctx context.Context, in VerifySignInIntentInput) (string, error) {
	ak, err := s.apiKeyRepo.FindByEmail(ctx, in.Email)

	if err != nil {
		return "", custom_err.NewErrDatasourceOperationFailed("find api key by email", err)
	}

	if ak.ID == "" {
		return "", custom_err.NewErrResourceNotFound("api key")
	}

	otp, err := s.otpRepo.FindMostRecentByCorrelationID(ctx, ak.ID)

	if err != nil {
		return "", custom_err.NewErrDatasourceOperationFailed("find otp by correlation id", err)
	}

	if otp.ID == "" {
		return "", custom_err.NewErrResourceNotFound("otp")
	}

	if otp.Code != in.OTP {
		return "", custom_err.NewErrInvalidOTPCode("does not match")
	}

	if otp.ExpiresAt.Before(time.Now().UTC()) {
		return "", custom_err.NewErrInvalidOTPCode("expired")
	}

	return s.tokenClient.Generate(ak.ID, security.Payload{
		Tier: ak.Tier,
	})
}
