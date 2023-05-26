package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/werbenhu/digo"
)

func main() {

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "path",
				Value: "",
				Usage: "language for the greeting",
			},
		},
		Action: func(cCtx *cli.Context) error {
			parser := digo.NewParser()
			parser.Start()
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
