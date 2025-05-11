package endpoint

import (
	"context"

	"github.com/charmingruby/doris/lib/delivery/http/rest"
	"github.com/charmingruby/doris/service/scribe/internal/quota/core/model"
	"github.com/charmingruby/doris/service/scribe/internal/quota/core/usecase"
	"github.com/gin-gonic/gin"
)

type ModifyQuotaLimitRequest struct {
	Kind     string `json:"kind" binding:"max=255"`
	MaxValue int    `json:"max_value"`
	Unit     string `json:"unit" binding:"max=255"`
	IsActive bool   `json:"is_active"`
}

func (e *Endpoint) makeModifyQuotaLimit(c *gin.Context) {
	var req ModifyQuotaLimitRequest
	if err := c.BindJSON(&req); err != nil {
		reasons := e.val.UnwrapValidationErr(err)

		rest.NewPayloadErrResponse(c, reasons)
		return
	}

	quotaLimitID := c.Param("limit_id")
	if quotaLimitID == "" {
		rest.NewPayloadErrResponse(c, []string{"quota limit id is required"})
		return
	}

	if err := e.uc.ModifyQuotaLimit(context.Background(), usecase.ModifyQuotaLimitInput{
		ID: quotaLimitID,
		NewState: model.ModifyQuotaLimitInput{
			Kind:     req.Kind,
			MaxValue: req.MaxValue,
			Unit:     req.Unit,
			IsActive: req.IsActive,
		},
	}); err != nil {
		rest.HandleHTTPError(c, e.logger, err)
		return
	}

	rest.NewNoContentResponse(c)
}
