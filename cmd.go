package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/urfave/cli"
)

var reqParam string

var paramFlag = cli.StringFlag{
	Name:        "param, p",
	Usage:       "HTTP request param",
	Destination: &reqParam,
}

var (
	// GetCommand http get
	GetCommand = cli.Command{
		Name:    "get",
		Aliases: []string{"g"},
		Usage:   "HTTP get method",
		Flags: []cli.Flag{
			paramFlag,
		},
		Action: func(c *cli.Context) {
			url := c.Args().Get(0)
			if url == "" {
				fmt.Println("param error: url empty")
				return
			}
			url = GenerateURL(url)
			body := GenerateParam(reqParam)
			method := strings.ToUpper(c.Command.Name)

			req, err := NewRequest(method, url, body)
			if err != nil {
				fmt.Println(err)
				return
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			resp, err := NewClient(15 * time.Second).Do(req)
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
	}

	// PostCommand http post
	PostCommand = cli.Command{
		Name:    "post",
		Aliases: []string{"po"},
		Usage:   "HTTP post method",
		Flags: []cli.Flag{
			paramFlag,
		},
		Action: func(c *cli.Context) {

		},
	}

	// PutCommand http put
	PutCommand = cli.Command{
		Name:    "put",
		Aliases: []string{"pt"},
		Usage:   "HTTP put method",
		Action: func(c *cli.Context) {

		},
	}

	// DeleteCommand http delete
	DeleteCommand = cli.Command{
		Name:    "delete",
		Aliases: []string{"d"},
		Usage:   "HTTP delete method",
		Action: func(c *cli.Context) {

		},
	}
)
