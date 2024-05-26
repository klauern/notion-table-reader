package pkg

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/sashabaranov/go-openai"
)

type MockOpenAIClient struct {
	CreateChatCompletionFunc func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}

func (m *MockOpenAIClient) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	return m.CreateChatCompletionFunc(ctx, req)
}

func TestIdentifyTags(t *testing.T) {
	mockClient := &MockOpenAIClient{
		CreateChatCompletionFunc: func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
			return openai.ChatCompletionResponse{
				Choices: []openai.ChatCompletionChoice{
					{
						Message: openai.ChatCompletionMessage{
							Content: "tag1\ntag2",
						},
					},
				},
			}, nil
		},
	}

	client := Client{
		llmClient: mockClient,
		context:   context.Background(),
		Model:     "test-model",
		MaxTokens: 100,
	}

	tagInput := &TagInput{
		Title: "Test Title",
		URL:   "http://test.url",
		Raw:   "Test Raw",
	}

	tags, err := client.IdentifyTags(tagInput, []string{"tag1", "tag2", "tag3"})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedTags := []string{"tag1", "tag2"}
	if !reflect.DeepEqual(tags, expectedTags) {
		t.Errorf("Expected tags %v, but got %v", expectedTags, tags)
	}

	// Test with error
	mockClient = &MockOpenAIClient{
		CreateChatCompletionFunc: func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
			return openai.ChatCompletionResponse{}, errors.New("error")
		},
	}

	client = Client{
		llmClient: mockClient,
		context:   context.Background(),
		Model:     "test-model",
		MaxTokens: 100,
	}

	_, err = client.IdentifyTags(tagInput, []string{"tag1", "tag2", "tag3"})
	if err == nil || err.Error() != "error creating chat completion request after 3 attempts: error" {
		t.Errorf("Expected error 'error creating chat completion request after 3 attempts: error', but got: %v", err)
	}
}
