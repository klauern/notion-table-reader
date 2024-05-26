package main

import (
	"context"

	"github.com/klauern/notion-table-reader/pkg"
)

func main() {
	pkg.NewClient(context.Background(), "", "").ListDatabases("")
}
