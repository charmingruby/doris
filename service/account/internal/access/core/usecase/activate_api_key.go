package usecase

import (
	"context"
	"time"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/security"
	"github.com/charmingruby/doris/service/account/internal/access/core/event"
	"github.com/charmingruby/doris/service/account/internal/access/core/model"
	"github.com/charmingruby/doris/service/account/internal/access/core/repository"
)

type ActivateAPIKeyInput struct {
	APIKeyID string `json:"api_key_id"`
	OTP      string `json:"otp"`
}

func (uc *UseCase) ActivateAPIKey(ctx context.Context, in ActivateAPIKeyInput) (string, error) {
	ak, err := uc.apiKeyRepo.FindByID(ctx, in.APIKeyID)

	if err != nil {
		return "", custom_err.NewErrDatasourceOperationFailed("find api key by id", err)
	}

	if ak.ID == "" {
		return "", custom_err.NewErrResourceNotFound("api key")
	}

	if ak.Status == model.API_KEY_STATUS_ACTIVE {
		return "", custom_err.NewErrAPIKeyAlreadyActivated()
	}

	otp, err := uc.otpRepo.FindMostRecentByCorrelationID(ctx, ak.ID)

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

	if err := uc.txManager.Transact(func(tx repository.TransactionManager) error {
		ak.Status = model.API_KEY_STATUS_ACTIVE
		if err := tx.APIKeyRepo.Save(ctx, ak); err != nil {
			return custom_err.NewErrDatasourceOperationFailed("save api key", err)
		}

		event := event.APIKeyActivated{
			ID:     ak.ID,
			Tier:   ak.Tier,
			SentAt: time.Now(),
		}

		if err := uc.event.DispatchAPIKeyActivated(ctx, event); err != nil {
			return custom_err.NewErrMessagingWrapper(err)
		}

		return nil
	}); err != nil {
		return "", err
	}

	return uc.tokenClient.Generate(ak.ID, security.Payload{
		Tier: ak.Tier,
	})
}
