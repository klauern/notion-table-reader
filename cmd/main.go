package main

import (
	"fmt"
	"os"

	"github.com/dstotijn/go-notion"
	"github.com/klauern/notion-table-reader/pkg"
	"github.com/urfave/cli/v2"
)

var client *pkg.Client

func init() {
	client = pkg.NewClient()
}

func main() {
	e := &cli.App{
		Name: "notion",
		Commands: []*cli.Command{
			{
				Name:    "database",
				Aliases: []string{"db"},
				Subcommands: []*cli.Command{
					{
						Name:        "query",
						Description: "Query databases by Database name",
						Args:        true,
						Action:      QueryDatabase,
					},
				},
			},
		},
	}
	err := e.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func QueryDatabase(context *cli.Context) error {
	dbs, err := client.ListDatabases(context.Args().First())
	if err != nil {
		return fmt.Errorf("failed to query databases: %w", err)
	}
	for _, db := range dbs {
		// Print the parent of the database
		parent := db.Parent
		switch db.Parent.Type {
		case notion.ParentTypeDatabase:
			fmt.Printf("Database(%s) -> %s: %s\n", parent.DatabaseID, db.Title[0].PlainText, db.ID)
		case notion.ParentTypePage:
			fmt.Printf("Page(%s) -> %s: %s\n", parent.PageID, db.Title[0].PlainText, db.ID)
		case notion.ParentTypeWorkspace:
			fmt.Printf("Workspace -> %s: %s\n", db.Title[0].PlainText, db.ID)
		default:
			fmt.Println("Unknown parent type")
		}
	}
	return nil
}
