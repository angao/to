package main

import (
	"net/http"
	"time"
)

type req struct {
	method string
	url    string
	*http.Client
}

// New create req
func New(method, url string) *req {
	return &req{
		method: method,
		url:    url,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// NewRequest create request
func (r *req) NewRequest() (*http.Request, error) {
	rq, err := http.NewRequest(r.method, r.url, nil)
	if err != nil {
		return nil, err
	}
	rq.Header.Set("Host", "angao.xyz")
	rq.Header.Set("Referer", "http://angao.xyz")
	rq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.162 Safari/537.36")
	return rq, nil
}
