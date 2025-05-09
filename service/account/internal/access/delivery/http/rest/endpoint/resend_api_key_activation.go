package endpoint

import (
	"context"

	"github.com/charmingruby/doris/lib/delivery/http/rest"
	"github.com/charmingruby/doris/service/account/internal/access/core/usecase"
	"github.com/gin-gonic/gin"
)

func (e *Endpoint) makeResendAPIKeyActivation(c *gin.Context) {
	apiKeyID := c.Param("id")
	if apiKeyID == "" {
		rest.NewPayloadErrResponse(c, []string{"api key id is required"})
		return
	}

	if err := e.uc.ResendAPIKeyActivation(context.Background(), usecase.ResendAPIKeyActivationInput{
		APIKeyID: apiKeyID,
	}); err != nil {
		rest.HandleHTTPError(c, e.logger, err)
		return
	}

	rest.NewOKResponse(c, "confirmation code sent", nil)
}
