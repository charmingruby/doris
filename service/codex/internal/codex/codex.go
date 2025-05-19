package codex

import (
	"github.com/charmingruby/doris/lib/delivery/http/rest"
	"github.com/charmingruby/doris/lib/delivery/messaging"
	"github.com/charmingruby/doris/lib/instrumentation"
	persistenceLib "github.com/charmingruby/doris/lib/persistence"
	"github.com/charmingruby/doris/lib/storage"
	"github.com/charmingruby/doris/lib/validation"
	"github.com/charmingruby/doris/service/codex/config"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/client"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/repository"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/usecase"
	"github.com/charmingruby/doris/service/codex/internal/codex/delivery/event"
	"github.com/charmingruby/doris/service/codex/internal/codex/delivery/http/rest/endpoint"
	"github.com/charmingruby/doris/service/codex/internal/codex/persistence"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Datasource struct {
	codexRepo              repository.CodexRepository
	codexDocumentRepo      repository.CodexDocumentRepository
	codexDocumentChunkRepo repository.CodexDocumentChunkRepository
	txManager              persistenceLib.TransactionManager[repository.TransactionManager]
}

func NewDatasource(db *sqlx.DB) (*Datasource, error) {
	codexRepo, err := persistence.NewCodexRepository(db)
	if err != nil {
		return nil, err
	}

	codexDocumentRepo, err := persistence.NewCodexDocumentRepository(db)
	if err != nil {
		return nil, err
	}

	codexDocumentChunkRepo, err := persistence.NewCodexDocumentChunkRepository(db)
	if err != nil {
		return nil, err
	}

	txManager, err := persistence.NewTransactionManager(db)
	if err != nil {
		return nil, err
	}

	return &Datasource{
		codexRepo:              codexRepo,
		codexDocumentRepo:      codexDocumentRepo,
		codexDocumentChunkRepo: codexDocumentChunkRepo,
		txManager:              txManager,
	}, nil
}

func NewUseCase(
	logger *instrumentation.Logger,
	datasource *Datasource,
	eventHandler *event.Handler,
	storage storage.Storage,
	quotaUsageManagementClient client.QuotaUsageManagement,
	embeddingSourceDocsBucket string,
) *usecase.UseCase {
	return usecase.New(
		logger,
		datasource.codexRepo,
		datasource.codexDocumentRepo,
		datasource.codexDocumentChunkRepo,
		storage,
		eventHandler,
		datasource.txManager,
		quotaUsageManagementClient,
		embeddingSourceDocsBucket,
	)
}

func NewEventHandler(logger *instrumentation.Logger, pub messaging.Publisher, sub messaging.Subscriber, cfg config.Config) (*event.Handler, error) {
	eventHander := event.NewHandler(logger, pub, sub, event.TopicInput{
		CodexDocumentUploaded: cfg.Custom.CodexDocumentUploadedTopic,
	})

	return eventHander, nil
}

func SubscribeEventHandler(eventHandler *event.Handler, uc *usecase.UseCase) error {
	if err := eventHandler.Subscribe(uc); err != nil {
		return err
	}

	return nil
}

func NewHTTPHandler(logger *instrumentation.Logger, r *gin.Engine, mw *rest.Middleware, val *validation.Validator, uc *usecase.UseCase) {
	endpoint.New(logger, r, mw, val, uc).Register()
}
