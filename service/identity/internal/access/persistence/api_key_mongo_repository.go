package persistence

import (
	"context"

	"github.com/charmingruby/doris/service/identity/internal/access/core/model"
)

type APIKeyMongoRepo struct{}

func NewAPIKeyMongoRepo() *APIKeyMongoRepo {
	return &APIKeyMongoRepo{}
}

func (r *APIKeyMongoRepo) FindByID(ctx context.Context, id string) (model.APIKey, error) {
	return model.APIKey{}, nil
}

func (r *APIKeyMongoRepo) FindByEmail(ctx context.Context, email string) (model.APIKey, error) {
	return model.APIKey{}, nil
}

func (r *APIKeyMongoRepo) FindByKey(ctx context.Context, key string) (model.APIKey, error) {
	return model.APIKey{}, nil
}

func (r *APIKeyMongoRepo) Create(ctx context.Context, apiKey model.APIKey) error {
	return nil
}

func (r *APIKeyMongoRepo) Update(ctx context.Context, apiKey model.APIKey) error {
	return nil
}
