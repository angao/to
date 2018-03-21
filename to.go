package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
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

				rq := New(http.MethodGet, url)
				req, err := rq.NewRequest()
				if err != nil {
					fmt.Println(err)
					return
				}
				resp, err := rq.Do(req)
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

// generateURL
func generateURL(url string) string {
	if !(strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")) {
		url = "http://" + url
	}
	return url
}

// printHeader print response header
func printHeader(resp *http.Response) {
	color.Yellow("%s %s\n", resp.Proto, resp.Status)
	for key, values := range resp.Header {
		for i := range values {
			color.Blue("%s: %s\n", key, values[i])
		}
	}
}

// printBody print response body
func printBody(resp *http.Response) {
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}
