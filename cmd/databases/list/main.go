package main

import "github.com/klauern/notion-table-reader/pkg"

const DatabaseID = "2ce556682898478d8e9d175badac759e"

func main() {
	pkg.NewClient().ListDatabases()
}
