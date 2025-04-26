package notification

import (
	"github.com/charmingruby/doris/lib/delivery/messaging/nats"
	"github.com/charmingruby/doris/lib/instrumentation/logger"
	"github.com/charmingruby/doris/service/notification/config"
	"github.com/charmingruby/doris/service/notification/internal/notification/core/service"
	"github.com/charmingruby/doris/service/notification/internal/notification/delivery/event"
)

func NewService(log *logger.Logger) *service.Service {
	return service.New(log)
}

func NewEventHandler(log *logger.Logger, sub *nats.Subscriber, cfg config.Config, svc *service.Service) {
	evtHandler := event.NewHandler(log, sub, svc, event.TopicInput{
		ReceiveNotificationTopic: cfg.Custom.NotificationsSendTopic,
	})

	evtHandler.Subscribe()
}
