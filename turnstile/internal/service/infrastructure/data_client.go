package infrastructure

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"turnstile/internal/models"
	"turnstile/pkg/client"
	"turnstile/pkg/logging"
)

type DataClient struct {
	logger logging.Logger
	client *client.Client
	url    string
}

func NewDataClient(logger logging.Logger, client *client.Client, url string) *DataClient {
	return &DataClient{
		logger: logger,
		client: client,
		url:    url,
	}
}

func (dc *DataClient) GetData(rv uint64) (map[string][]models.Employee, error) {
	req, err := http.NewRequest(http.MethodGet, dc.url, nil)
	q := req.URL.Query()
	q.Add("rv", strconv.FormatUint(rv, 10))
	req.URL.RawQuery = q.Encode()

	res, err := dc.client.HttpClient.Do(req)
	if res == nil {
		return nil, errors.New("DataService is unreachable")
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var employees map[string][]models.Employee

	err = json.Unmarshal(body, &employees)
	if err != nil {
		return nil, err
	}

	return employees, nil
}
