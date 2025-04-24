package endpoint

import (
	"context"
	"net/http"
	"time"

	"github.com/charmingruby/doris/service/gateway/internal/identity/core/service"
	"github.com/gin-gonic/gin"
)

type GenerateAPIKeyRequest struct {
	FirstName string `json:"first_name" binding:"required,min=2,max=255"`
	LastName  string `json:"last_name" binding:"required,min=1,max=255"`
	Email     string `json:"email" binding:"required,email"`
}

func (e *Endpoint) makeGenerateAPIKey(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req GenerateAPIKeyRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": e.val.UnwrapValidationErr(err),
		})

		return
	}

	if err := e.svc.GenerateAPIKey(ctx, service.GenerateAPIKeyInput{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
