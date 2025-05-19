package endpoint

import (
	"github.com/charmingruby/doris/lib/delivery/http/rest"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/usecase"

	"github.com/gin-gonic/gin"
)

func (e *Endpoint) makeUploadCodexDocuments(c *gin.Context) {
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

	files, err := rest.HandleMultipartFormFiles(c, e.logger, rest.HandleMultipartFormFilesInput{
		MaxFileSize:  5 * 1024 * 1024, // 5MB in bytes
		AllowedTypes: ".md",
		MaxFiles:     5,
	})
	if err != nil {
		rest.NewPayloadErrResponse(c, []string{err.Error()})
		return
	}

	op, err := e.uc.UploadDocuments(c, usecase.UploadDocumentsInput{
		CodexID:       codexID,
		CorrelationID: apiKeyID,
		Documents:     files.ValidFiles,
	})
	if err != nil {
		rest.HandleHTTPError(c, e.logger, err)
		return
	}

	rest.NewAcceptedResponse(c, []rest.AcceptedResponseIdentifier{
		{
			Key:   "uploaded",
			Value: op.UploadedDocs,
		},
		{
			Key:   "failed",
			Value: op.FailedDocs,
		},
	})
}
