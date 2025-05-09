package usecase

import (
	"context"
	"time"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/core/id"
	"github.com/charmingruby/doris/service/account/internal/access/core/event"
	"github.com/charmingruby/doris/service/account/internal/access/core/model"
	"github.com/charmingruby/doris/service/account/internal/access/core/repository"
)

type GenerateAPIKeyInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func (uc *UseCase) GenerateAPIKey(ctx context.Context, in GenerateAPIKeyInput) (string, error) {
	apiKey, err := uc.apiKeyRepo.FindByEmail(ctx, in.Email)

	if err != nil {
		return "", custom_err.NewErrDatasourceOperationFailed("find api key by email", err)
	}

	if apiKey.ID != "" {
		return "", custom_err.NewErrResourceAlreadyExists("api key")
	}

	ak := model.NewAPIKey(model.APIKeyInput{
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Email:     in.Email,
		Key:       id.New(),
	})

	if err := uc.txManager.Transact(func(tx repository.TransactionManager) error {
		if err := tx.APIKeyRepo.Create(ctx, *ak); err != nil {
			return custom_err.NewErrDatasourceOperationFailed("create api key", err)
		}

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

		ak.Status = model.API_KEY_STATUS_PENDING

		if err := tx.APIKeyRepo.Update(ctx, *ak); err != nil {
			return custom_err.NewErrDatasourceOperationFailed("update api key", err)
		}

		return nil
	}); err != nil {
		return "", err
	}

	return ak.ID, nil
}
