package main

import (
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	e := &cli.App{
		Name: "notion",
		Commands: []*cli.Command{
			{
				Name:   "query",
				Action: QueryCommand,
			},
		},
	}
	err := e.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func QueryCommand(context *cli.Context) error {
	return nil
}
