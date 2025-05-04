package endpoint

import (
	"context"
	"errors"
	"time"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/delivery/http/rest"
	"github.com/charmingruby/doris/service/account/internal/access/core/service"
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
		reasons := e.val.UnwrapValidationErr(err)

		rest.NewPayloadErrResponse(c, reasons)
		return
	}

	id, err := e.svc.GenerateAPIKey(ctx, service.GenerateAPIKeyInput{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	})
	if err != nil {
		var errResourceAlreadyExists *custom_err.ErrResourceAlreadyExists
		if errors.As(err, &errResourceAlreadyExists) {
			rest.NewResourceAlreadyExistsResponse(c, "api key")
			return
		}

		e.logger.Error("error on generate api key", "error", err)

		rest.NewUncaughtErrResponse(c, err)
		return
	}

	rest.NewCreatedResponse(c, "confirmation code sent", id, "api key")
}
