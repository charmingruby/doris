package quota

import (
	"github.com/charmingruby/doris/lib/delivery/http/rest"
	"github.com/charmingruby/doris/lib/delivery/messaging/nats"
	"github.com/charmingruby/doris/lib/instrumentation"
	persistenceLib "github.com/charmingruby/doris/lib/persistence"
	"github.com/charmingruby/doris/lib/validation"
	"github.com/charmingruby/doris/service/scribe/config"
	"github.com/charmingruby/doris/service/scribe/internal/quota/core/repository"
	"github.com/charmingruby/doris/service/scribe/internal/quota/core/usecase"
	"github.com/charmingruby/doris/service/scribe/internal/quota/delivery/event"
	"github.com/charmingruby/doris/service/scribe/internal/quota/delivery/http/rest/endpoint"
	"github.com/charmingruby/doris/service/scribe/internal/quota/persistence"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Datasource struct {
	quotaRepo      repository.QuotaRepository
	quotaUsageRepo repository.QuotaUsageRepository
	txManager      persistenceLib.TransactionManager[repository.TransactionManager]
}

func NewDatasource(db *sqlx.DB) (*Datasource, error) {
	quotaRepo, err := persistence.NewQuotaRepository(db)
	if err != nil {
		return nil, err
	}

	quotaUsageRepo, err := persistence.NewQuotaUsageRepository(db)
	if err != nil {
		return nil, err
	}

	txManager, err := persistence.NewTransactionManager(db)
	if err != nil {
		return nil, err
	}

	return &Datasource{
		quotaRepo:      quotaRepo,
		quotaUsageRepo: quotaUsageRepo,
		txManager:      txManager,
	}, nil
}

func NewUseCase(
	logger *instrumentation.Logger,
	datasource *Datasource,
) *usecase.UseCase {
	return usecase.New(
		logger,
		datasource.quotaRepo,
		datasource.quotaUsageRepo,
		datasource.txManager,
	)
}

func NewEventHandler(logger *instrumentation.Logger, sub *nats.Subscriber, cfg config.Config, uc *usecase.UseCase) {
	evtHandler := event.NewHandler(logger, sub, uc, event.TopicInput{
		APIKeyDelegated: cfg.Custom.APIKeyDelegatedTopic,
		APIKeyActivated: cfg.Custom.APIKeyActivatedTopic,
	})

	evtHandler.Subscribe()
}

func NewHTTPHandler(logger *instrumentation.Logger, r *gin.Engine, mw *rest.Middleware, val *validation.Validator, uc *usecase.UseCase) {
	endpoint.New(logger, r, mw, val, uc).Register()
}
