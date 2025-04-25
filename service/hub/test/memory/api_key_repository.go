package memory

import (
	"context"
	"errors"

	"github.com/charmingruby/doris/service/hub/internal/identity/core/model"
)

var ErrUnhealthyDatasource = errors.New("datasource is unhealthy")

type APIKeyRepository struct {
	Items     []model.APIKey
	IsHealthy bool
}

func NewAPIKeyRepository() *APIKeyRepository {
	return &APIKeyRepository{
		Items:     []model.APIKey{},
		IsHealthy: true,
	}
}

func (r *APIKeyRepository) FindByID(ctx context.Context, id string) (model.APIKey, error) {
	for _, i := range r.Items {
		if i.ID == id {
			return i, nil
		}
	}

	if !r.IsHealthy {
		return model.APIKey{}, ErrUnhealthyDatasource
	}

	return model.APIKey{}, nil
}

func (r *APIKeyRepository) FindByEmail(ctx context.Context, email string) (model.APIKey, error) {
	for _, i := range r.Items {
		if i.Email == email {
			return i, nil
		}
	}

	if !r.IsHealthy {
		return model.APIKey{}, ErrUnhealthyDatasource
	}

	return model.APIKey{}, nil
}

func (r *APIKeyRepository) FindByKey(ctx context.Context, key string) (model.APIKey, error) {
	for _, i := range r.Items {
		if i.Key == key {
			return i, nil
		}
	}

	if !r.IsHealthy {
		return model.APIKey{}, ErrUnhealthyDatasource
	}

	return model.APIKey{}, nil
}

func (r *APIKeyRepository) Create(ctx context.Context, apiKey model.APIKey) error {
	if !r.IsHealthy {
		return ErrUnhealthyDatasource
	}

	r.Items = append(r.Items, apiKey)

	return nil
}

func (r *APIKeyRepository) Update(ctx context.Context, apiKey model.APIKey) error {
	if !r.IsHealthy {
		return ErrUnhealthyDatasource
	}

	for idx, i := range r.Items {
		if i.ID == apiKey.ID {
			r.Items[idx] = apiKey
		}
	}

	return nil
}
