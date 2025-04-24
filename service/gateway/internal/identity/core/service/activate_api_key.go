package service

import (
	"context"
	"time"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/service/gateway/internal/identity/core/model"
)

type ActivateAPIKeyInput struct {
	Key            string `json:"key"`
	ActivationCode string `json:"activation_code"`
}

func (s *Service) ActivateAPIKey(ctx context.Context, in ActivateAPIKeyInput) error {
	ak, err := s.apiKeyRepo.FindByKey(ctx, in.Key)

	if err != nil {
		return custom_err.NewErrDatasourceOperationFailed("find api key by key", err)
	}

	if ak.ID == "" {
		return custom_err.NewErrResourceNotFound("api key")
	}

	if ak.ActivationCode != in.ActivationCode {
		return custom_err.NewErrInvalidConfirmationCode("does not match")
	}

	if ak.ActivationCodeExpiresAt.Before(time.Now()) {
		return custom_err.NewErrInvalidConfirmationCode("expired")
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
