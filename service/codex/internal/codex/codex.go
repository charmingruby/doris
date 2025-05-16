package codex

import (
	"github.com/charmingruby/doris/lib/delivery/http/rest"
	"github.com/charmingruby/doris/lib/delivery/messaging"
	"github.com/charmingruby/doris/lib/instrumentation"
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
	codexRepo         repository.CodexRepository
	codexDocumentRepo repository.CodexDocumentRepository
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

	return &Datasource{
		codexRepo:         codexRepo,
		codexDocumentRepo: codexDocumentRepo,
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
		storage,
		eventHandler,
		quotaUsageManagementClient,
		embeddingSourceDocsBucket,
	)
}

func NewEventHandler(logger *instrumentation.Logger, pub messaging.Publisher, cfg config.Config) *event.Handler {
	return event.NewHandler(logger, pub, event.TopicInput{
		CodexDocumentUploaded: cfg.Custom.CodexDocumentUploadedTopic,
	})
}

func NewHTTPHandler(logger *instrumentation.Logger, r *gin.Engine, mw *rest.Middleware, val *validation.Validator, uc *usecase.UseCase) {
	endpoint.New(logger, r, mw, val, uc).Register()
}
