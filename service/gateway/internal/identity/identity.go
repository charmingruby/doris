package identity

import (
	"github.com/charmingruby/doris/lib/delivery/messaging"
	"github.com/charmingruby/doris/lib/instrumentation/logger"
	"github.com/charmingruby/doris/service/gateway/config"
	"github.com/charmingruby/doris/service/gateway/internal/identity/core/repository"
	"github.com/charmingruby/doris/service/gateway/internal/identity/core/service"
	"github.com/charmingruby/doris/service/gateway/internal/identity/delivery/event"
	"github.com/charmingruby/doris/service/gateway/internal/identity/delivery/http/rest/endpoint"
	"github.com/gin-gonic/gin"
)

func NewEventHandler(pub messaging.Publisher, cfg config.Config) *event.Handler {
	return event.NewHandler(pub, event.HandlerInput{
		APIKeyRequestTopic: cfg.Custom.NotificationsSendTopic,
	})
}

func NewService(log *logger.Logger, apiKeyRepo repository.APIKeyRepository, eventHandler *event.Handler) *service.Service {
	return service.New(log, apiKeyRepo, eventHandler)
}

func NewHTTPHandler(log *logger.Logger, r *gin.Engine, svc *service.Service) {
	endpoint.New(log, r, svc).Register()
}
