package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
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

// generateURL
func generateURL(url string) string {
	if !(strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")) {
		url = "http://" + url
	}
	return url
}

// printHeader print response header
func printHeader(resp *http.Response) {
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

// printBody print response body
func printBody(resp *http.Response) {
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

func download(resp *http.Response) {
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
			typ := contentType(resp)
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
func contentType(resp *http.Response) string {
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
