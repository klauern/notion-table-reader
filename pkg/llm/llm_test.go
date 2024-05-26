package llm_test

import (
	"bytes"
	"context"
	"testing"
	"text/template"

	"github.com/klauern/notion-table-reader/pkg/llm"
	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
)

func TestGenerateSystemPrompt(t *testing.T) {
	tags := []string{"tag1", "tag2", "tag3"}
	expected := `
		You are a command-line app that responds with only a list of tags that categorize the content of the messages being sent to you.
		You can only provide AT MOST 3 tags, and at least 1 TAG.  Less is preferable.  Do not infer tags. The list of tags you output are:

		- tag1
		- tag2
		- tag3
	`
	result := llm.GenerateSystemPrompt(tags)
	assert.Equal(t, expected, result)
}

func TestGenerateTagInputMessage(t *testing.T) {
	input := &llm.TagInput{
		Title: "Test Title",
		URL:   "http://example.com",
		Raw:   "Test content",
	}
	expected := `
		Title: Test Title
		URL: http://example.com

		Content Raw: Test content
	`
	result := llm.GenerateTagInputMessage(input, 1000)
	assert.Equal(t, expected, result)
}

func TestGenerateTagInputMessage_Truncate(t *testing.T) {
	input := &llm.TagInput{
		Title: "Test Title",
		URL:   "http://example.com",
		Raw:   "Test content",
	}
	expected := `
		Title: Test Title
		URL: http://example.com

		Content Raw: Test content
	`
	result := llm.GenerateTagInputMessage(input, 10)
	assert.Equal(t, expected[:10], result)
}

func TestSplitResponse(t *testing.T) {
	response := "line1\nline2\nline3"
	expected := []string{"line1", "line2", "line3"}
	result := llm.SplitResponse(response)
	assert.Equal(t, expected, result)
}

type mockOpenAIClient struct{}

func (m *mockOpenAIClient) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	return openai.ChatCompletionResponse{
		Choices: []openai.ChatCompletionChoice{
			{Message: openai.ChatCompletionMessage{Content: "response"}},
		},
	}, nil
}

func TestRequestChatCompletion(t *testing.T) {
	client := &llm.Client{
		LLMClient: &mockOpenAIClient{},
		Context:   context.Background(),
		Model:     "test-model",
	}

	messages := []openai.ChatCompletionMessage{
		{Role: "user", Content: "test message"},
	}
	expected := "response"
	result, err := client.RequestChatCompletion(messages)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
