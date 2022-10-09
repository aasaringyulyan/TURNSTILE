package client

import (
	"net/http"
	"time"
)

type Config struct {
	MaxConnsPerHost int
	Timeout         time.Duration
}

type Client struct {
	HttpClient *http.Client
}

func NewClient(cfg Config) *Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxConnsPerHost: cfg.MaxConnsPerHost,
		},
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       cfg.Timeout,
	}

	return &Client{
		HttpClient: client,
	}
}
