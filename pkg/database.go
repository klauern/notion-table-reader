package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"

	"github.com/dstotijn/go-notion"
)

const DatabaseID = "2ce556682898478d8e9d175badac759e"

type PageWithBlocks struct {
	Page   *notion.Page
	Blocks []notion.Block
}

func extractRichText(richText []notion.RichText) string {
	var buf bytes.Buffer
	for _, t := range richText {
		buf.WriteString(t.PlainText)
	}
	return buf.String()
}

func blockToMarkdown(block notion.Block) string {
	switch b := block.(type) {
	case *notion.ParagraphBlock:
		return extractRichText(b.RichText)
	case *notion.Heading1Block:
		return extractRichText(b.RichText)
	case *notion.Heading2Block:
		return extractRichText(b.RichText)
	case *notion.Heading3Block:
		return extractRichText(b.RichText)
	case *notion.BulletedListItemBlock:
		return extractRichText(b.RichText)
	case *notion.NumberedListItemBlock:
		return extractRichText(b.RichText)
	case *notion.ToDoBlock:
		return extractRichText(b.RichText)
	case *notion.CalloutBlock:
		return extractRichText(b.RichText)
	default:
		return ""
	}
}

func (c *Client) ListMultiSelectProps(databaseId, columnName string) ([]string, error) {
	database, err := c.client.FindDatabaseByID(c.context, databaseId)
	if err != nil {
		return nil, fmt.Errorf("can't retrieve database: %w", err)
	}
	var props []string
	for _, prop := range database.Properties {
		if prop.Type == notion.DBPropTypeMultiSelect && prop.Name == columnName {
			for _, p := range prop.MultiSelect.Options {
				props = append(props, p.Name)
			}
			return props, nil
		}
	}
	return nil, fmt.Errorf("Unable to find column %s", columnName)
}

func (c *Client) ListDatabases(query string) ([]notion.Database, error) {
	resp, err := c.client.Search(c.context, &notion.SearchOpts{
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

func (c *Client) ListTagsForDatabaseColumn(databaseId, columnName string) ([]string, error) {
	database, err := c.client.FindDatabaseByID(c.context, DatabaseID)
	if err != nil {
		return nil, fmt.Errorf("Error finding database: %w", err)
	}

	var columns []string
	for _, prop := range database.Properties {
		if prop.Type == notion.DBPropTypeMultiSelect {
			for _, opt := range prop.MultiSelect.Options {
				columns = append(columns, opt.Name)
			}
			return columns, nil
		}
	}

	return nil, errors.New("No columns found")
}

func (c *Client) ListPages(notTagged bool) ([]notion.Page, error) {
	results, err := c.client.QueryDatabase(c.context, DatabaseID, &notion.DatabaseQuery{
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

func (c *Client) GetPage(pageId string) (*PageWithBlocks, error) {
	page, err := c.client.FindPageByID(c.context, pageId)
	if err != nil {
		return nil, fmt.Errorf("Error finding page: %w", err)
	}
	Log.Debug("page", "id", page.ID, "parent_id", page.Parent.PageID)
	blocks, err := c.client.FindBlockChildrenByID(c.context, page.ID, &notion.PaginationQuery{})
	if err != nil {
		return nil, fmt.Errorf("Error finding blocks: %w", err)
	}

	pageBlocks := PageWithBlocks{}
	pageBlocks.Page = &page
	for _, block := range blocks.Results {
		// If block is of BlockType Paragraph, NumberedListItem, Heading1, Heading2, or BulletedListItem, parse it and store
		var validBlock interface{}
		switch block := block.(type) {
		case *notion.ParagraphBlock:
			validBlock = *block
		case *notion.NumberedListItemBlock:
			validBlock = *block
		case *notion.BulletedListItemBlock:
			validBlock = *block
		case *notion.Heading1Block:
			validBlock = *block
		case *notion.Heading2Block:
			validBlock = *block
		case *notion.Heading3Block:
			validBlock = *block
		case *notion.CalloutBlock:
			validBlock = *block
		default:
			Log.Debug("unrecognized", "type", reflect.TypeOf(block))
			continue
		}
		pageBlocks.Blocks = append(pageBlocks.Blocks, validBlock.(notion.Block))
	}

	return &pageBlocks, nil
}

func (p PageWithBlocks) NormalizeBody() string {
	var buf bytes.Buffer
	for _, block := range p.Blocks {
		buf.WriteString(blockToMarkdown(block))
	}
	return buf.String()
}

func tagsToNotionProps(tags []string) []notion.SelectOptions {
	var notionTags []notion.SelectOptions
	for _, tag := range tags {
		notionTags = append(notionTags, notion.SelectOptions{
			Name: tag,
		})
	}
	return notionTags
}

func (c *Client) TagDatabasePage(pageId string, tags []string) error {
	_, err := c.client.UpdatePage(c.context, pageId, notion.UpdatePageParams{
		DatabasePageProperties: notion.DatabasePageProperties{
			"Tags": notion.DatabasePageProperty{
				MultiSelect: tagsToNotionProps(tags),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to update page %s with tags: %w", pageId, err)
	}
	return nil
}
