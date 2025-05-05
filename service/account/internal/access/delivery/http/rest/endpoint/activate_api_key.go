package endpoint

import (
	"context"

	"github.com/charmingruby/doris/lib/delivery/http/rest"
	"github.com/charmingruby/doris/service/account/internal/access/core/service"
	"github.com/gin-gonic/gin"
)

type ActivateAPIKeyRequest struct {
	OTP string `json:"otp" binding:"required,min=1,max=255"`
}

type ActivateAPIKeyResponse struct {
	AccessToken string `json:"access_token"`
}

func (e *Endpoint) makeActivateAPIKey(c *gin.Context) {
	var req ActivateAPIKeyRequest
	if err := c.BindJSON(&req); err != nil {
		reasons := e.val.UnwrapValidationErr(err)

		rest.NewPayloadErrResponse(c, reasons)
		return
	}

	apiKeyID := c.Param("id")
	if apiKeyID == "" {
		rest.NewPayloadErrResponse(c, []string{"api key id is required"})
		return
	}

	token, err := e.svc.ActivateAPIKey(context.Background(), service.ActivateAPIKeyInput{
		APIKeyID: apiKeyID,
		OTP:      req.OTP,
	})

	if err != nil {
		rest.HandleHTTPError(c, e.logger, err)
		return
	}

	rest.NewOKResponse(c, "api key activated successfully", ActivateAPIKeyResponse{
		AccessToken: token,
	})
}
