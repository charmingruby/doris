package repository

import (
	"context"

	"github.com/charmingruby/doris/service/hub/internal/identity/core/model"
)

type APIKeyRepository interface {
	FindByID(ctx context.Context, id string) (model.APIKey, error)
	FindByEmail(ctx context.Context, email string) (model.APIKey, error)
	FindByKey(ctx context.Context, key string) (model.APIKey, error)
	Create(ctx context.Context, apiKey model.APIKey) error
	Update(ctx context.Context, apiKey model.APIKey) error
}
