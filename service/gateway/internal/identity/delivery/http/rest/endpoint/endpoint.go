package endpoint

import (
	"github.com/charmingruby/doris/lib/instrumentation/logger"
	"github.com/charmingruby/doris/service/gateway/internal/identity/core/service"
	"github.com/gin-gonic/gin"
)

type Endpoint struct {
	logger *logger.Logger
	r      *gin.Engine
	svc    *service.Service
}

func New(
	logger *logger.Logger,
	r *gin.Engine,
	svc *service.Service,
) *Endpoint {
	return &Endpoint{
		logger: logger,
		r:      r,
		svc:    svc,
	}
}

func (e *Endpoint) Register() {
	api := e.r.Group("/api")

	api.GET("/api-key/request", e.makeRequestAPIKey)
}
