package endpoint

import (
	"context"

	"github.com/charmingruby/doris/lib/delivery/http/rest"
	"github.com/charmingruby/doris/service/account/internal/access/core/service"
	"github.com/gin-gonic/gin"
)

type SignInIntentRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (e *Endpoint) makeSignInIntent(c *gin.Context) {
	var req SignInIntentRequest
	if err := c.BindJSON(&req); err != nil {
		reasons := e.val.UnwrapValidationErr(err)

		rest.NewPayloadErrResponse(c, reasons)
		return
	}

	err := e.svc.SignInIntent(context.Background(), service.SignInIntentInput{
		Email: req.Email,
	})
	if err != nil {
		rest.HandleHTTPError(c, e.logger, err)
		return
	}

	rest.NewOKResponse(c, "confirmation code sent", nil)
}
