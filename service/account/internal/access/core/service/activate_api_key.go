package service

import (
	"context"
	"time"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/security"
	"github.com/charmingruby/doris/service/account/internal/access/core/model"
)

type ActivateAPIKeyInput struct {
	APIKeyID string `json:"api_key_id"`
	OTP      string `json:"otp"`
}

func (s *Service) ActivateAPIKey(ctx context.Context, in ActivateAPIKeyInput) (string, error) {
	ak, err := s.apiKeyRepo.FindByID(ctx, in.APIKeyID)

	if err != nil {
		return "", custom_err.NewErrDatasourceOperationFailed("find api key by id", err)
	}

	if ak.ID == "" {
		return "", custom_err.NewErrResourceNotFound("api key")
	}

	if ak.Status == model.API_KEY_STATUS_ACTIVE {
		return "", custom_err.NewErrAPIKeyAlreadyActivated()
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

	ak.Status = model.API_KEY_STATUS_ACTIVE

	if err := s.apiKeyRepo.Update(ctx, ak); err != nil {
		return "", custom_err.NewErrDatasourceOperationFailed("update api key", err)
	}

	return s.tokenClient.Generate(ak.ID, security.Payload{
		Tier: ak.Tier,
	})
}
