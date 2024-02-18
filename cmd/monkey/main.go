package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/zhulik/monkey/repl"
)

func main() {
	app := cli.App{
		Name:  "monkey",
		Usage: "Monkey interpreter",
		Action: func(ctx *cli.Context) error {
			file := ctx.Args().Get(0)
			if file == "" {
				err := repl.Start()
				if err != nil {
					return fmt.Errorf("repl error: %w", err)
				}
			}

			return nil // read and execute file
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
