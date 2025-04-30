package notification

import (
	"github.com/charmingruby/doris/lib/delivery/messaging/nats"
	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/service/notification/config"
	"github.com/charmingruby/doris/service/notification/internal/notification/core/client"
	"github.com/charmingruby/doris/service/notification/internal/notification/core/repository"
	"github.com/charmingruby/doris/service/notification/internal/notification/core/service"
	"github.com/charmingruby/doris/service/notification/internal/notification/delivery/event"
	"github.com/charmingruby/doris/service/notification/test/memory"
)

type Datasource struct {
	notificationRepo repository.NotificationRepository
}

func NewDatasource() (*Datasource, error) {
	notificationRepo := memory.NewNotificationRepository()

	return &Datasource{
		notificationRepo: notificationRepo,
	}, nil
}

func NewService(logger *instrumentation.Logger, datasource *Datasource, notifier client.Notifier) *service.Service {
	return service.New(logger, datasource.notificationRepo, notifier)
}

func NewEventHandler(logger *instrumentation.Logger, sub *nats.Subscriber, cfg config.Config, svc *service.Service) {
	evtHandler := event.NewHandler(logger, sub, svc, event.TopicInput{
		OTPNotification: cfg.Custom.SendOTPNotificationTopic,
	})

	evtHandler.Subscribe()
}
