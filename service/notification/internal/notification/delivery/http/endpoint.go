package endpoint

import (
	"github.com/charmingruby/doris/lib/delivery/http/rest"
	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/lib/validation"
	"github.com/charmingruby/doris/service/notification/internal/notification/core/usecase"
	"github.com/gin-gonic/gin"
)

type Endpoint struct {
	logger *instrumentation.Logger
	r      *gin.Engine
	mw     *rest.Middleware
	val    *validation.Validator
	uc     *usecase.UseCase
}

func New(
	logger *instrumentation.Logger,
	r *gin.Engine,
	mw *rest.Middleware,
	val *validation.Validator,
	uc *usecase.UseCase,
) *Endpoint {
	return &Endpoint{
		logger: logger,
		r:      r,
		mw:     mw,
		val:    val,
		uc:     uc,
	}
}

func (e *Endpoint) Register() {
	api := e.r.Group("/api")
	notifications := api.Group("/notifications")

	notifications.GET("/", e.mw.Auth(), e.makeFetchCorrelatedNotifications)
}
