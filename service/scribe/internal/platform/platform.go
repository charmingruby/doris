package platform

import (
	"github.com/charmingruby/doris/service/scribe/internal/platform/delivery/http/rest/endpoint"
	"github.com/gin-gonic/gin"
)

func NewHTTPHandler(r *gin.Engine) {
	endpoint := endpoint.New(r)
	endpoint.Register()
}
