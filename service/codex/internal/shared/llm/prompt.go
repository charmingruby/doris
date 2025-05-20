package llm

import (
	"fmt"
	"strings"

	"github.com/charmingruby/doris/service/codex/internal/codex/core/model"
)

const (
	LIMIT_USER      = "<|user|>"
	LIMIT_SYSTEM    = "<|system|>"
	LIMIT_ASSISTANT = "<|assistant|>"

	CONTEXT_HEADER  = "Context:\n"
	CHUNK_FORMAT    = "%s\n"
	CONTEXT_INTRO   = "Context:\n"
	QUESTION_PREFIX = "Question: "
	ASSISTANT_INTRO = "Answer:"

	DEFAULT_SYSTEM_PROMPT = `You are a strict answer generator. Follow these rules EXACTLY:

1. FORMAT:
   - Use bullet points for multiple items
   - Maximum 3 sentences per point
   - Maximum 5 bullet points total

2. CONTENT RULES:
   - ONLY use information from the provided context
   - NO explanations or introductions
   - NO assumptions or inferences
   - NO personal opinions or suggestions
   - If information is not in context, respond with "No relevant information found in context"

3. LANGUAGE:
   - Use simple, direct language
   - Avoid technical jargon unless present in context
   - No conversational language
   - No questions or uncertainties

4. STRUCTURE:
   - One clear answer per bullet point
   - No nested information
   - No cross-references between points
   - No conclusions or summaries`
)

type PromptInput struct {
	SystemPrompt string
	UserPrompt   string
	Chunks       []model.CodexDocumentChunk
}

func BuildPrompt(in PromptInput) string {
	context := buildContextSection(in.Chunks)

	systemPrompt := in.SystemPrompt
	if systemPrompt == "" {
		systemPrompt = DEFAULT_SYSTEM_PROMPT
	}

	return fmt.Sprintf(`%s
%s
%s

%s
%s
%s
%s%s
%s

%s
%s
%s
`,
		LIMIT_SYSTEM, systemPrompt, LIMIT_SYSTEM,
		LIMIT_USER, CONTEXT_INTRO, context,
		QUESTION_PREFIX, in.UserPrompt, LIMIT_USER,
		LIMIT_ASSISTANT, ASSISTANT_INTRO, LIMIT_ASSISTANT,
	)
}

func buildContextSection(chunks []model.CodexDocumentChunk) string {
	if len(chunks) == 0 {
		return ""
	}

	var builder strings.Builder
	builder.WriteString(CONTEXT_HEADER)

	for _, chunk := range chunks {
		builder.WriteString(fmt.Sprintf(CHUNK_FORMAT, chunk.Content))
	}

	return builder.String()
}
