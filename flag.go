package main

import "github.com/urfave/cli"

var (
	headerFlag = cli.BoolFlag{
		Name:  "header, h",
		Usage: "Just show http header",
	}

	helpFlag = cli.BoolFlag{
		Name:  "help",
		Usage: "Show help info",
	}

	versionFlag = cli.BoolFlag{
		Name:  "version, v",
		Usage: "Print the versions",
	}
)
