package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/fs"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/event"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/model"
)

type UploadCodexDocumentsInput struct {
	CodexID       string
	CorrelationID string
	Documents     []fs.File
}

type UploadCodexDocumentsOutput struct {
	UploadedDocs []string
	FailedDocs   []string
}

func (u *UseCase) UploadCodexDocuments(ctx context.Context, in UploadCodexDocumentsInput) (UploadCodexDocumentsOutput, error) {
	codex, err := u.codexRepo.FindByID(ctx, in.CodexID)

	if err != nil {
		return UploadCodexDocumentsOutput{}, custom_err.NewErrDatasourceOperationFailed("find codex by id", err)
	}

	if codex.CorrelationID != in.CorrelationID {
		return UploadCodexDocumentsOutput{}, custom_err.NewErrResourceNotFound("codex")
	}

	uploadedDocs := []string{}
	failedDocs := []string{}

	for _, doc := range in.Documents {
		key := u.codexDocumentKey(in.CorrelationID, doc)

		imageURL, err := u.storage.Upload(
			ctx,
			u.embeddingSourceDocsBucket,
			key,
			doc.File,
		)
		if err != nil {
			u.logger.Error("failed to upload file", "error", err)
			failedDocs = append(failedDocs, doc.Filename)
			continue
		}

		codexDocument := model.NewCodexDocument(model.CodexDocumentInput{
			CodexID:  in.CodexID,
			Title:    doc.Filename,
			ImageURL: imageURL,
		})

		if err := u.codexDocumentRepo.Create(ctx, *codexDocument); err != nil {
			u.logger.Error("failed to create codex document", "error", err)
			failedDocs = append(failedDocs, doc.Filename)
			continue
		}

		event := event.CodexDocumentUploaded{
			ID:            codexDocument.ID,
			CodexID:       codexDocument.CodexID,
			CorrelationID: in.CorrelationID,
			ImageURL:      codexDocument.ImageURL,
			SentAt:        time.Now(),
		}

		if err := u.eventHandler.DispatchCodexDocumentUploaded(ctx, event); err != nil {
			u.logger.Error("failed to dispatch codex document uploaded event", "error", err)
			failedDocs = append(failedDocs, doc.Filename)
			continue
		}

		uploadedDocs = append(uploadedDocs, codexDocument.ID)
	}

	return UploadCodexDocumentsOutput{
		UploadedDocs: uploadedDocs,
		FailedDocs:   failedDocs,
	}, nil
}

func (u *UseCase) codexDocumentKey(correlationID string, file fs.File) string {
	return fmt.Sprintf("%s/%d_%s", correlationID, time.Now().Unix(), file.Filename)
}
