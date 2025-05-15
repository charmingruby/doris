package usecase

import (
	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/lib/storage"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/repository"
)

type UseCase struct {
	logger                    *instrumentation.Logger
	codexRepo                 repository.CodexRepository
	storage                   storage.Storage
	embeddingSourceDocsBucket string
}

func New(logger *instrumentation.Logger, codexRepo repository.CodexRepository, storage storage.Storage, embeddingSourceDocsBucket string) *UseCase {
	return &UseCase{
		logger:                    logger,
		codexRepo:                 codexRepo,
		storage:                   storage,
		embeddingSourceDocsBucket: embeddingSourceDocsBucket,
	}
}
