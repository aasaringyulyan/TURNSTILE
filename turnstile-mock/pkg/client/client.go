package client

import (
	"net/http"
	"time"
	"turnstile-mock/pkg/logging"
	"turnstile-mock/pkg/service"
)

type Config struct {
	MaxConnsPerHost int
	Timeout         time.Duration
}

type Client struct {
	HttpClient *http.Client
	logger     logging.Logger
	services   *service.Service
}

func NewClient(logger logging.Logger, services *service.Service, cfg Config) *Client {
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
		logger:     logger,
		services:   services,
	}
}
