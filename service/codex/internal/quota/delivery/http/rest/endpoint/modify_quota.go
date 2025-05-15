package endpoint

import (
	"context"

	"github.com/charmingruby/doris/lib/delivery/http/rest"
	"github.com/charmingruby/doris/service/codex/internal/quota/core/model"
	"github.com/charmingruby/doris/service/codex/internal/quota/core/usecase"
	"github.com/gin-gonic/gin"
)

type ModifyQuotaRequest struct {
	Tier     string `json:"tier" binding:"max=255"`
	Kind     string `json:"kind" binding:"max=255"`
	MaxValue int    `json:"max_value"`
	Unit     string `json:"unit" binding:"max=255"`
	Status   string `json:"status" binding:"max=255"`
}

func (e *Endpoint) makeModifyQuota(c *gin.Context) {
	var req ModifyQuotaRequest
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

	if err := e.uc.ModifyQuota(context.Background(), usecase.ModifyQuotaInput{
		ID: quotaID,
		NewState: model.ModifyQuotaInput{
			Tier:     req.Tier,
			Status:   req.Status,
			Kind:     req.Kind,
			MaxValue: req.MaxValue,
			Unit:     req.Unit,
		},
	}); err != nil {
		rest.HandleHTTPError(c, e.logger, err)
		return
	}

	rest.NewNoContentResponse(c)
}
