package notification

import (
	"github.com/charmingruby/doris/lib/delivery/messaging/nats"
	"github.com/charmingruby/doris/lib/instrumentation/logger"
	"github.com/charmingruby/doris/service/notification/config"
	"github.com/charmingruby/doris/service/notification/internal/notification/delivery/event"
)

func NewEventHandler(log *logger.Logger, sub *nats.Subscriber, cfg config.Config) {
	evtHandler := event.NewHandler(log, sub, event.HandlerInput{
		ReceiveNotificationTopic: cfg.Custom.NotificationsSendTopic,
	})

	evtHandler.Subscribe()
}
