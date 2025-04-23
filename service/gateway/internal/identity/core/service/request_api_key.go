package service

import (
	"context"
	"time"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/core/id"
	"github.com/charmingruby/doris/service/gateway/internal/identity/core/model"
)

type RequestAPIKeyInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func (s *Service) RequestAPIKey(ctx context.Context, in RequestAPIKeyInput) error {
	apiKey, err := s.apiKeyRepo.FindByEmail(ctx, in.Email)

	if err != nil {
		s.logger.Error("error on find by email", "error", err)

		return custom_err.NewErrDatasourceOperationFailed("find api key by email", err)
	}

	if apiKey.ID != "" {
		return custom_err.NewErrResourceAlreadyExists("api key")
	}

	key := id.New()

	expirationDelay := 10 * time.Minute

	expirationDate := time.Now().Add(expirationDelay)

	ak := model.NewAPIKey(model.APIKeyInput{
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Email:     in.Email,
		Key:       key,
		ExpiresAt: expirationDate,
	})

	if err := s.apiKeyRepo.Create(ctx, *ak); err != nil {
		return custom_err.NewErrDatasourceOperationFailed("create api key", err)
	}

	if err := s.eventHandler.PublishRequestAPIKey(ctx, "notification", *ak); err != nil {
		return err
	}

	ak.Status = model.API_KEY_STATUS_PENDING

	if err := s.apiKeyRepo.Update(ctx, *ak); err != nil {
		return custom_err.NewErrDatasourceOperationFailed("update api key", err)
	}

	return nil
}
