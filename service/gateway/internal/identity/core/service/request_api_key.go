package service

import (
	"context"
	"errors"
	"time"

	"github.com/charmingruby/doris/lib/core"
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
		return err
	}

	if apiKey.ID != "" {
		return errors.New("api key already exists")
	}

	key := core.NewID()

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
		return err
	}

	// publish to the notification queue

	ak.Status = model.API_KEY_STATUS_PENDING

	if err := s.apiKeyRepo.Update(ctx, *ak); err != nil {
		return err
	}

	return nil
}
