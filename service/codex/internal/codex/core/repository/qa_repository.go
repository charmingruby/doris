package repository

import (
	"context"

	"github.com/charmingruby/doris/service/codex/internal/codex/core/model"
)

type QARepository interface {
	Create(ctx context.Context, qa model.QA) error
}
