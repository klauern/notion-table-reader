package pkg

import (
	"fmt"

	"github.com/dstotijn/go-notion"
)

const DatabaseID = "2ce556682898478d8e9d175badac759e"

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
	// fmt.Println(resp.Results)
	for _, result := range resp.Results {
		if database, ok := result.(notion.Database); ok {
			if database.Title != nil && len(database.Title) > 0 {
				if database.Title[0].PlainText != "" {
					databases = append(databases, database)
					fmt.Println(database.Title[0].PlainText)
				}
			}
		}
	}

	return databases, nil
}
