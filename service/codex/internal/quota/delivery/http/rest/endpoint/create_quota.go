package endpoint

import (
	"context"

	"github.com/charmingruby/doris/lib/delivery/http/rest"
	"github.com/charmingruby/doris/service/codex/internal/quota/core/usecase"
	"github.com/gin-gonic/gin"
)

type CreateQuotaRequest struct {
	Tier     string `json:"tier" binding:"required,min=1,max=255"`
	Kind     string `json:"kind" binding:"required,min=1,max=255"`
	MaxValue int    `json:"max_value" binding:"required,min=1"`
	Unit     string `json:"unit" binding:"required,min=1,max=255"`
}

func (e *Endpoint) makeCreateQuota(c *gin.Context) {
	var req CreateQuotaRequest
	if err := c.BindJSON(&req); err != nil {
		reasons := e.val.UnwrapValidationErr(err)

		rest.NewPayloadErrResponse(c, reasons)
		return
	}

	id, err := e.uc.CreateQuota(context.Background(), usecase.CreateQuotaInput{
		Tier:     req.Tier,
		Kind:     req.Kind,
		MaxValue: req.MaxValue,
		Unit:     req.Unit,
	})
	if err != nil {
		rest.HandleHTTPError(c, e.logger, err)
		return
	}

	rest.NewCreatedResponse(c, "", id, "quota")
}
