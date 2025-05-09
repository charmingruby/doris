package endpoint

import (
	"context"

	"github.com/charmingruby/doris/lib/delivery/http/rest"
	"github.com/charmingruby/doris/service/account/internal/access/core/usecase"
	"github.com/gin-gonic/gin"
)

type GenerateAPIKeyRequest struct {
	FirstName string `json:"first_name" binding:"required,min=2,max=255"`
	LastName  string `json:"last_name" binding:"required,min=1,max=255"`
	Email     string `json:"email" binding:"required,email"`
}

func (e *Endpoint) makeGenerateAPIKey(c *gin.Context) {
	var req GenerateAPIKeyRequest
	if err := c.BindJSON(&req); err != nil {
		reasons := e.val.UnwrapValidationErr(err)

		rest.NewPayloadErrResponse(c, reasons)
		return
	}

	id, err := e.uc.GenerateAPIKey(context.Background(), usecase.GenerateAPIKeyInput{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	})
	if err != nil {
		rest.HandleHTTPError(c, e.logger, err)
		return
	}

	rest.NewCreatedResponse(c, "confirmation code sent", id, "api key")
}
