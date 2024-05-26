package pkg

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/dstotijn/go-notion"
)

type NotionClient interface {
	FetchPages(databaseID string, untagged bool) ([]PageDetail, error)
	TagPage(id string, availableTags []string) error
}

type PageDetail struct {
	ID   string
	Name string
}

// FetchPages returns a list of page details from the database.
func (l *Client) FetchPages(databaseID string, untagged bool) ([]PageDetail, error) {
	pages, err := l.ListPages(databaseID, untagged)
	if err != nil {
		return nil, fmt.Errorf("failed to query pages: %w", err)
	}

	var pageDetails []PageDetail
	for _, page := range pages {
		pageProps, ok := page.Properties.(notion.DatabasePageProperties)
		if !ok {
			return nil, fmt.Errorf("failed to convert page properties to notion.DatabasePageProperties")
		}
		name := pageProps["Name"].Title[0].PlainText
		pageDetails = append(pageDetails, PageDetail{
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

	tagList, err := l.IdentifyTags(NewTagInput(p), availableTags)
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

func NewTagInput(page *PageWithBlocks) *TagInput {
	tag := &TagInput{
		Title: page.Page.Properties.(notion.DatabasePageProperties)["Name"].Title[0].PlainText,
		URL:   page.Page.URL,
		Raw:   page.NormalizeBody(),
	}
	return tag
}

// PrintPageDetails prints the details of the provided pages.
func PrintPageDetails(pages []PageDetail) {
	for _, page := range pages {
		fmt.Printf("Page(%s): %s\n", page.ID, page.Name)
	}
}
