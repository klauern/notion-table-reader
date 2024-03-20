package pkg

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/dstotijn/go-notion"
)

const DatabaseID = "2ce556682898478d8e9d175badac759e"

type PageWithBlocks struct {
	Page   *notion.Page
	Blocks []*notion.ParagraphBlock
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
	blocks, err := c.client.FindBlockChildrenByID(c.context, page.ID, &notion.PaginationQuery{})
	if err != nil {
		return nil, fmt.Errorf("Error finding blocks: %w", err)
	}

	pageBlocks := PageWithBlocks{}
	pageBlocks.Page = &page
	for _, block := range blocks.Results {
		// If block is of BlockType Paragraph, parse it and store
		switch block := block.(type) {
		case *notion.ParagraphBlock:
			pageBlocks.Blocks = append(pageBlocks.Blocks, block)
		default:
			fmt.Printf("Unsupported block type: %T\n", block)
		}
	}

	return &pageBlocks, nil
}

func (p PageWithBlocks) NormalizeBody() string {
	var buf bytes.Buffer
	for _, block := range p.Blocks {
		for _, text := range block.RichText {
			buf.WriteString(text.PlainText)
		}
		buf.WriteString("\n")
	}
	return buf.String()
}
