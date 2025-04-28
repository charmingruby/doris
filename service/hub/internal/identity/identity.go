package identity

import (
	"github.com/charmingruby/doris/lib/delivery/messaging"
	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/lib/validation"
	"github.com/charmingruby/doris/service/hub/config"
	"github.com/charmingruby/doris/service/hub/internal/identity/core/repository"
	"github.com/charmingruby/doris/service/hub/internal/identity/core/service"
	"github.com/charmingruby/doris/service/hub/internal/identity/delivery/event"
	"github.com/charmingruby/doris/service/hub/internal/identity/delivery/http/rest/endpoint"
	"github.com/gin-gonic/gin"
)

func NewEventHandler(pub messaging.Publisher, cfg config.Config) *event.Handler {
	return event.NewHandler(pub, event.TopicInput{
		OTPNotification: cfg.Custom.SendOTPNotificationTopic,
	})
}

func NewService(logger *instrumentation.Logger, apiKeyRepo repository.APIKeyRepository, otpRepo repository.OTPRepository, eventHandler *event.Handler) *service.Service {
	return service.New(logger, apiKeyRepo, otpRepo, eventHandler)
}

func NewHTTPHandler(logger *instrumentation.Logger, r *gin.Engine, val *validation.Validator, svc *service.Service) {
	endpoint.New(logger, r, val, svc).Register()
}
