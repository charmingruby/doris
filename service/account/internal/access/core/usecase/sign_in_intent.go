package usecase

import (
	"context"
	"time"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/service/account/internal/access/core/event"
	"github.com/charmingruby/doris/service/account/internal/access/core/model"
	"github.com/charmingruby/doris/service/account/internal/access/core/repository"
)

type SignInIntentInput struct {
	Email string `json:"email"`
}

func (uc *UseCase) SignInIntent(ctx context.Context, in SignInIntentInput) error {
	ak, err := uc.apiKeyRepo.FindByEmail(ctx, in.Email)

	if err != nil {
		return custom_err.NewErrDatasourceOperationFailed("find api key by email", err)
	}

	if ak.ID == "" {
		return custom_err.NewErrResourceNotFound("api key")
	}

	hasSufficientPermission := ak.Status == model.API_KEY_STATUS_ACTIVE || ak.Status == model.API_KEY_STATUS_DEFAULTER

	if !hasSufficientPermission {
		return custom_err.NewErrInsufficientPermission()
	}

	if err := uc.txManager.Transact(func(tx repository.TransactionManager) error {
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

		if err := uc.event.DispatchSendOTPNotification(ctx, event); err != nil {
			return custom_err.NewErrMessagingWrapper(err)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
