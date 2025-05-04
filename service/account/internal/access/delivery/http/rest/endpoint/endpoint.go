package endpoint

import (
	"github.com/charmingruby/doris/lib/delivery/http/rest"
	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/lib/validation"
	"github.com/charmingruby/doris/service/account/internal/access/core/model"
	"github.com/charmingruby/doris/service/account/internal/access/core/service"
	"github.com/gin-gonic/gin"
)

type Endpoint struct {
	logger *instrumentation.Logger
	r      *gin.Engine
	mw     *rest.Middleware
	val    *validation.Validator
	svc    *service.Service
}

func New(
	logger *instrumentation.Logger,
	r *gin.Engine,
	mw *rest.Middleware,
	val *validation.Validator,
	svc *service.Service,
) *Endpoint {
	return &Endpoint{
		logger: logger,
		r:      r,
		mw:     mw,
		val:    val,
		svc:    svc,
	}
}

func (e *Endpoint) Register() {
	api := e.r.Group("/api")

	api.POST("/api-key/generate", e.makeGenerateAPIKey)
	api.POST("/api-key/:id/activate", e.makeActivateAPIKey)

	api.POST("/api-key/:id/delegate",
		e.mw.RBAC(
			model.API_KEY_TIER_MANAGER,
			model.API_KEY_TIER_ADMIN,
		),
		e.makeDelegateAPIKeyTier,
	)

	// POST signin
}
