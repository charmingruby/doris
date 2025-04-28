package notification

import (
	"github.com/charmingruby/doris/lib/delivery/messaging/nats"
	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/service/notification/config"
	"github.com/charmingruby/doris/service/notification/internal/notification/core/repository"
	"github.com/charmingruby/doris/service/notification/internal/notification/core/service"
	"github.com/charmingruby/doris/service/notification/internal/notification/delivery/event"
)

func NewService(logger *instrumentation.Logger, notificationRepo repository.NotificationRepository) *service.Service {
	return service.New(logger, notificationRepo)
}

func NewEventHandler(logger *instrumentation.Logger, sub *nats.Subscriber, cfg config.Config, svc *service.Service) {
	evtHandler := event.NewHandler(logger, sub, svc, event.TopicInput{
		OTPNotification: cfg.Custom.SendOTPNotificationTopic,
	})

	evtHandler.Subscribe()
}
