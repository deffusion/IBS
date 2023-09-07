package main

import (
	"github.com/deffusion/IBS/cmd"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := cli.App{
		Name:     "IBS",
		Commands: cmd.Root,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
