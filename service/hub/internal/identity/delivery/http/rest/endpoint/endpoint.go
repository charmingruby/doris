package endpoint

import (
	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/lib/validation"
	"github.com/charmingruby/doris/service/hub/internal/identity/core/service"
	"github.com/gin-gonic/gin"
)

type Endpoint struct {
	logger *instrumentation.Logger
	r      *gin.Engine
	val    *validation.Validator
	svc    *service.Service
}

func New(
	logger *instrumentation.Logger,
	r *gin.Engine,
	val *validation.Validator,
	svc *service.Service,
) *Endpoint {
	return &Endpoint{
		logger: logger,
		r:      r,
		val:    val,
		svc:    svc,
	}
}

func (e *Endpoint) Register() {
	api := e.r.Group("/api")

	api.POST("/api-key/generate", e.makeGenerateAPIKey)
}
