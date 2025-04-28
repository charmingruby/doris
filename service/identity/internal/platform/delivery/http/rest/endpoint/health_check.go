package endpoint

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (e *Endpoint) makeHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
