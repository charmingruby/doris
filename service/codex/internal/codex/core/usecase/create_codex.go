package usecase

import (
	"context"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/model"
)

type CreateCodexInput struct {
	CorrelationID string
	Name          string
	Description   string
}

func (u *UseCase) CreateCodex(ctx context.Context, in *CreateCodexInput) (string, error) {
	codex, err := u.codexRepo.FindByCorrelationIDAndName(ctx, in.CorrelationID, in.Name)

	if err != nil {
		return "", custom_err.NewErrDatasourceOperationFailed("create codex", err)
	}

	if codex.ID != "" {
		return "", custom_err.NewErrResourceAlreadyExists("codex")
	}

	c := model.NewCodex(model.CodexInput{
		CorrelationID: in.CorrelationID,
		Name:          in.Name,
		Description:   in.Description,
	})

	if err := u.codexRepo.Create(ctx, *c); err != nil {
		return "", custom_err.NewErrDatasourceOperationFailed("create codex", err)
	}

	return c.ID, nil
}
