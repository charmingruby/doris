package endpoint

import (
	"github.com/charmingruby/doris/lib/delivery/http/rest"
	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/lib/validation"
	"github.com/gin-gonic/gin"
)

type Endpoint struct {
	logger *instrumentation.Logger
	r      *gin.Engine
	mw     *rest.Middleware
	val    *validation.Validator
}

func New(
	logger *instrumentation.Logger,
	r *gin.Engine,
	mw *rest.Middleware,
	val *validation.Validator,
) *Endpoint {
	return &Endpoint{
		logger: logger,
		r:      r,
		mw:     mw,
		val:    val,
	}
}

func (e *Endpoint) Register() {
	// protected := e.r.Group("/api", e.mw.Auth())
	// protected.POST("/codex/", e.makeCreateCodex)
	// protected.POST("/codex/:id/documents/upload", e.makeUploadDocuments)
	// protected.POST("/codex/:id/chat", e.makeCodexChat)
}
