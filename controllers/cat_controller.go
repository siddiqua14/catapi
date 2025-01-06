package controllers

import (
	"net/http"
	"runtime"
	"github.com/beego/beego/v2/server/web"
	
)
type HTTPClient interface {
    Do(req *http.Request) (*http.Response, error)
}

type CatController struct {
	web.Controller
	httpClient HTTPClient
}
// SetHTTPClient allows injection of mock client for testing
func (c *CatController) SetHTTPClient(client HTTPClient) {
    c.httpClient = client
}

// getHTTPClient returns the http client to use
func (c *CatController) getHTTPClient() HTTPClient {
    if c.httpClient != nil {
        return c.httpClient
    }
    return &http.Client{} // default client
}

// getCallerInfo returns the calling function's file and line number
func getCallerInfo() (string, int) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return "unknown", 0
	}
	return file, line
}

type MockHTTPClient struct {
    DoFunc func(req *http.Request) (*http.Response, error)
}



type CatImage struct {
	ID        string   `json:"id"`
	URL       string   `json:"url"`
	Width     int      `json:"width"`
	Height    int      `json:"height"`
	MimeType  string   `json:"mime_type"`
	Breeds    []Breed  `json:"breeds"`
	Categories []string `json:"categories"`
}
type Vote struct {
	ImageID string `json:"image_id"`
	Value   int    `json:"value"`
}

type Breed struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Origin       string `json:"origin"`
	WikipediaURL string `json:"wikipedia_url"`
}




