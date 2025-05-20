package external

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/tmc/langchaingo/textsplitter"
)

type Ollama struct {
	logger *instrumentation.Logger
	client *http.Client

	EmbeddingModel  string
	CompletionModel string
	BaseURL         string
}

type OllamaInput struct {
	EmbeddingModel  string
	CompletionModel string
	BaseURL         string
}

func NewOllama(logger *instrumentation.Logger, in OllamaInput) *Ollama {
	return &Ollama{
		logger:          logger,
		EmbeddingModel:  in.EmbeddingModel,
		CompletionModel: in.CompletionModel,
		BaseURL:         in.BaseURL,
		client:          &http.Client{},
	}
}

func (o *Ollama) ChunkText(text string) ([]string, error) {
	splitter := textsplitter.NewRecursiveCharacter(
		textsplitter.WithChunkSize(1000),
		textsplitter.WithChunkOverlap(200),
	)

	chunks, err := splitter.SplitText(text)
	if err != nil {
		return nil, err
	}

	return chunks, nil
}

type OllamaGenerateEmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type OllamaGenerateEmbeddingResponse struct {
	Embedding []float64 `json:"embedding"`
}

func (o *Ollama) GenerateEmbedding(ctx context.Context, text string) ([]float64, error) {
	req := OllamaGenerateEmbeddingRequest{
		Model:  o.EmbeddingModel,
		Prompt: text,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/embeddings", o.BaseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	res, err := o.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, err
	}

	var embeddingRes OllamaGenerateEmbeddingResponse
	if err := json.NewDecoder(res.Body).Decode(&embeddingRes); err != nil {
		return nil, err
	}

	return embeddingRes.Embedding, nil
}

type OllamaCompletionRequest struct {
	Model   string         `json:"model"`
	Prompt  string         `json:"prompt"`
	Stream  bool           `json:"stream"`
	Options map[string]any `json:"options,omitempty"`
}

type OllamaCompletionResponse struct {
	Model           string `json:"model"`
	CreatedAt       string `json:"created_at"`
	Response        string `json:"response"`
	Done            bool   `json:"done"`
	DoneReason      string `json:"done_reason"`
	Context         []int  `json:"context"`
	TotalDuration   int64  `json:"total_duration"`
	LoadDuration    int64  `json:"load_duration"`
	PromptEvalCount int    `json:"prompt_eval_count"`
	EvalCount       int    `json:"eval_count"`
}

func (o *Ollama) GenerateCompletion(
	ctx context.Context,
	prompt string,
	limits []string,
) (string, error) {
	req := OllamaCompletionRequest{
		Model:  o.CompletionModel,
		Prompt: prompt,
		Stream: false,
		Options: map[string]any{
			"temperature": 0.1,
			"top_p":       0.5,
			"stop":        limits,
		},
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/api/generate", o.BaseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return "", err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	res, err := o.client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", err
	}

	var completionRes OllamaCompletionResponse
	if err := json.NewDecoder(res.Body).Decode(&completionRes); err != nil {
		return "", err
	}

	return completionRes.Response, nil
}
