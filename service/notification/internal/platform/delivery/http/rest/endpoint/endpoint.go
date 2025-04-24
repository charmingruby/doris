package endpoint

import "github.com/gin-gonic/gin"

type Endpoint struct {
	r *gin.Engine
}

func New(r *gin.Engine) *Endpoint {
	return &Endpoint{
		r: r,
	}
}

func (e *Endpoint) Register() {
	api := e.r.Group("/api")

	api.GET("/health-check", e.makeHealthCheck)
}
