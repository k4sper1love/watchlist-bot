package client

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func SendRequest(request *http.Request) (*http.Response, error) {
	client := &http.Client{}
	return client.Do(request)
}

func SendRequestWithOptions(requestURL, method string, body any, headers map[string]string) (*http.Response, error) {
	req, err := prepareRequest(requestURL, method, body)
	if err != nil {
		return nil, err
	}

	AddRequestHeaders(req, headers)

	return SendRequest(req)
}

func AddRequestHeaders(req *http.Request, headers map[string]string) {
	req.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}
}

func prepareRequest(requestURL, method string, data any) (*http.Request, error) {
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

	return http.NewRequest(method, requestURL, requestBody)
}
