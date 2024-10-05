package watchlist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func SendRequest(baseURL, requestURL, method, token string, data any) (*http.Response, error) {
	req, err := prepareRequest(baseURL, requestURL, method, data)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	return doRequest(req)
}

func prepareRequest(baseURL, requestURL, method string, data any) (*http.Request, error) {
	var requestBody *bytes.Buffer

	if data != nil {
		body, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		requestBody = bytes.NewBuffer(body)
	} else {
		requestBody = bytes.NewBuffer([]byte{})
	}

	return http.NewRequest(method, baseURL+requestURL, requestBody)
}

func doRequest(request *http.Request) (*http.Response, error) {
	client := &http.Client{}
	return client.Do(request)
}
