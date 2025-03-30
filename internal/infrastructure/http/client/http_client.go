// MIT License
// Copyright (c) 2025 Toni Liesche
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

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
