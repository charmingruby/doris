package access

import (
	"github.com/charmingruby/doris/lib/delivery/messaging"
	"github.com/charmingruby/doris/lib/instrumentation"
	persistenceLib "github.com/charmingruby/doris/lib/persistence"
	"github.com/charmingruby/doris/lib/validation"
	"github.com/charmingruby/doris/service/account/config"
	"github.com/charmingruby/doris/service/account/internal/access/core/repository"
	"github.com/charmingruby/doris/service/account/internal/access/core/service"
	"github.com/charmingruby/doris/service/account/internal/access/delivery/event"
	"github.com/charmingruby/doris/service/account/internal/access/delivery/http/rest/endpoint"
	"github.com/charmingruby/doris/service/account/internal/access/persistence"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Datasource struct {
	apiKeyRepo repository.APIKeyRepository
	otpRepo    repository.OTPRepository
	txManager  persistenceLib.TransactionManager[repository.TransactionManager]
}

func NewDatasource(db *sqlx.DB) (*Datasource, error) {
	apiKeyRepo, err := persistence.NewAPIKeyRepo(db)
	if err != nil {
		return nil, err
	}

	otpRepo, err := persistence.NewOTPRepo(db)
	if err != nil {
		return nil, err
	}

	accessTxManager, err := persistence.NewTransactionManager(db)
	if err != nil {
		return nil, err
	}

	return &Datasource{
		apiKeyRepo: apiKeyRepo,
		otpRepo:    otpRepo,
		txManager:  accessTxManager,
	}, nil
}

func NewEventHandler(pub messaging.Publisher, cfg config.Config) *event.Handler {
	return event.NewHandler(pub, event.TopicInput{
		OTPNotification: cfg.Custom.SendOTPNotificationTopic,
	})
}

func NewService(
	logger *instrumentation.Logger,
	datasource *Datasource,
	eventHandler *event.Handler,
) *service.Service {
	return service.New(logger, datasource.apiKeyRepo, datasource.otpRepo, datasource.txManager, eventHandler)
}

func NewHTTPHandler(logger *instrumentation.Logger, r *gin.Engine, val *validation.Validator, svc *service.Service) {
	endpoint.New(logger, r, val, svc).Register()
}
