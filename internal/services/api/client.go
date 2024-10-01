package api

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func SendRequest(baseURL, requestURL, method string, data any) (*http.Response, error) {
	requestBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, baseURL+requestURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return client.Do(req)
}
