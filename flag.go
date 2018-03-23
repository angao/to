package main

import "github.com/urfave/cli"

var (
	// HeaderFlag show http header
	HeaderFlag = cli.BoolFlag{
		Name:  "header, h",
		Usage: "Just show http header",
	}

	// HelpFlag change default help
	HelpFlag = cli.BoolFlag{
		Name:  "help",
		Usage: "Show help info",
	}

	// VersionFlag change default version
	VersionFlag = cli.BoolFlag{
		Name:  "version, v",
		Usage: "Print the versions",
	}
)
