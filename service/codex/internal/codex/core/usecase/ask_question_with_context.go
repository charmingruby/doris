package usecase

import (
	"context"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/model"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/repository"
	"github.com/charmingruby/doris/service/codex/internal/shared/core/kind"
	"github.com/charmingruby/doris/service/codex/internal/shared/llm"
)

type AskQuestionWithContextInput struct {
	Question      string
	CodexID       string
	CorrelationID string
}

type AskQuestionWithContextOutput struct {
	Answer string `json:"answer"`
}

func (u *UseCase) AskQuestionWithContext(ctx context.Context, in AskQuestionWithContextInput) (AskQuestionWithContextOutput, error) {
	availableQuota, err := u.quotaUsageManagementClient.CheckQuotaAvailability(ctx, in.CorrelationID, kind.QUOTA_LIMIT_PROMPT, 1)

	if err != nil {
		u.logger.Error("failed to check quota availability", "error", err, "correlation_id", in.CorrelationID)
		return AskQuestionWithContextOutput{}, err
	}

	if !availableQuota {
		return AskQuestionWithContextOutput{}, custom_err.NewErrQuotaExceeded()
	}

	codex, err := u.codexRepo.FindByIDAndCorrelationID(ctx, in.CodexID, in.CorrelationID)
	if err != nil {
		u.logger.Error("failed to find codex", "error", err, "codex_id", in.CodexID, "correlation_id", in.CorrelationID)
		return AskQuestionWithContextOutput{}, custom_err.NewErrDatasourceOperationFailed("find codex by id", err)
	}

	if codex.ID == "" {
		return AskQuestionWithContextOutput{}, custom_err.NewErrResourceNotFound("codex")
	}

	var answer string
	var qa *model.QA

	if err := u.txManager.Transact(func(tx repository.TransactionManager) error {
		embedding, err := u.llm.GenerateEmbedding(ctx, in.Question)
		if err != nil {
			u.logger.Error("failed to generate embedding", "error", err, "correlation_id", in.CorrelationID)
			return custom_err.NewErrExternalService(err)
		}

		chunks, err := tx.CodexDocumentChunkRepo.FindSimilarChunks(ctx, codex.ID, embedding, 5)
		if err != nil {
			u.logger.Error("failed to find similar chunks", "error", err, "codex_id", in.CodexID)
			return custom_err.NewErrDatasourceOperationFailed("find similar chunks", err)
		}

		if len(chunks) == 0 {
			answer = "Couldn't find relevant information to answer your question."
			return nil
		}

		prompt := llm.BuildPrompt(llm.PromptInput{
			UserPrompt: in.Question,
			Chunks:     chunks,
		})

		answer, err = u.llm.GenerateCompletion(ctx, prompt, []string{
			llm.LIMIT_USER,
			llm.LIMIT_SYSTEM,
			llm.LIMIT_ASSISTANT,
		})
		if err != nil {
			u.logger.Error("failed to generate completion", "error", err, "correlation_id", in.CorrelationID)
			return custom_err.NewErrExternalService(err)
		}

		qa = model.NewQA(model.QAInput{
			CodexID:  codex.ID,
			Question: in.Question,
			Answer:   answer,
		})

		if err := tx.QARepo.Create(ctx, *qa); err != nil {
			u.logger.Error("failed to create qa", "error", err, "codex_id", in.CodexID)
			return custom_err.NewErrDatasourceOperationFailed("create qa", err)
		}

		return nil
	}); err != nil {
		return AskQuestionWithContextOutput{}, err
	}

	if err := u.quotaUsageManagementClient.ConsumeQuota(ctx, in.CorrelationID, kind.QUOTA_LIMIT_PROMPT, 1); err != nil {
		u.logger.Error("failed to consume quota", "error", err, "correlation_id", in.CorrelationID)
	}

	return AskQuestionWithContextOutput{
		Answer: qa.Answer,
	}, nil
}
