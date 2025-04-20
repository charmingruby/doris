package repository

import (
	"context"

	"github.com/charmingruby/doris/service/gateway/internal/identity/core/model"
)

type APIKeyRepository interface {
	Create(ctx context.Context, apiKey model.APIKey) error
	Update(ctx context.Context, apiKey model.APIKey) error
}
