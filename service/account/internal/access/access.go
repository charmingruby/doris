package access

import (
	"github.com/charmingruby/doris/lib/delivery/messaging"
	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/lib/persistence"
	"github.com/charmingruby/doris/lib/validation"
	"github.com/charmingruby/doris/service/account/config"
	"github.com/charmingruby/doris/service/account/internal/access/core/repository"
	"github.com/charmingruby/doris/service/account/internal/access/core/service"
	"github.com/charmingruby/doris/service/account/internal/access/delivery/event"
	"github.com/charmingruby/doris/service/account/internal/access/delivery/http/rest/endpoint"
	"github.com/gin-gonic/gin"
)

func NewEventHandler(pub messaging.Publisher, cfg config.Config) *event.Handler {
	return event.NewHandler(pub, event.TopicInput{
		OTPNotification: cfg.Custom.SendOTPNotificationTopic,
	})
}

func NewService(
	logger *instrumentation.Logger,
	apiKeyRepo repository.APIKeyRepository,
	otpRepo repository.OTPRepository,
	eventHandler *event.Handler,
	txManager persistence.TransactionManager[repository.TransactionManager],
) *service.Service {
	return service.New(logger, apiKeyRepo, otpRepo, txManager, eventHandler)
}

func NewHTTPHandler(logger *instrumentation.Logger, r *gin.Engine, val *validation.Validator, svc *service.Service) {
	endpoint.New(logger, r, val, svc).Register()
}
