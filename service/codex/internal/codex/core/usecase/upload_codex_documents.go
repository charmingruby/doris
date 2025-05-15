package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/charmingruby/doris/lib/fs"
)

type UploadCodexDocumentsInput struct {
	CodexID       string
	CorrelationID string
	Documents     []fs.File
}

func (u *UseCase) UploadCodexDocuments(ctx context.Context, in UploadCodexDocumentsInput) ([]string, []string, error) {
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

		// codexDocument := model.NewCodexDocument(model.CodexDocumentInput{
		// 	CodexID:  in.CodexID,
		// 	Title:    doc.Filename,
		// 	ImageURL: imageURL,
		// })

		fmt.Println(imageURL)

		uploadedDocs = append(uploadedDocs, key)
	}

	return uploadedDocs, failedDocs, nil
}

func (u *UseCase) codexDocumentKey(correlationID string, file fs.File) string {
	return fmt.Sprintf("%s/%d_%s", correlationID, time.Now().Unix(), file.Filename)
}
