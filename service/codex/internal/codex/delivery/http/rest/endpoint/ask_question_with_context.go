package endpoint

import (
	"github.com/charmingruby/doris/lib/delivery/http/rest"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/usecase"

	"github.com/gin-gonic/gin"
)

type AskQuestionWithContextRequest struct {
	Question string `json:"question" binding:"required,min=1,max=255"`
}

func (e *Endpoint) makeAskQuestionWithContext(c *gin.Context) {
	codexID := c.Param("id")
	if codexID == "" {
		rest.NewPayloadErrResponse(c, []string{"codex id is required"})
		return
	}

	apiKeyID := c.GetString("api-key-id")
	if apiKeyID == "" {
		rest.NewUnauthorizedResponse(c)
		return
	}

	var req AskQuestionWithContextRequest
	if err := c.BindJSON(&req); err != nil {
		reasons := e.val.UnwrapValidationErr(err)

		rest.NewPayloadErrResponse(c, reasons)
		return
	}

	op, err := e.uc.AskQuestionWithContext(c, usecase.AskQuestionWithContextInput{
		CodexID:       codexID,
		CorrelationID: apiKeyID,
		Question:      req.Question,
	})
	if err != nil {
		rest.HandleHTTPError(c, e.logger, err)
		return
	}

	rest.NewOKResponse(c, "", op)
}
