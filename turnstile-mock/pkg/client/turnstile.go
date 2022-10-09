package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
	"turnstile-mock/models"
)

const (
	checkUrl = "http://localhost:8000/check"
	logsUrl  = "http://localhost:8000/logs"
)

type ConfigTimeout struct {
	DeltaTime time.Duration
}

func (c *Client) PostPassage(cfg ConfigTimeout) error {
	logger := c.logger.Logger

	for {
		passage, err := c.services.Passage.GeneratePassage()
		if err != nil {
			return newErrorResponse(http.StatusInternalServerError, err.Error(), c.logger)
		}

		data, err := json.Marshal(passage)
		if err != nil {
			return newErrorResponse(http.StatusInternalServerError, err.Error(), c.logger)
		}

		req, err := http.NewRequest(http.MethodPost, checkUrl, bytes.NewReader(data))
		if err != nil {
			return newErrorResponse(http.StatusInternalServerError, err.Error(), c.logger)
		}

		logger.Infof("trying to send passage: %v", passage)
		res, err := c.HttpClient.Do(req)
		if err != nil {
			return newErrorResponse(http.StatusInternalServerError, err.Error(), c.logger)
		}

		if res.Body != nil {
			defer res.Body.Close()
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return newErrorResponse(http.StatusInternalServerError, err.Error(), c.logger)
		}

		var isAccepted map[string]bool

		err = json.Unmarshal(body, &isAccepted)
		if err != nil {
			return newErrorResponse(http.StatusInternalServerError, err.Error(), c.logger)
		}

		logger.Infof("response: %v", isAccepted)

		if isAccepted["accepted"] {
			err = c.PostLogs(passage)
			if err != nil {
				return newErrorResponse(http.StatusInternalServerError, err.Error(), c.logger)
			}
		}

		time.Sleep(cfg.DeltaTime)
	}
}

func (c *Client) PostLogs(passage models.PassageCheck) error {
	logger := c.logger.Logger

	logs, err := c.services.Passage.GenerateLogs(passage)
	if err != nil {
		return newErrorResponse(http.StatusInternalServerError, err.Error(), c.logger)
	}

	data, err := json.Marshal(logs)
	if err != nil {
		return newErrorResponse(http.StatusInternalServerError, err.Error(), c.logger)
	}

	req, err := http.NewRequest(http.MethodPost, logsUrl, bytes.NewReader(data))
	if err != nil {
		return newErrorResponse(http.StatusInternalServerError, err.Error(), c.logger)
	}

	logger.Infof("trying to send logs: %v", logs)
	res, err := c.HttpClient.Do(req)
	if err != nil {
		return newErrorResponse(http.StatusInternalServerError, err.Error(), c.logger)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return newErrorResponse(http.StatusInternalServerError, err.Error(), c.logger)
	}

	logger.Infof("response: %v", body)

	return nil
}
