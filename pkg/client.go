package pkg

import (
	"context"
	"os"

	"github.com/dstotijn/go-notion"
)

type Client struct {
	client  *notion.Client
	context context.Context
}

func NewClient() *Client {
	return &Client{
		client:  notion.NewClient(os.Getenv("NOTION_API_KEY")),
		context: context.Background(),
	}
}
