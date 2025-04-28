package service

import (
	"context"
	"time"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/service/hub/internal/identity/core/model"
)

type ActivateAPIKeyInput struct {
	APIKeyID string `json:"api_key_id"`
	OTPCode  string `json:"otp_code"`
}

func (s *Service) ActivateAPIKey(ctx context.Context, in ActivateAPIKeyInput) error {
	ak, err := s.apiKeyRepo.FindByID(ctx, in.APIKeyID)

	if err != nil {
		return custom_err.NewErrDatasourceOperationFailed("find api key by id", err)
	}

	if ak.ID == "" {
		return custom_err.NewErrResourceNotFound("api key")
	}

	otp, err := s.otpRepo.FindByCorrelationID(ctx, ak.ID)

	if err != nil {
		return custom_err.NewErrDatasourceOperationFailed("find otp by correlation id", err)
	}

	if otp.ID == "" {
		return custom_err.NewErrResourceNotFound("otp")
	}

	if otp.Code != in.OTPCode {
		return custom_err.NewErrInvalidOTPCode("does not match")
	}

	if otp.ExpiresAt.Before(time.Now()) {
		return custom_err.NewErrInvalidOTPCode("expired")
	}

	if ak.Status == model.API_KEY_STATUS_ACTIVE {
		return custom_err.NewErrAPIKeyAlreadyConfirmed()
	}

	ak.Status = model.API_KEY_STATUS_ACTIVE

	if err := s.apiKeyRepo.Update(ctx, ak); err != nil {
		return custom_err.NewErrDatasourceOperationFailed("update api key", err)
	}

	return nil
}
