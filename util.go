package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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
		color.Green("%s %s\n", resp.Proto, resp.Status)
	} else if code >= 300 && code < 400 {
		color.Yellow("%s %s\n", resp.Proto, resp.Status)
	} else {
		color.Red("%s %s\n", resp.Proto, resp.Status)
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
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
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
