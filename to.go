package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

// VERSION version
const (
	NAME    = "To"
	USAGE   = " To is a command line for http."
	VERSION = "0.0.1"
)

func main() {
	app := cli.NewApp()

	app.Name = NAME
	app.Usage = USAGE
	app.Version = VERSION
	cli.AppHelpTemplate = HelpTemplate()

	app.Flags = []cli.Flag{
		headerFlag,
	}

	cli.HelpFlag = helpFlag
	cli.VersionFlag = versionFlag

	app.Commands = []cli.Command{
		getCommand,
		postCommand,
		putCommand,
		deleteCommand,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
