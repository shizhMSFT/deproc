package main

import (
	"os"

	"github.com/apex/log"
	apexcli "github.com/apex/log/handlers/cli"
	"github.com/urfave/cli/v2"
)

func main() {
	log.SetHandler(apexcli.New(os.Stderr))
	app := &cli.App{
		Name:    "deproc",
		Usage:   "Reserve Proxy for Debugging",
		Version: "0.1.0",
		Authors: []*cli.Author{
			{
				Name:  "Shiwei Zhang",
				Email: "shizh@microsoft.com",
			},
		},
		Commands: []*cli.Command{
			serveCommand,
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err.Error())
	}
}
