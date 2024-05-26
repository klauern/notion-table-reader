package pkg

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/dstotijn/go-notion"
	readNotion "github.com/klauern/notion-table-reader/pkg/notion"
)

func (l *Client) ListMultiSelectProps(databaseId, columnName string) ([]string, error) {
	if l.NotionClient == nil {
		return nil, errors.New("NotionClient is not initialized")
	}
	if l.context == nil {
		return nil, errors.New("context is not initialized")
	}
	if l.NotionClient == nil {
		return nil, errors.New("NotionClient is not initialized")
	}
	if l.context == nil {
		return nil, errors.New("context is not initialized")
	}
	if l.NotionClient == nil {
		return nil, errors.New("NotionClient is not initialized")
	}
	if l.context == nil {
		return nil, errors.New("context is not initialized")
	}
	database, err := l.NotionClient.FindDatabaseByID(l.context, databaseId)
	if err != nil {
		return nil, fmt.Errorf("can't retrieve database: %w", err)
	}
	var props []string
	for _, prop := range database.Properties {
		if prop.Type == notion.DBPropTypeMultiSelect && prop.Name == columnName {
			for _, p := range prop.Select.Options {
				props = append(props, p.Name)
			}
			return props, nil
		}
	}
	return nil, fmt.Errorf("Unable to find column %s", columnName)
}

func (l *Client) ListDatabases(query string) ([]notion.Database, error) {
	resp, err := l.NotionClient.Search(l.context, &notion.SearchOpts{
		Query: query,
		Filter: &notion.SearchFilter{
			Value:    "database",
			Property: "object",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("Error querying for databases: %w", err)
	}

	var databases []notion.Database
	for _, result := range resp.Results {
		if database, ok := result.(notion.Database); ok {
			if database.Title != nil && len(database.Title) > 0 && database.Title[0].PlainText != "" {
				databases = append(databases, database)
			}
		}
	}

	return databases, nil
}

func (l *Client) ListTagsForDatabaseColumn(databaseId, columnName string) ([]string, error) {
	database, err := l.NotionClient.FindDatabaseByID(l.context, databaseId)
	if err != nil {
		return nil, fmt.Errorf("Error finding database: %w", err)
	}

	var columns []string
	for _, prop := range database.Properties {
		if prop.Type == notion.DBPropTypeMultiSelect {
			for _, opt := range prop.Select.Options {
				columns = append(columns, opt.Name)
			}
			return columns, nil
		}
	}

	return nil, errors.New("No columns found")
}

func (l *Client) ListPages(databaseId string, notTagged bool) ([]notion.Page, error) {
	results, err := l.NotionClient.QueryDatabase(l.context, databaseId, &notion.DatabaseQuery{
		Filter: &notion.DatabaseQueryFilter{
			Property: "Tags",
			DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
				MultiSelect: &notion.MultiSelectDatabaseQueryFilter{
					IsEmpty: true,
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("Error querying database: %w", err)
	}
	return results.Results, nil
}

func (l *Client) GetPage(pageId string) (*readNotion.PageWithBlocks, error) {
	page, err := l.NotionClient.FindPageByID(l.context, pageId)
	if err != nil {
		return nil, fmt.Errorf("Error finding page: %w", err)
	}
	slog.Debug("page", "id", page.ID, "parent_id", page.Parent.PageID)
	blocks, err := l.NotionClient.FindBlockChildrenByID(l.context, page.ID, &notion.PaginationQuery{})
	if err != nil {
		return nil, fmt.Errorf("Error finding blocks: %w", err)
	}

	pageWithBlocks := readNotion.PageWithBlocks{
		Page:   &page,
		Blocks: blocks.Results,
	}

	return &pageWithBlocks, nil
}

func TagsToNotionProps(tags []string) []notion.SelectOptions {
	var notionTags []notion.SelectOptions
	for _, tag := range tags {
		notionTags = append(notionTags, notion.SelectOptions{
			Name: tag,
		})
	}
	return notionTags
}

func (l *Client) TagDatabasePage(pageId string, tags []string) error {
	_, err := l.NotionClient.UpdatePage(l.context, pageId, notion.UpdatePageParams{
		DatabasePageProperties: notion.DatabasePageProperties{
			"Tags": notion.DatabasePageProperty{
				MultiSelect: TagsToNotionProps(tags),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to update page %s with tags: %w", pageId, err)
	}
	return nil
}
