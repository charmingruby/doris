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
	apiKey, err := s.repo.FindByEmail(ctx, in.Email)

	if err != nil {
		s.logger.Error("error on find by email", "error", err)

		return "", custom_err.NewErrDatasourceOperationFailed("find api key by email", err)
	}

	if apiKey.ID != "" {
		return "", custom_err.NewErrResourceAlreadyExists("api key")
	}

	key := id.New()

	expirationDelay := 10 * time.Minute

	expirationDate := time.Now().Add(expirationDelay)

	ak := model.NewAPIKey(model.APIKeyInput{
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Email:     in.Email,
		Key:       key,
		// TODO: create a confirmation code logic
		ActivationCode:          id.New(),
		ActivationCodeExpiresAt: expirationDate,
	})

	if err := s.repo.Create(ctx, *ak); err != nil {
		return "", custom_err.NewErrDatasourceOperationFailed("create api key", err)
	}

	event := &event.APIKeyActivation{
		ID:             ak.ID,
		To:             ak.Email,
		RecipientName:  ak.FirstName + " " + ak.LastName,
		ActivationCode: ak.ActivationCode,
		SentAt:         time.Now(),
	}

	if err := s.event.SendAPIKeyActivation(ctx, event); err != nil {
		return "", err
	}

	ak.Status = model.API_KEY_STATUS_PENDING

	if err := s.repo.Update(ctx, *ak); err != nil {
		return "", custom_err.NewErrDatasourceOperationFailed("update api key", err)
	}

	return ak.ID, nil
}
