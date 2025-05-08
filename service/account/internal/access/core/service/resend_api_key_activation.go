package service

import (
	"context"
	"time"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/service/account/internal/access/core/event"
	"github.com/charmingruby/doris/service/account/internal/access/core/model"
	"github.com/charmingruby/doris/service/account/internal/access/core/repository"
)

type ResendAPIKeyActivationInput struct {
	APIKeyID string `json:"api_key_id"`
}

func (s *Service) ResendAPIKeyActivation(ctx context.Context, in ResendAPIKeyActivationInput) error {
	ak, err := s.apiKeyRepo.FindByID(ctx, in.APIKeyID)

	if err != nil {
		return custom_err.NewErrDatasourceOperationFailed("find api key by id", err)
	}

	if ak.ID == "" {
		return custom_err.NewErrResourceNotFound("api key")
	}

	isOTPAlreadySent := true

	otp, err := s.otpRepo.FindMostRecentByCorrelationID(ctx, ak.ID)

	if err != nil {
		return custom_err.NewErrDatasourceOperationFailed("find otp by correlation id", err)
	}

	if otp.ID == "" {
		isOTPAlreadySent = false
	}

	if isOTPAlreadySent {
		otpGenerationDelay := time.Second * 30

		if otp.CreatedAt.Add(otpGenerationDelay).After(time.Now()) {
			return custom_err.NewErrOTPGenerationCooldown()
		}
	}

	if err := s.txManager.Transact(func(tx repository.TransactionManager) error {
		otp, err := model.NewOTP(model.OTPInput{
			Purpose:       model.OTP_PURPOSE_API_KEY_ACTIVATION,
			CorrelationID: ak.ID,
			ExpiresAt:     time.Now().UTC().Add(30 * time.Minute),
		})

		if err != nil {
			return custom_err.NewErrInvalidEntity(err.Error())
		}

		if err := tx.OTPRepo.Create(ctx, *otp); err != nil {
			return custom_err.NewErrDatasourceOperationFailed("create otp", err)
		}

		event := event.SendOTPNotification{
			ID:            ak.ID,
			To:            ak.Email,
			RecipientName: ak.FirstName + " " + ak.LastName,
			Code:          otp.Code,
			SentAt:        time.Now(),
		}

		if err := s.event.DispatchSendOTPNotification(ctx, event); err != nil {
			return custom_err.NewErrMessagingWrapper(err)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
