package pkg

import (
	"context"
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
