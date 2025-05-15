package memory

import (
	"context"

	"github.com/charmingruby/doris/service/codex/internal/codex/core/model"
)

type CodexRepository struct {
	Items     []model.Codex
	IsHealthy bool
}

func NewCodexRepository() *CodexRepository {
	return &CodexRepository{
		Items:     []model.Codex{},
		IsHealthy: true,
	}
}

func (r *CodexRepository) FindByCorrelationIDAndName(ctx context.Context, correlationID, name string) (model.Codex, error) {
	if !r.IsHealthy {
		return model.Codex{}, ErrUnhealthyDatasource
	}

	for _, i := range r.Items {
		if i.CorrelationID == correlationID && i.Name == name {
			return i, nil
		}
	}

	return model.Codex{}, nil
}

func (r *CodexRepository) Create(ctx context.Context, codex model.Codex) error {
	if !r.IsHealthy {
		return ErrUnhealthyDatasource
	}

	r.Items = append(r.Items, codex)

	return nil
}
