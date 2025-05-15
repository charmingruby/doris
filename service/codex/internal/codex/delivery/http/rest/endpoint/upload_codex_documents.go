package endpoint

import (
	"github.com/charmingruby/doris/lib/delivery/http/rest"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/usecase"

	"github.com/gin-gonic/gin"
)

const (
	maxFileSize  = 5 * 1024 * 1024 // 5MB in bytes
	maxFiles     = 5
	allowedTypes = ".pdf,.txt"
	s3BucketName = "embeddings-source-docs" // Replace with your actual bucket name
)

type UploadResponse struct {
	Success bool     `json:"success"`
	Files   []string `json:"files"`
	Error   string   `json:"error,omitempty"`
}

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
		MaxFileSize:  maxFileSize,
		AllowedTypes: allowedTypes,
		MaxFiles:     maxFiles,
	})
	if err != nil {
		rest.NewPayloadErrResponse(c, []string{err.Error()})
		return
	}

	filesUploaded, filesFailed, err := e.uc.UploadCodexDocuments(c, usecase.UploadCodexDocumentsInput{
		CodexID:       codexID,
		CorrelationID: apiKeyID,
		Files:         files.ValidFiles,
	})
	if err != nil {
		rest.HandleHTTPError(c, e.logger, err)
		return
	}

	rest.NewAcceptedResponse(c, []rest.AcceptedResponseIdentifier{
		{
			Key:   "successfully_uploaded_files",
			Value: filesUploaded,
		},
		{
			Key:   "failed_to_upload_files",
			Value: filesFailed,
		},
	})
}
