package pkg

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/klauern/notion-table-reader/pkg/llm"
	"github.com/klauern/notion-table-reader/pkg/mocks"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/mock/gomock"
)

func TestRequestChatCompletion_Success(t *testing.T) {
	// Create a new instance of our mock object
	ctrl := gomock.NewController(t)

	mockClient := mocks.NewMockLLMClient(ctrl)
	mockClient.EXPECT().CreateChatCompletion(gomock.Any(), gomock.Any()).Return(openai.ChatCompletionResponse{
		Choices: []openai.ChatCompletionChoice{
			{
				Message: openai.ChatCompletionMessage{
					Content: "Test content",
				},
			},
		},
	}, nil)

	client := Client{
		LLMClient: mockClient,
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
	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockLLMClient(ctrl)
	mockClient.EXPECT().CreateChatCompletion(gomock.Any(), gomock.Any()).Return(openai.ChatCompletionResponse{}, errors.New("error")).AnyTimes()

	client := Client{
		LLMClient: mockClient,
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
	client := NewClient(context.Background(), "openai_key", "notion_api_key")
	if client == nil {
		t.Error("Expected a non-nil client, but got nil")
	}
}

// func TestListTagsForDatabaseColumn(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	mockClient := mocks.NewMockNotionClient(ctrl)
// 	mockClient.EXPECT().FindDatabaseByID(gomock.Any(), gomock.Any()).Return(notion.Database{
// 		Properties: map[string]notion.DatabaseProperty{
// 			"Tags": {
// 				Type: notion.DBPropTypeMultiSelect,
// 				MultiSelect: &notion.SelectMetadata{
// 					Options: []notion.SelectOptions{
// 						{Name: "tag1"},
// 						{Name: "tag2"},
// 					},
// 				},
// 			},
// 		},
// 	}, nil,
// 	).AnyTimes()

// 	tags, err := mockClient.ListTagsForDatabaseColumn("databaseId", "Tags")
// 	if err != nil {
// 		t.Errorf("Unexpected error: %v", err)
// 	}

// 	expectedTags := []string{"tag1", "tag2"}
// 	if !reflect.DeepEqual(tags, expectedTags) {
// 		t.Errorf("Expected tags %v, but got %v", expectedTags, tags)
// 	}
// }

func TestIdentifyTags(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockOpenAIClient(ctrl)
	mockClient.EXPECT().CreateChatCompletion(gomock.Any(), gomock.Any()).Return(openai.ChatCompletionResponse{
		Choices: []openai.ChatCompletionChoice{
			{
				Message: openai.ChatCompletionMessage{
					Content: "tag1\ntag2",
				},
			},
		},
	}, nil).AnyTimes()

	client := Client{
		LLMClient: mockClient,
		context:   context.Background(),
		Model:     "test-model",
		MaxTokens: 100,
	}

	tagInput := &llm.TagInput{
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

	mockClient = mocks.NewMockOpenAIClient(ctrl)
	mockClient.EXPECT().CreateChatCompletion(gomock.Any(), gomock.Any()).Return(openai.ChatCompletionResponse{}, errors.New("error")).AnyTimes()

	client = Client{
		LLMClient: mockClient, context: context.Background(),
		Model:     "test-model",
		MaxTokens: 100,
	}

	_, err = client.IdentifyTags(tagInput, []string{"tag1", "tag2", "tag3"})
	if err == nil || err.Error() != "error creating chat completion request after 3 attempts: error" {
		t.Errorf("Expected error 'error creating chat completion request after 3 attempts: error', but got: %v", err)
	}
}
