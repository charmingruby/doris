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

type SignInIntentRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (e *Endpoint) makeSignInIntent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req SignInIntentRequest
	if err := c.BindJSON(&req); err != nil {
		reasons := e.val.UnwrapValidationErr(err)

		rest.NewPayloadErrResponse(c, reasons)
		return
	}

	err := e.svc.SignInIntent(ctx, service.SignInIntentInput{
		Email: req.Email,
	})
	if err != nil {
		var errInvalidEntity *custom_err.ErrInvalidEntity
		if errors.As(err, &errInvalidEntity) {
			rest.NewUnprocessableEntityResponse(c, errInvalidEntity.Error())
			return
		}

		var errResourceNotFound *custom_err.ErrResourceNotFound
		if errors.As(err, &errResourceNotFound) {
			rest.NewResourceNotFoundResponse(c, "api key")
			return
		}

		e.logger.Error("error on generate api key", "error", err)

		rest.NewUncaughtErrResponse(c, err)
		return
	}

	rest.NewOKResponse(c, "confirmation code sent", nil)
}
