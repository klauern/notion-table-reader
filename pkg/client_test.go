package pkg

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/dstotijn/go-notion"
	"github.com/sashabaranov/go-openai"
)

// MockClient is a mock type for pkg.Client
type MockClient struct {
	CreateChatCompletionFunc  func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
	FindDatabaseByIDFunc      func(ctx context.Context, databaseId string) (notion.Database, error)
	ListTagsForDatabaseColumnFunc func(dbId, colName string) ([]string, error)
}

// CreateChatCompletion is a mock method for CreateChatCompletion method of pkg.Client
func (m *MockClient) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	return m.CreateChatCompletionFunc(ctx, req)
}

func TestRequestChatCompletion_Success(t *testing.T) {
	// Create a new instance of our mock object
	mockClient := &MockClient{
		CreateChatCompletionFunc: func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
			return openai.ChatCompletionResponse{
				Choices: []openai.ChatCompletionChoice{
					{
						Message: openai.ChatCompletionMessage{
							Content: "Test content",
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

	// Test successful request
	messages := []openai.ChatCompletionMessage{{Content: "Test message"}}
	resp, err := client.RequestChatCompletion(messages)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if resp != "Test content" {
		t.Errorf("Expected 'Test content', but got: %v", resp)
	}
}

func TestRequestChatCompletion_Failure(t *testing.T) {
	// Create a new instance of our mock object
	mockClient := &MockClient{
		CreateChatCompletionFunc: func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
			return openai.ChatCompletionResponse{}, errors.New("error")
		},
	}

	client := Client{
		llmClient: mockClient,
		context:   context.Background(),
		Model:     "test-model",
		MaxTokens: 100,
	}

	// Test failed request
	messages := []openai.ChatCompletionMessage{{Content: "Test message"}}
	_, err := client.RequestChatCompletion(messages)
	if err == nil || err.Error() != "error creating chat completion request after 3 attempts: error" {
		t.Errorf("Expected error 'error creating chat completion request after 3 attempts: error', but got: %v", err)
	}
}

// Write rest of tests
func TestNewClient(t *testing.T) {
	// Test NewClient function
	client := NewClient("openai_key", "notion_api_key")
	if client == nil {
		t.Error("Expected a non-nil client, but got nil")
	}
}

func (m *MockClient) ListTagsForDatabaseColumn(dbId, colName string) ([]string, error) {
	return m.ListTagsForDatabaseColumnFunc(dbId, colName)
}

func TestListTagsForDatabaseColumn(t *testing.T) {
	mockClient := &MockClient{
		FindDatabaseByIDFunc: func(ctx context.Context, databaseId string) (notion.Database, error) {
			return notion.Database{
				Properties: map[string]notion.DatabaseProperty{
					"Tags": {
						Type: notion.DBPropTypeMultiSelect,
						MultiSelect: &notion.SelectMetadata{
							Options: []notion.SelectOptions{
								{Name: "tag1"},
								{Name: "tag2"},
							},
						},
					},
				},
			}, nil
		},
		ListTagsForDatabaseColumn: func(dbId, colName string) ([]string, error) {
			return []string{"tag1", "tag2"}, nil
		},
	}

	client := Client{
		notionClient: &notion.Client{}, // Use an empty notion.Client for testing
		context:      context.Background(),
	}

	tags, err := client.ListTagsForDatabaseColumn("databaseId", "Tags")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedTags = []string{"tag1", "tag2"}
	if !reflect.DeepEqual(tags, expectedTags) {
		t.Errorf("Expected tags %v, but got %v", expectedTags, tags)
	}
	mockClient := &MockClient{
		FindDatabaseByIDFunc: func(ctx context.Context, databaseId string) (notion.Database, error) {
			return notion.Database{
				Properties: map[string]notion.DatabaseProperty{
					"Tags": {
						Type: notion.DBPropTypeMultiSelect,
						MultiSelect: &notion.SelectMetadata{
							Options: []notion.SelectOptions{
								{Name: "tag1"},
								{Name: "tag2"},
							},
						},
					},
				},
			}, nil
		},
	}

	tags, err = mockClient.ListTagsForDatabaseColumn("databaseId", "Tags")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedTags := []string{"tag1", "tag2"}
	if !reflect.DeepEqual(tags, expectedTags) {
		t.Errorf("Expected tags %v, but got %v", expectedTags, tags)
	}
}
