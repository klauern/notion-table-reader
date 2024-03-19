package pkg

import (
	"context"
	"os"

	"github.com/dstotijn/go-notion"
	"github.com/sashabaranov/go-openai"
)

type Client struct {
	client    *notion.Client
	context   context.Context
	llmClient *openai.Client
}

func NewClient() *Client {
	return &Client{
		client:    notion.NewClient(os.Getenv("NOTION_API_KEY")),
		context:   context.Background(),
		llmClient: openai.NewClient(os.Getenv("OPENAI_API_KEY")),
	}
}
