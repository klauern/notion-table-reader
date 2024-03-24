package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dstotijn/go-notion"
	"github.com/klauern/notion-table-reader/pkg"
	"github.com/urfave/cli/v2"
)

const DatabaseID = "2ce556682898478d8e9d175badac759e"

var (
	client        *pkg.Client
	availableTags []string
	version       = "dev"
	commit        = "none"
	date          = "unknown"
)

func init() {
	client = pkg.NewClient("", "")
	tags, err := client.ListTagsForDatabaseColumn(DatabaseID, "Tags")
	if err != nil {
		panic(err)
	}
	availableTags = tags
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
					{
						Name:        "tags",
						Action:      ListTags,
						Description: "List all tags for the default DatabaseId",
					},
				},
			},
			{
				Name:    "pages",
				Aliases: []string{"p"},
				Subcommands: []*cli.Command{
					{
						Name:        "query",
						Description: "Query pages by type",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:  "untagged",
								Usage: "Include untagged pages in results",
							},
						},
						Action: QueryPages,
					},
					{
						Name:        "tag",
						Description: "Tag a page using the LLM results",
						Flags: []cli.Flag{
							&cli.StringSliceFlag{
								Name:  "page_id",
								Usage: "Page ID to tag",
							},
						},
						Action: TagPages,
					},
				},
			},
			{
				Name:    "version",
				Aliases: []string{"v"},
				Action: func(ctx *cli.Context) error {
					fmt.Printf("notion-tagger v%s, commit %s, built at %s\n", version, commit, date)
					return nil
				},
				Description: "print the version of the binary",
			},
		},
	}
	err := e.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func ListTags(context *cli.Context) error {
	tags, err := client.ListTagsForDatabaseColumn(DatabaseID, "Tags")
	if err != nil {
		return fmt.Errorf("failed to list tags for column: %w", err)
	}
	for _, tag := range tags {
		fmt.Println(tag)
	}
	return nil
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

func QueryPages(context *cli.Context) error {
	pages, err := client.ListPages(context.Bool("untagged"))
	if err != nil {
		return fmt.Errorf("failed to query pages: %w", err)
	}
	for _, page := range pages {
		pageProps, ok := page.Properties.(notion.DatabasePageProperties)
		if !ok {
			return fmt.Errorf("failed to convert page properties to notion.PageProperties")
		}
		fmt.Printf("Page(%s): %s\n", page.ID, pageProps["Name"].Title[0].PlainText)
	}
	return nil
}

func TagPages(context *cli.Context) error {
	errs := make([]error, 0)
	// for each 'page_id' in the stringslice in context, pull the page from Notion

	for _, id := range context.StringSlice("page_id") {
		page, err := client.GetPage(id)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to get page %s: %w", id, err))
		}
		tagList, err := client.IdentifyTags(&pkg.TagInput{
			Title: page.Page.Properties.(notion.DatabasePageProperties)["Name"].Title[0].PlainText,
			URL:   page.Page.URL,
			Raw:   page.NormalizeBody(),
		}, availableTags)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to identify tags for page %s: %w", id, err))
		}
		fmt.Printf("Tagging page %s with tags: %s\n", id, strings.Join(tagList, ", "))
		if err := client.TagDatabasePage(id, tagList); err != nil {
			errs = append(errs, fmt.Errorf("failed to tag page %s: %w", id, err))
		}
	}
	if len(errs) != 0 {
		// return all the errors wrapped in an error:
		return fmt.Errorf("%v", errs)
	}
	return nil
}
