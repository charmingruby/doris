package endpoint

import (
	"github.com/charmingruby/doris/lib/core/privilege"
	"github.com/charmingruby/doris/lib/delivery/http/rest"
	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/lib/validation"
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
	apiKey := api.Group("/api-keys")
	session := api.Group("/sessions")

	apiKey.POST("/generate", e.makeGenerateAPIKey)
	apiKey.PATCH("/:id/activate", e.makeActivateAPIKey)
	apiKey.POST("/:id/resend", e.makeResendAPIKeyActivation)
	apiKey.PATCH("/:id/delegate",
		e.mw.RBAC(
			privilege.TIER_MANAGER,
			privilege.TIER_ADMIN,
		),
		e.makeDelegateAPIKeyTier,
	)

	session.POST("/", e.makeSignInIntent)
	session.POST("/confirm", e.makeVerifySignInIntent)
}
