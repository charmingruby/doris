package endpoint

import (
	"context"

	"github.com/charmingruby/doris/lib/delivery/http/rest"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/usecase"
	"github.com/gin-gonic/gin"
)

type CreateCodexRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=255"`
	Description string `json:"description"`
}

func (e *Endpoint) makeCreateCodex(c *gin.Context) {
	apiKeyID := c.GetString("api-key-id")
	if apiKeyID == "" {
		rest.NewUnauthorizedResponse(c)
		return
	}

	var req CreateCodexRequest
	if err := c.BindJSON(&req); err != nil {
		reasons := e.val.UnwrapValidationErr(err)

		rest.NewPayloadErrResponse(c, reasons)
		return
	}

	id, err := e.uc.CreateCodex(context.Background(), usecase.CreateCodexInput{
		CorrelationID: apiKeyID,
		Name:          req.Name,
		Description:   req.Description,
	})
	if err != nil {
		rest.HandleHTTPError(c, e.logger, err)
		return
	}

	rest.NewCreatedResponse(c, "", id, "codex")
}
