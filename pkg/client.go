package pkg

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/dstotijn/go-notion"
	"github.com/sashabaranov/go-openai"
)

type Client struct {
	llmClient    OpenAIClient
	context      context.Context
	Model        string
	MaxTokens    int
	notionClient *notion.Client
}

var tokenMax map[string]int = map[string]int{
	openai.GPT4o: 4096,
}

// NewClient creates a new client for the given API keys and returns a *Client.
func NewClient(openai_key string, notion_api_key string) *Client {
	if openai_key == "" {
		openai_key = os.Getenv("OPENAI_API_KEY")
	}
	if notion_api_key == "" {
		notion_api_key = os.Getenv("NOTION_API_KEY")
	}
	config := openai.DefaultConfig(openai_key)
	model := openai.GPT4TurboPreview
	maxToken := tokenMax[model]
	return &Client{
		context:      context.Background(),
		llmClient:    openai.NewClientWithConfig(config),
		notionClient: notion.NewClient(notion_api_key),
		Model:        model,
		MaxTokens:    maxToken,
	}
}

// RequestChatCompletion returns a chat completion response.
func (l *Client) RequestChatCompletion(messages []openai.ChatCompletionMessage) (string, error) {
	var resp openai.ChatCompletionResponse
	var err error

	retries := 3
	for i := 0; i < retries; i++ {
		resp, err = l.llmClient.CreateChatCompletion(l.context, openai.ChatCompletionRequest{
			Model:     l.Model,
			Messages:  messages,
			MaxTokens: l.MaxTokens,
		})
		if err == nil {
			break
		}
		slog.Error("Getting chat completion", "resp", err)

		fmt.Printf("Error creating chat completion request (attempt %d): %v\n", i+1, err)
		if i < retries-1 {
			fmt.Println("Retrying...")
		}
	}

	if err != nil {
		return "", fmt.Errorf("error creating chat completion request after %d attempts: %w", retries, err)
	}

	slog.Debug("number of responses", "count", len(resp.Choices))
	return resp.Choices[0].Message.Content, nil
}
