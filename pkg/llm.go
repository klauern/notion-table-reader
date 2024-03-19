package pkg

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/sashabaranov/go-openai"
)

const (
	SystemPromptTemplate = `
		You are a command-line app that responds with only a list of tags that categorize the content of the messages being sent to you.
		You can only provide AT MOST 3 tags, and at least 1 TAG.  Less is preferable.  Do not infer tags. The list of tags you output are:

		{{-range .}}
		- {{.}}
		{{-end}}
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

func GenerateTagInputMessage(input *TagInput) string {
	tmpl, err := template.New("tag-input").Parse(TagInputTemplate)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, input)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

type LLMClient struct {
	openAI  *openai.Client
	context context.Context
}

func NewOpenAIClient(ctx context.Context) LLMClient {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	return LLMClient{
		openAI:  client,
		context: ctx,
	}
}

func (l Client) IdentifyTags(messageContent *TagInput, tagOptions []string) ([]string, error) {
	resp, err := l.llmClient.CreateChatCompletion(l.context, openai.ChatCompletionRequest{
		Model: openai.GPT4TurboPreview,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: GenerateSystemPrompt(tagOptions),
			},
			{
				Role:    "user",
				Content: GenerateTagInputMessage(messageContent),
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error creating chat completion request: %w", err)
	}
	return splitResponse(resp.Choices[0].Message.Content), nil
}

func splitResponse(response string) []string {
	return strings.Split(response, "\n")
}
