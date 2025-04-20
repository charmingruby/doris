package usecase

import (
	"context"
	"time"

	"github.com/charmingruby/doris/lib/core"
	"github.com/charmingruby/doris/service/gateway/internal/identity/core/model"
)

type RequestAPIKeyInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func (u *UseCase) RequestAPIKey(ctx context.Context, in RequestAPIKeyInput) error {
	key := core.NewID()

	expirationDate := time.Now().AddDate(0, 0, 1)

	ak := model.NewAPIKey(model.APIKeyInput{
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Email:     in.Email,
		Key:       key,
		ExpiresAt: expirationDate,
	})

	if err := u.apiKeyRepo.Create(ctx, *ak); err != nil {
		return err
	}

	// publish to the notification queue

	ak.Status = model.API_KEY_STATUS_PENDING

	if err := u.apiKeyRepo.Update(ctx, *ak); err != nil {
		return err
	}

	return nil
}
