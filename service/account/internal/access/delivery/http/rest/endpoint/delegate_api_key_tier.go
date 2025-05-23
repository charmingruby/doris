package endpoint

import (
	"context"

	"github.com/charmingruby/doris/lib/delivery/http/rest"
	"github.com/charmingruby/doris/service/account/internal/access/core/usecase"
	"github.com/gin-gonic/gin"
)

type DelegateAPIKeyTier struct {
	NewTier string `json:"new_tier" binding:"required,min=1,max=16"`
}

func (e *Endpoint) makeDelegateAPIKeyTier(c *gin.Context) {
	var req DelegateAPIKeyTier
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

	managerAPIKeyID := c.GetString("api-key-id")
	if managerAPIKeyID == "" {
		rest.NewUnauthorizedResponse(c)
		return
	}

	if err := e.uc.DelegateAPIKeyTier(context.Background(), usecase.DelegateAPIKeyTierInput{
		ManagerAPIKeyID:  managerAPIKeyID,
		APIKeyIDToChange: apiKeyID,
		NewTier:          req.NewTier,
	}); err != nil {
		rest.HandleHTTPError(c, e.logger, err)
		return
	}

	rest.NewOKResponse(c, "delegated api key new tier sucessfully", nil)
}
