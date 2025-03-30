package client

import (
	"io"
	"net/http"
	"net/url"
	"time"
)

type HttpClient struct {
	baseUrl *url.URL
	client  *http.Client
}

func (c *HttpClient) NewRequest(method string, path string, body io.Reader) (*http.Request, error) {
	fullURL := c.baseUrl.ResolveReference(&url.URL{Path: path}).String()

	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (c *HttpClient) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

func NewHttpClient(baseURL string, timeout time.Duration) (*HttpClient, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return &HttpClient{
		baseUrl: parsedURL,
		client: &http.Client{
			Timeout: timeout,
		},
	}, nil
}
