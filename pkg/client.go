package pkg

import (
	"context"
	"fmt"
	"os"

	"github.com/dstotijn/go-notion"
	"github.com/sashabaranov/go-openai"
)

const MyOrgID = "org-f6O3xgu6hFn7sGzM2dwRPt0v"

type Client struct {
	client    *notion.Client
	context   context.Context
	llmClient *openai.Client
}

func NewClient() *Client {
	config := openai.DefaultConfig(os.Getenv("OPENAI_API_KEY"))
	config.OrgID = MyOrgID
	client := openai.NewClientWithConfig(config)
	return &Client{
		client:    notion.NewClient(os.Getenv("NOTION_API_KEY")),
		context:   context.Background(),
		llmClient: client,
	}
}

func (l Client) RequestChatCompletion(messages []openai.ChatCompletionMessage, maxTokens int) (string, error) {
	var resp openai.ChatCompletionResponse
	var err error

	retries := 3
	for i := 0; i < retries; i++ {
		resp, err = l.llmClient.CreateChatCompletion(l.context, openai.ChatCompletionRequest{
			Model:     "gpt-3.5-turbo-16k-0613",
			Messages:  messages,
			MaxTokens: maxTokens,
		})

		if err == nil {
			break
		}

		fmt.Printf("Error creating chat completion request (attempt %d): %v\n", i+1, err)
		if i < retries-1 {
			fmt.Println("Retrying...")
		}
	}

	if err != nil {
		return "", fmt.Errorf("error creating chat completion request after %d attempts: %w", retries, err)
	}

	return resp.Choices[0].Message.Content, nil
}
