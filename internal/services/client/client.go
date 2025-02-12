package client

import (
	"bytes"
	"encoding/json"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"net/http"
)

const (
	HeaderAuthorization  = "Authorization"
	HeaderVerification   = "Verification"
	HeaderExternalAPIKey = "X-API-KEY"
	HeaderContentType    = "Content-Type"
)

type CustomRequest struct {
	HeaderType         string
	HeaderValue        string
	Method             string
	URL                string
	Body               any
	ExpectedStatusCode int
}

func SendRequest(request *http.Request) (*http.Response, error) {
	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		utils.LogRequestError("failed to send request", err, request.Method, request.URL.String())
		return nil, err
	}

	return resp, nil
}

func SendRequestWithOptions(requestURL, method string, body any, headers map[string]string) (*http.Response, error) {
	req, err := prepareRequest(requestURL, method, body)
	if err != nil {
		utils.LogRequestError("failed to prepare request", err, method, requestURL)
		return nil, err
	}

	AddRequestHeaders(req, headers)

	resp, err := SendRequest(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
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

func Do(request *CustomRequest) (*http.Response, error) {
	headers := map[string]string{
		request.HeaderType: request.HeaderValue,
	}

	resp, err := SendRequestWithOptions(request.URL, request.Method, request.Body, headers)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != request.ExpectedStatusCode {
		defer utils.CloseBody(resp.Body)
		return nil, utils.LogResponseError(request.URL, request.Method, resp.StatusCode, resp.Status)
	}

	return resp, nil
}
