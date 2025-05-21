package usecase

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/fs"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/event"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/model"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/repository"
	"github.com/charmingruby/doris/service/codex/internal/shared/core/kind"
)

type UploadDocumentsInput struct {
	CodexID       string
	CorrelationID string
	Documents     []fs.File
}

type UploadDocumentsOutput struct {
	UploadedDocs []string
	FailedDocs   []string
}

type documentProcessingResult struct {
	docID    string
	filename string
	error    error
}

type documentProcessingJob struct {
	doc           fs.File
	codexID       string
	correlationID string
}

func (u *UseCase) UploadDocuments(ctx context.Context, in UploadDocumentsInput) (UploadDocumentsOutput, error) {
	docsCount := len(in.Documents)

	availableQuota, err := u.quotaUsageManagementClient.CheckQuotaAvailability(ctx, in.CorrelationID, kind.QUOTA_LIMIT_DOCUMENT, docsCount)
	if err != nil {
		return UploadDocumentsOutput{}, err
	}

	if !availableQuota {
		return UploadDocumentsOutput{}, custom_err.NewErrQuotaExceeded()
	}

	codex, err := u.codexRepo.FindByIDAndCorrelationID(ctx, in.CodexID, in.CorrelationID)
	if err != nil {
		return UploadDocumentsOutput{}, custom_err.NewErrDatasourceOperationFailed("find codex by id", err)
	}

	if codex.CorrelationID != in.CorrelationID {
		return UploadDocumentsOutput{}, custom_err.NewErrResourceNotFound("codex")
	}

	jobs := make(chan documentProcessingJob, docsCount)
	results := make(chan documentProcessingResult, docsCount)

	numWorkers := min(docsCount, 5)

	var wg sync.WaitGroup
	for range numWorkers {
		wg.Add(1)
		go u.documentProcessor(ctx, &wg, jobs, results)
	}

	go func() {
		for _, doc := range in.Documents {
			jobs <- documentProcessingJob{
				doc:           doc,
				codexID:       in.CodexID,
				correlationID: in.CorrelationID,
			}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	uploadedDocs := []string{}
	failedDocs := []string{}

	for result := range results {
		if result.error != nil {
			failedDocs = append(failedDocs, result.filename)
			continue
		}

		uploadedDocs = append(uploadedDocs, result.docID)
	}

	if len(uploadedDocs) > 0 {
		if err := u.quotaUsageManagementClient.ConsumeQuota(ctx, in.CorrelationID, kind.QUOTA_LIMIT_DOCUMENT, len(uploadedDocs)); err != nil {
			u.logger.Error("failed to consume quota", "error", err)
		}
	}

	return UploadDocumentsOutput{
		UploadedDocs: uploadedDocs,
		FailedDocs:   failedDocs,
	}, nil
}

func (u *UseCase) documentProcessor(
	ctx context.Context,
	wg *sync.WaitGroup,
	jobs <-chan documentProcessingJob,
	results chan<- documentProcessingResult,
) {
	defer wg.Done()

	for job := range jobs {
		key := u.codexDocumentKey(job.correlationID, job.doc)

		imageURL, err := u.storage.Upload(
			ctx,
			u.embeddingSourceDocsBucket,
			key,
			job.doc.File,
		)
		if err != nil {
			u.logger.Error("failed to upload file", "error", err)
			results <- documentProcessingResult{
				filename: job.doc.Filename,
				error:    err,
			}

			continue
		}

		if err := u.txManager.Transact(func(txManager repository.TransactionManager) error {
			codexDocument := model.NewCodexDocument(model.CodexDocumentInput{
				CodexID:  job.codexID,
				Title:    job.doc.Filename,
				ImageURL: imageURL,
			})

			if err := txManager.CodexDocumentRepo.Create(ctx, *codexDocument); err != nil {
				u.logger.Error("failed to create codex document", "error", err)
				return err
			}

			event := event.CodexDocumentUploaded{
				ID:            codexDocument.ID,
				CodexID:       codexDocument.CodexID,
				CorrelationID: job.correlationID,
				ImageURL:      codexDocument.ImageURL,
				SentAt:        time.Now(),
			}

			if err := u.eventHandler.DispatchCodexDocumentUploaded(ctx, event); err != nil {
				u.logger.Error("failed to dispatch codex document uploaded event", "error", err)
				return err
			}

			results <- documentProcessingResult{
				docID:    codexDocument.ID,
				filename: job.doc.Filename,
			}
			return nil
		}); err != nil {
			results <- documentProcessingResult{
				filename: job.doc.Filename,
				error:    err,
			}
		}
	}
}

func (u *UseCase) codexDocumentKey(correlationID string, file fs.File) string {
	return fmt.Sprintf("%s/%d_%s", correlationID, time.Now().Unix(), file.Filename)
}
