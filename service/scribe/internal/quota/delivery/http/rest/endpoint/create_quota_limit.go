package endpoint

import (
	"context"

	"github.com/charmingruby/doris/lib/delivery/http/rest"
	"github.com/charmingruby/doris/service/scribe/internal/quota/core/usecase"
	"github.com/gin-gonic/gin"
)

type CreateQuotaLimitRequest struct {
	Kind     string `json:"kind" binding:"required,min=1,max=255"`
	MaxValue int    `json:"max_value" binding:"required,min=1"`
	Unit     string `json:"unit" binding:"required,min=1,max=255"`
}

func (e *Endpoint) makeCreateQuotaLimit(c *gin.Context) {
	var req CreateQuotaLimitRequest
	if err := c.BindJSON(&req); err != nil {
		reasons := e.val.UnwrapValidationErr(err)

		rest.NewPayloadErrResponse(c, reasons)
		return
	}

	quotaID := c.Param("id")
	if quotaID == "" {
		rest.NewPayloadErrResponse(c, []string{"quota id is required"})
		return
	}

	id, err := e.uc.CreateQuotaLimit(context.Background(), usecase.CreateQuotaLimitInput{
		QuotaID:  quotaID,
		Kind:     req.Kind,
		MaxValue: req.MaxValue,
		Unit:     req.Unit,
	})
	if err != nil {
		rest.HandleHTTPError(c, e.logger, err)
		return
	}

	rest.NewCreatedResponse(c, "", id, "quota limit")
}
