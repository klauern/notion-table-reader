package main

import "github.com/klauern/notion-table-reader/pkg"


func main() {
	pkg.NewClient().ListDatabases("")
}
