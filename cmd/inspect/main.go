package main

import (
	"context"
	"fmt"

	"github.com/klauern/notion-table-reader/pkg"
)

const DatabaseID = "2ce556682898478d8e9d175badac759e"

func main() {
	tags, err := pkg.NewClient(context.Background(), "", "").ListTagsForDatabaseColumn(DatabaseID, "Tags")
	if err != nil {
		panic(err)
	}

	for _, tag := range tags {
		fmt.Println(tag)
	}
}
