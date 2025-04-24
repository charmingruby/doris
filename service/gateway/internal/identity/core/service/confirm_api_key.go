package service

import (
	"context"
	"time"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/service/gateway/internal/identity/core/model"
)

type ConfirmAPIKeyInput struct {
	Key              string `json:"key"`
	ConfirmationCode string `json:"confirmation_code"`
}

func (s *Service) ConfirmAPIKey(ctx context.Context, in ConfirmAPIKeyInput) error {
	ak, err := s.apiKeyRepo.FindByKey(ctx, in.Key)

	if err != nil {
		return custom_err.NewErrDatasourceOperationFailed("find api key by key", err)
	}

	if ak.ID == "" {
		return custom_err.NewErrResourceNotFound("api key")
	}

	if ak.ConfirmationCode != in.ConfirmationCode {
		return custom_err.NewErrInvalidConfirmationCode("does not match")
	}

	if ak.ConfirmationCodeExpiresAt.Before(time.Now()) {
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
