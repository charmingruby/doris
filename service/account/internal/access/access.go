package access

import (
	"github.com/charmingruby/doris/lib/delivery/http/rest"
	"github.com/charmingruby/doris/lib/delivery/messaging"
	"github.com/charmingruby/doris/lib/instrumentation"
	persistenceLib "github.com/charmingruby/doris/lib/persistence"
	"github.com/charmingruby/doris/lib/security"
	"github.com/charmingruby/doris/lib/validation"
	"github.com/charmingruby/doris/service/account/config"
	"github.com/charmingruby/doris/service/account/internal/access/core/repository"
	"github.com/charmingruby/doris/service/account/internal/access/core/usecase"
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
	apiKeyRepo, err := persistence.NewAPIKeyRepository(db)
	if err != nil {
		return nil, err
	}

	otpRepo, err := persistence.NewOTPRepository(db)
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

func NewEventHandler(logger *instrumentation.Logger, pub messaging.Publisher, cfg config.Config) *event.Handler {
	return event.NewHandler(logger, pub, event.TopicInput{
		SendOTPNotification: cfg.Custom.SendOTPNotificationTopic,
		APIKeyDelegated:     cfg.Custom.APIKeyDelegatedTopic,
		APIKeyActivated:     cfg.Custom.APIKeyActivatedTopic,
	})
}

func NewUseCase(
	logger *instrumentation.Logger,
	datasource *Datasource,
	eventHandler *event.Handler,
	tokenClient security.Token,
) *usecase.UseCase {
	return usecase.New(logger, datasource.apiKeyRepo, datasource.otpRepo, datasource.txManager, tokenClient, eventHandler)
}

func NewHTTPHandler(logger *instrumentation.Logger, r *gin.Engine, mw *rest.Middleware, val *validation.Validator, uc *usecase.UseCase) {
	endpoint.New(logger, r, mw, val, uc).Register()
}
