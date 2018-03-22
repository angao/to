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
		Usage: "show help info",
	}
	app := cli.NewApp()

	app.Name = "To"
	app.Usage = " To is a command line for http."
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "header, h",
			Usage: "show http header",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "get",
			Usage: "http get method",
			Action: func(c *cli.Context) {
				url := generateURL(c.Args().Get(0))
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
				printHeader(resp)
				if !c.GlobalBool("header") {
					printBody(resp)
				}
			},
		},
		{
			Name:  "post",
			Usage: "http post method",
		},
		{
			Name:  "put",
			Usage: "http put method",
		},
		{
			Name:  "delete",
			Usage: "http delete method",
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
