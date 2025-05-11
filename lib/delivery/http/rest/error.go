package rest

import (
	"errors"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/gin-gonic/gin"
)

func HandleHTTPError(c *gin.Context, logger *instrumentation.Logger, err error) {
	switch e := err.(type) {
	case *custom_err.ErrResourceNotFound:
		NewResourceNotFoundResponse(c, e.Error())
	case *custom_err.ErrInvalidOTPCode:
		NewConflictResponse(c, e.Error())
	case *custom_err.ErrAPIKeyAlreadyActivated:
		NewConflictResponse(c, e.Error())
	case *custom_err.ErrNothingToChange:
		NewUnprocessableEntityResponse(c, e.Error())
	case *custom_err.ErrInsufficientPermission:
		NewForbiddenResponse(c)
	case *custom_err.ErrInvalidEntity:
		NewUnprocessableEntityResponse(c, e.Error())
	case *custom_err.ErrResourceAlreadyExists:
		NewResourceAlreadyExistsResponse(c, e.Error())
	case *custom_err.ErrOTPGenerationCooldown:
		NewTooManyRequestsResponse(c, err.Error())

	default:
		logger.Error("uncaught error",
			"error", err,
			"error_type", errors.Unwrap(err),
			"path", c.Request.URL.Path,
			"method", c.Request.Method,
		)
		NewUncaughtErrResponse(c, err)
	}
}
