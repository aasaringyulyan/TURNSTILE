package infrastructure

import (
	"bytes"
	"encoding/json"
	"net/http"
	"turnstile/internal/models"
	"turnstile/pkg/client"
	"turnstile/pkg/logging"
)

type LogClient struct {
	logger logging.Logger
	client *client.Client
	url    string
}

func NewLogClient(logger logging.Logger, client *client.Client, url string) *LogClient {
	return &LogClient{
		logger: logger,
		client: client,
		url:    url,
	}
}

func (lc *LogClient) SendLogs(logs models.PassageLogsForApi) error {
	data, err := json.Marshal(logs)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, lc.url, bytes.NewReader(data))
	if err != nil {
		return err
	}

	res, err := lc.client.HttpClient.Do(req)
	if err != nil {
		return err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	return nil
}
