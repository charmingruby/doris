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
	Files         []fs.File
}

func (u *UseCase) UploadCodexDocuments(ctx context.Context, in UploadCodexDocumentsInput) ([]string, []string, error) {
	filesUploaded := []string{}
	filesFailed := []string{}

	for _, file := range in.Files {
		key := u.codexDocumentKey(in.CorrelationID, file)

		err := u.storage.Upload(
			ctx,
			u.embeddingSourceDocsBucket,
			key,
			file.File,
		)
		if err != nil {
			u.logger.Error("failed to upload file", "error", err)
			filesFailed = append(filesFailed, file.Filename)
			continue
		}

		filesUploaded = append(filesUploaded, key)
	}

	return filesUploaded, filesFailed, nil
}

func (u *UseCase) codexDocumentKey(correlationID string, file fs.File) string {
	return fmt.Sprintf("%s/%d_%s", correlationID, time.Now().Unix(), file.Filename)
}
