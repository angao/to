package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli"
)

func main() {
	cli.HelpFlag = cli.BoolFlag{
		Name:  "help",
		Usage: "Show help info",
	}
	app := cli.NewApp()

	app.Name = "To"
	app.Usage = " To is a command line for http."
	app.Version = "0.0.1"
	cli.AppHelpTemplate = HelpTemplate()

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "header, h",
			Usage: "Just show http header",
		},
	}

	var reqParam string
	app.Commands = []cli.Command{
		{
			Name:    "get",
			Aliases: []string{"g"},
			Usage:   "HTTP get method",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "param, p",
					Usage:       "HTTP request param",
					Destination: &reqParam,
				},
			},
			Action: func(c *cli.Context) {
				url := c.Args().Get(0)
				if url == "" {
					fmt.Println("param error: url empty")
					return
				}
				url = GenerateURL(url, reqParam)
				method := strings.ToUpper(c.Command.Name)

				req, err := NewRequest(method, url, nil)
				if err != nil {
					fmt.Println(err)
					return
				}
				resp, err := NewClient(10 * time.Second).Do(req)
				if err != nil {
					fmt.Println(err)
					return
				}
				defer resp.Body.Close()
				PrintHeader(resp)
				typ := ContentType(resp)
				if !c.GlobalBool("header") {
					if Text(typ) {
						PrintBody(resp)
						return
					}
				}
				if !Text(typ) && OctetStream(typ) && Image(typ) {
					fmt.Println("\ndownloading file, wait a moment...")
					Download(resp)
				} else if !Text(typ) {
					PrintUndefined()
					return
				}
			},
		},
		{
			Name:    "post",
			Aliases: []string{"po"},
			Usage:   "HTTP post method",
		},
		{
			Name:    "put",
			Aliases: []string{"pt"},
			Usage:   "HTTP put method",
		},
		{
			Name:    "delete",
			Aliases: []string{"d"},
			Usage:   "HTTP delete method",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
