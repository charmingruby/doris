package notification

import (
	"github.com/charmingruby/doris/lib/delivery/http/rest"
	"github.com/charmingruby/doris/lib/delivery/messaging/nats"
	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/lib/persistence/dynamo"
	"github.com/charmingruby/doris/lib/validation"
	"github.com/charmingruby/doris/service/notification/config"
	"github.com/charmingruby/doris/service/notification/internal/notification/core/client"
	"github.com/charmingruby/doris/service/notification/internal/notification/core/repository"
	"github.com/charmingruby/doris/service/notification/internal/notification/core/service"
	"github.com/charmingruby/doris/service/notification/internal/notification/delivery/event"
	endpoint "github.com/charmingruby/doris/service/notification/internal/notification/delivery/http"
	"github.com/charmingruby/doris/service/notification/internal/notification/persistence"
	"github.com/gin-gonic/gin"
)

type Datasource struct {
	notificationRepo repository.NotificationRepository
}

func NewDatasource(cfg config.Config, db *dynamo.Client) (*Datasource, error) {
	notificationRepo := persistence.NewNotificationRepository(db.Conn, persistence.NotificationRepositoryInput{
		TableName:          cfg.Custom.NotificatiosnDynamoTable,
		CorrelationIDIndex: cfg.Custom.CorrelationIDDynamoIndex,
	})

	return &Datasource{
		notificationRepo: notificationRepo,
	}, nil
}

func NewService(logger *instrumentation.Logger, datasource *Datasource, notifier client.Notifier) *service.Service {
	return service.New(logger, datasource.notificationRepo, notifier)
}

func NewEventHandler(logger *instrumentation.Logger, sub *nats.Subscriber, cfg config.Config, svc *service.Service) {
	evtHandler := event.NewHandler(logger, sub, svc, event.TopicInput{
		SendOTPNotification: cfg.Custom.SendOTPNotificationTopic,
	})

	evtHandler.Subscribe()
}

func NewHTTPHandler(logger *instrumentation.Logger, r *gin.Engine, mw *rest.Middleware, val *validation.Validator, svc *service.Service) {
	endpoint.New(logger, r, mw, val, svc).Register()
}
