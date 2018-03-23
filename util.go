package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
)

// NewRequest wrap http.NewRequest
func NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Referer", "http://angao.xyz")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.162 Safari/537.36")
	return req, nil
}

// NewClient custom http.Client
func NewClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
	}
}

// GenerateURL generate url by rawurl and param
func GenerateURL(rawurl, param string) string {
	if !(strings.HasPrefix(rawurl, "http://") || strings.HasPrefix(rawurl, "https://")) {
		rawurl = "http://" + rawurl
	}
	if param != "" {
		values := &url.Values{}
		// if param start with '{' and end with '}', then it's json
		if strings.HasPrefix(param, "{") && strings.HasSuffix(param, "}") {
			p := make(map[string]string)
			err := json.Unmarshal([]byte(param), &p)
			if err != nil {
				// param not a json string
				values = generateQuery(param)
			} else {
				for key, val := range p {
					values.Add(key, val)
				}
			}
		} else {
			values = generateQuery(param)
		}
		return fmt.Sprintf("%s?%s", rawurl, values.Encode())
	}
	return rawurl
}

func generateQuery(param string) *url.Values {
	values := &url.Values{}
	subs := strings.Split(param, "&")
	for _, s := range subs {
		sub := strings.SplitN(s, "=", 2)
		if len(sub) > 1 {
			values.Add(sub[0], sub[1])
		} else {
			values.Add(sub[0], "")
		}
	}
	return values
}

// PrintHeader print response header
func PrintHeader(resp *http.Response) {
	yellow := color.New(color.FgYellow)
	blue := color.New(color.FgBlue)
	code := resp.StatusCode
	if code >= 100 && code < 300 {
		color.Green("\n%s %s\n\n", resp.Proto, resp.Status)
	} else if code >= 300 && code < 400 {
		color.Yellow("\n%s %s\n\n", resp.Proto, resp.Status)
	} else {
		color.Red("\n%s %s\n\n", resp.Proto, resp.Status)
	}
	for key, values := range resp.Header {
		for i := range values {
			yellow.Printf("%s: ", key)
			blue.Printf("%s\n", values[i])
		}
	}
}

// PrintBody print response body
func PrintBody(resp *http.Response) {
	fmt.Printf("\n-------------------------------------------------\n")
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
		}
		fmt.Print(line)
	}
}

// Download download resource
func Download(resp *http.Response) {
	name, err := getFileName(resp)
	if err != nil {
		fmt.Println(err)
		return
	}

	path, _ := os.Getwd()
	f, err := os.Create(path + "/" + name)
	if err != nil {
		fmt.Println(err)
		return
	}

	writer := bufio.NewWriter(f)
	reader := bufio.NewReader(resp.Body)

	for {
		b, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		writer.WriteByte(b)
	}
	writer.Flush()
	f.Close()
}

func getFileName(resp *http.Response) (string, error) {
	var name string
	h := resp.Header.Get("Content-Disposition")
	if h == "" {
		p := resp.Request.URL.Path
		if p != "" {
			typ := ContentType(resp)
			var suffix string
			if Image(typ) {
				suffix = "." + strings.Split(typ, "image/")[1]
			}
			ss := strings.Split(p, "/")
			return ss[len(ss)-1] + suffix, nil
		}
		return "", errors.New("cann't found filename")
	}
	name = strings.Replace(h, "attachment;filename=", "", -1)
	return name, nil
}

// ContentType get content-type
func ContentType(resp *http.Response) string {
	return resp.Header.Get("Content-Type")
}

// OctetStream content-type is application/octet-stream
func OctetStream(typ string) bool {
	if typ == "application/octet-stream" {
		return true
	}
	return false
}

// Image content-type is image/*
func Image(typ string) bool {
	reg := regexp.MustCompile("image/.*")
	if len(reg.FindString(typ)) > 0 {
		return true
	}
	return false
}

// Text content-type is application/json, text/*
func Text(typ string) bool {
	if strings.Contains(typ, "application/json") {
		return true
	}
	if strings.Contains(typ, "application/xml") {
		return true
	}
	reg := regexp.MustCompile("text/.*")
	if len(reg.FindString(typ)) > 0 {
		return true
	}
	return false
}

// PrintUndefined print undefined content type
func PrintUndefined() {
	fmt.Println()
	fmt.Println("        |-------------------------------------|")
	fmt.Println("        |       Undefined Content-Type        |")
	fmt.Println("        |-------------------------------------|")
}

// HelpTemplate return help template
func HelpTemplate() string {
	return `Name:
	{{.Name}}{{if .Usage}} - {{.Usage}}{{end}}
 
Usage:
	{{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}
 
Version:
	{{.Version}}{{end}}{{end}}{{if .Description}}
 
Description:
	{{.Description}}{{end}}{{if len .Authors}}
 
Author{{with $length := len .Authors}}{{if ne 1 $length}}S{{end}}{{end}}:
	{{range $index, $author := .Authors}}{{if $index}}
	{{end}}{{$author}}{{end}}{{end}}{{if .VisibleCommands}}
 
Commands:{{range .VisibleCategories}}{{if .Name}}
	{{.Name}}:{{end}}{{range .VisibleCommands}}
	  {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
 
Global Options:
	{{range $index, $option := .VisibleFlags}}{{if $index}}
	{{end}}{{$option}}{{end}}{{end}}{{if .Copyright}}
 
Copyright:
	{{.Copyright}}{{end}}
 `
}
