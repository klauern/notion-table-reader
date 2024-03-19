package main

import (
	"context"
	"fmt"
	"os"

	"github.com/dstotijn/go-notion"
)

const DatabaseID = "2ce556682898478d8e9d175badac759e"

func main() {
	client := notion.NewClient(os.Getenv("NOTION_API_KEY"))
	database, err := client.FindDatabaseByID(context.Background(), DatabaseID)
	if err != nil {
		panic(err)
	}

	for _, prop := range database.Properties {
		if prop.Type == notion.DBPropTypeMultiSelect {
			fmt.Println(prop.Name)
			for _, opt := range prop.MultiSelect.Options {
				fmt.Println(opt.ID, opt.Name)
			}
		}
	}

	result, err := client.QueryDatabase(context.Background(), DatabaseID, &notion.DatabaseQuery{
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
		panic(err)
	}

	for _, r := range result.Results {
		props, ok := r.Properties.(notion.DatabasePageProperties)
		if !ok {
			// Handle the case where the type assertion fails
			fmt.Println("Failed to convert r.Properties to notion.DatabasePageProperties")
			continue
		}
		fmt.Println(props["Name"].Title[0].PlainText)
	}
}
