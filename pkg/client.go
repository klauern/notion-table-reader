package pkg

import (
	"context"
	"fmt"
	"os"

	"github.com/dstotijn/go-notion"
	"github.com/sashabaranov/go-openai"
)

type Client struct {
	client    *notion.Client
	context   context.Context
	llmClient *openai.Client
	Model     string
	MaxTokens int
}

var tokenMax map[string]int = map[string]int{
	openai.GPT4TurboPreview: 4096,
}

func NewClient(openai_key string, notion_api_key string) *Client {
	if openai_key == "" {
		openai_key = os.Getenv("OPENAI_API_KEY")
	}
	if notion_api_key == "" {
		notion_api_key = os.Getenv("NOTION_API_KEY")
	}
	config := openai.DefaultConfig(openai_key)
	client := openai.NewClientWithConfig(config)
	model := openai.GPT4TurboPreview
	maxToken := tokenMax[model]
	return &Client{
		client:    notion.NewClient(notion_api_key),
		context:   context.Background(),
		llmClient: client,
		Model:     model,
		MaxTokens: maxToken,
	}
}

func (l Client) RequestChatCompletion(messages []openai.ChatCompletionMessage) (string, error) {
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
		Log.Error("Getting chat completion", "resp", err)

		fmt.Printf("Error creating chat completion request (attempt %d): %v\n", i+1, err)
		if i < retries-1 {
			fmt.Println("Retrying...")
		}
	}

	if err != nil {
		return "", fmt.Errorf("error creating chat completion request after %d attempts: %w", retries, err)
	}

	Log.Debug("number of responses", "count", len(resp.Choices))
	return resp.Choices[0].Message.Content, nil
}
