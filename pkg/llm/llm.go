package llm

import (
	"bytes"
	"context"
	"strings"
	"text/template"

	"github.com/sashabaranov/go-openai"
)

//go:generate mockgen -destination=../mocks/mock_llm.go -package=mocks . LLMClient,OpenAIClient
type OpenAIClient interface {
	CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}

type LLMClient interface {
	IdentifyTags(messageContent *TagInput, tagOptions []string) ([]string, error)
	RequestChatCompletion(messages []openai.ChatCompletionMessage) (string, error)
	OpenAIClient
}

const (
	SystemPromptTemplate = `
		You are a command-line app that responds with only a list of tags that categorize the content of the messages being sent to you.
		You can only provide AT MOST 3 tags, and at least 1 TAG.  Less is preferable.  Do not infer tags. The list of tags you output are:

		{{- range .}}
		- {{.}}
		{{- end}}
	`

	TagInputTemplate = `
		Title: {{.Title}}
		URL: {{.URL}}

		Content Raw: {{.Raw}}
	`
)

type TagInput struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Raw   string `json:"raw"`
}

func GenerateSystemPrompt(tags []string) string {
	tmpl, err := template.New("system-prompt").Parse(SystemPromptTemplate)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, tags)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func GenerateTagInputMessage(input *TagInput, tokenLimit int) string {
	tmpl, err := template.New("tag-input").Parse(TagInputTemplate)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, input)
	if err != nil {
		panic(err)
	}

	// Truncate the message to the token limit
	message := buf.String()
	if len(message) > tokenLimit {
		message = message[:tokenLimit]
	}

	return message
}

func SplitResponse(response string) []string {
	return strings.Split(response, "\n")
}
