package endpoint

import (
	"context"
	"time"

	"github.com/charmingruby/doris/lib/delivery/http/rest"
	"github.com/charmingruby/doris/service/account/internal/access/core/service"
	"github.com/gin-gonic/gin"
)

type VerifySignInIntentRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,min=6,max=6"`
}

type VerifySignInIntentResponse struct {
	AccessToken string `json:"access_token"`
}

func (e *Endpoint) makeVerifySignInIntent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req VerifySignInIntentRequest
	if err := c.BindJSON(&req); err != nil {
		reasons := e.val.UnwrapValidationErr(err)

		rest.NewPayloadErrResponse(c, reasons)
		return
	}

	token, err := e.svc.VerifySignInIntent(ctx, service.VerifySignInIntentInput{
		Email: req.Email,
		OTP:   req.OTP,
	})
	if err != nil {
		rest.HandleHTTPError(c, e.logger, err)
		return
	}

	rest.NewOKResponse(c, "signed in successfully", VerifySignInIntentResponse{
		AccessToken: token,
	})
}
