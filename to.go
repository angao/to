package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

// basic info
const (
	NAME    = "To"
	USAGE   = "To is a command line for HTTP."
	VERSION = "0.0.2"
	AUTHOR  = "Angao <gawaine2111@foxmail.com>"
)

func main() {
	app := cli.NewApp()

	app.Name = NAME
	app.Usage = USAGE
	app.Version = VERSION
	app.Author = AUTHOR

	cli.AppHelpTemplate = HelpTemplate
	cli.HelpFlag = HelpFlag
	cli.VersionFlag = VersionFlag

	app.Flags = []cli.Flag{
		HeaderFlag,
	}
	app.Commands = []cli.Command{
		GetCommand,
		PostCommand,
		PutCommand,
		DeleteCommand,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
