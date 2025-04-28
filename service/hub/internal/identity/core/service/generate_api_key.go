package service

import (
	"context"
	"time"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/core/id"
	"github.com/charmingruby/doris/service/hub/internal/identity/core/event"
	"github.com/charmingruby/doris/service/hub/internal/identity/core/model"
)

type GenerateAPIKeyInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func (s *Service) GenerateAPIKey(ctx context.Context, in GenerateAPIKeyInput) (string, error) {
	apiKey, err := s.apiKeyRepo.FindByEmail(ctx, in.Email)

	if err != nil {
		s.logger.Error("error on find by email", "error", err)

		return "", custom_err.NewErrDatasourceOperationFailed("find api key by email", err)
	}

	if apiKey.ID != "" {
		return "", custom_err.NewErrResourceAlreadyExists("api key")
	}

	key := id.New()

	expirationDelay := 30 * time.Minute

	expirationDate := time.Now().Add(expirationDelay)

	ak := model.NewAPIKey(model.APIKeyInput{
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Email:     in.Email,
		Key:       key,
	})

	if err := s.apiKeyRepo.Create(ctx, *ak); err != nil {
		return "", custom_err.NewErrDatasourceOperationFailed("create api key", err)
	}

	otp, err := model.NewOTP(model.OTPInput{
		Purpose:       model.OTP_PURPOSE_API_KEY_ACTIVATION,
		CorrelationID: ak.ID,
		ExpiresAt:     expirationDate,
	})

	if err != nil {
		return "", custom_err.NewErrInvalidEntity("otp")
	}

	if err := s.otpRepo.Create(ctx, *otp); err != nil {
		return "", custom_err.NewErrDatasourceOperationFailed("create otp", err)
	}

	event := &event.OTP{
		ID:            ak.ID,
		To:            ak.Email,
		RecipientName: ak.FirstName + " " + ak.LastName,
		Code:          otp.Code,
		SentAt:        time.Now(),
	}

	if err := s.event.SendOTP(ctx, event); err != nil {
		return "", err
	}

	ak.Status = model.API_KEY_STATUS_PENDING

	if err := s.apiKeyRepo.Update(ctx, *ak); err != nil {
		return "", custom_err.NewErrDatasourceOperationFailed("update api key", err)
	}

	return ak.ID, nil
}
