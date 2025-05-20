package usecase

import (
	"context"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/model"
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
	codex, err := u.codexRepo.FindByIDAndCorrelationID(ctx, in.CodexID, in.CorrelationID)

	if err != nil {
		return AskQuestionWithContextOutput{}, custom_err.NewErrDatasourceOperationFailed("find codex by id", err)
	}

	if codex.ID == "" {
		return AskQuestionWithContextOutput{}, custom_err.NewErrResourceNotFound("codex")
	}

	embedding, err := u.llm.GenerateEmbedding(ctx, in.Question)
	if err != nil {
		return AskQuestionWithContextOutput{}, custom_err.NewErrExternalService(err)
	}

	chunks, err := u.codexDocumentChunkRepo.FindSimilarChunks(ctx, codex.ID, embedding, 5)
	if err != nil {
		return AskQuestionWithContextOutput{}, custom_err.NewErrDatasourceOperationFailed("find similar chunks", err)
	}

	if len(chunks) == 0 {
		return AskQuestionWithContextOutput{
			Answer: "Couldn't find relevant information to answer your question.",
		}, nil
	}

	prompt := llm.BuildPrompt(llm.PromptInput{
		UserPrompt: in.Question,
		Chunks:     chunks,
	})

	answer, err := u.llm.GenerateCompletion(ctx, prompt, []string{
		llm.LIMIT_USER,
		llm.LIMIT_SYSTEM,
		llm.LIMIT_ASSISTANT,
	})
	if err != nil {
		return AskQuestionWithContextOutput{}, custom_err.NewErrExternalService(err)
	}

	qa := model.NewQA(model.QAInput{
		CodexID:  codex.ID,
		Question: in.Question,
		Answer:   answer,
	})

	if err := u.qaRepo.Create(ctx, *qa); err != nil {
		return AskQuestionWithContextOutput{}, custom_err.NewErrDatasourceOperationFailed("create qa", err)
	}

	return AskQuestionWithContextOutput{
		Answer: answer,
	}, nil
}
