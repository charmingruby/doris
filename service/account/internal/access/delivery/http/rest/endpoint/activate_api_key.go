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

type ActivateAPIKeyRequest struct {
	OTP string `json:"otp" binding:"required,min=1,max=255"`
}

type ActivateAPIKeyResponse struct {
	AccessToken string `json:"access_token"`
}

func (e *Endpoint) makeActivateAPIKey(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

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

	token, err := e.svc.ActivateAPIKey(ctx, service.ActivateAPIKeyInput{
		APIKeyID: apiKeyID,
		OTP:      req.OTP,
	})

	if err != nil {
		var errResourceNotFound *custom_err.ErrResourceNotFound
		if errors.As(err, &errResourceNotFound) {
			rest.NewResourceNotFoundResponse(c, "api key")
			return
		}

		var errInvalidOTPCode *custom_err.ErrInvalidOTPCode
		if errors.As(err, &errInvalidOTPCode) {
			rest.NewConflictResponse(c, errInvalidOTPCode.Error())
			return
		}

		var errAPIKeyAlreadyActivated *custom_err.ErrAPIKeyAlreadyActivated
		if errors.As(err, &errAPIKeyAlreadyActivated) {
			rest.NewConflictResponse(c, errAPIKeyAlreadyActivated.Error())
			return
		}

		e.logger.Error("error on generate api key", "error", err)

		rest.NewUncaughtErrResponse(c, err)
		return
	}

	rest.NewOKResponse(c, "api key activated successfully", ActivateAPIKeyResponse{
		AccessToken: token,
	})
}
