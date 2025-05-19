package usecase

import (
	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/lib/persistence"
	"github.com/charmingruby/doris/lib/storage"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/client"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/event"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/repository"
)

type UseCase struct {
	logger                     *instrumentation.Logger
	codexRepo                  repository.CodexRepository
	codexDocumentRepo          repository.CodexDocumentRepository
	codexDocumentChunkRepo     repository.CodexDocumentChunkRepository
	txManager                  persistence.TransactionManager[repository.TransactionManager]
	storage                    storage.Storage
	eventHandler               event.Handler
	quotaUsageManagementClient client.QuotaUsageManagement
	embeddingSourceDocsBucket  string
}

func New(
	logger *instrumentation.Logger,
	codexRepo repository.CodexRepository,
	codexDocumentRepo repository.CodexDocumentRepository,
	codexDocumentChunkRepo repository.CodexDocumentChunkRepository,
	storage storage.Storage,
	eventHandler event.Handler,
	txManager persistence.TransactionManager[repository.TransactionManager],
	quotaUsageManagementClient client.QuotaUsageManagement,
	embeddingSourceDocsBucket string,
) *UseCase {
	return &UseCase{
		logger:                     logger,
		codexRepo:                  codexRepo,
		codexDocumentRepo:          codexDocumentRepo,
		codexDocumentChunkRepo:     codexDocumentChunkRepo,
		storage:                    storage,
		eventHandler:               eventHandler,
		quotaUsageManagementClient: quotaUsageManagementClient,
		txManager:                  txManager,
		embeddingSourceDocsBucket:  embeddingSourceDocsBucket,
	}
}
