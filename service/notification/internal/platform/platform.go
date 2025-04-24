package platform

import (
	"github.com/charmingruby/doris/service/notification/internal/platform/delivery/http/rest/endpoint"
	"github.com/gin-gonic/gin"
)

func NewHTTPHandler(r *gin.Engine) {
	endpoint := endpoint.New(r)
	endpoint.Register()
}
