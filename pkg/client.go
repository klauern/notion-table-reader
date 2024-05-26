package pkg

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/dstotijn/go-notion"
	"github.com/klauern/notion-table-reader/pkg/llm"
	notionTypes "github.com/klauern/notion-table-reader/pkg/notion"
	"github.com/sashabaranov/go-openai"
)

type Client struct {
	LLMClient    llm.OpenAIClient
	context      context.Context
	Model        string
	MaxTokens    int
	NotionClient notionTypes.NotionClient
}

var tokenMax map[string]int = map[string]int{
	openai.GPT4o: 4096,
}

// NewClient creates a new client for the given API keys and returns a *Client.
func NewClient(ctx context.Context, openai_key string, notion_api_key string) *Client {
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
		context:      ctx,
		LLMClient:    openai.NewClientWithConfig(config),
		NotionClient: notion.NewClient(notion_api_key),
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
		resp, err = l.LLMClient.CreateChatCompletion(l.context, openai.ChatCompletionRequest{
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

func (l *Client) IdentifyTags(messageContent *llm.TagInput, tagOptions []string) ([]string, error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    "system",
			Content: llm.GenerateSystemPrompt(tagOptions),
		},
		{
			Role:    "user",
			Content: llm.GenerateTagInputMessage(messageContent, l.MaxTokens),
		},
	}

	response, err := l.RequestChatCompletion(messages)
	if err != nil {
		return nil, err
	}

	return llm.SplitResponse(response), nil
}

// FetchPages returns a list of page details from the database.
func (l *Client) FetchPages(databaseID string, untagged bool) ([]notionTypes.PageDetail, error) {
	pages, err := l.ListPages(databaseID, untagged)
	if err != nil {
		return nil, fmt.Errorf("failed to query pages: %w", err)
	}

	var pageDetails []notionTypes.PageDetail
	for _, page := range pages {
		pageProps, ok := page.Properties.(notion.DatabasePageProperties)
		if !ok {
			return nil, fmt.Errorf("failed to convert page properties to notion.DatabasePageProperties")
		}
		name := pageProps["Name"].Title[0].PlainText
		pageDetails = append(pageDetails, notionTypes.PageDetail{
			ID:   page.ID,
			Name: name,
		})
	}

	return pageDetails, nil
}

func (l *Client) TagPage(id string, availableTags []string) error {
	p, err := l.GetPage(id)
	if err != nil {
		return fmt.Errorf("failed to retrive Notion Page: %w", err)
	}

	tagList, err := l.IdentifyTags(notionTypes.NewTagInput(p), availableTags)
	if err != nil {
		return fmt.Errorf("failed to identify tags for page %s: %w", id, err)
	}

	slog.Info("Tagging page", "page", id, "tags", strings.Join(tagList, ", "))
	if err := l.TagDatabasePage(id, tagList); err != nil {
		slog.Error("Failed to tag page", "page", id, "err", err)
		return fmt.Errorf("failed to tag page %s: %w", id, err)
	}
	return nil
}
