package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"net/http"
	"strconv"
	"strings"
)

const (
	HeaderAuthorization  = "Authorization"
	HeaderVerification   = "Verification"
	HeaderExternalAPIKey = "X-API-KEY"
	HeaderContentType    = "Content-Type"
	ContentTypeJSON      = "application/json"
)

type CustomRequest struct {
	HeaderType         string
	HeaderValue        string
	Method             string
	URL                string
	Body               any
	ExpectedStatusCode int
	WithoutLog         bool
}

func SendRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		utils.LogRequestError("failed to send request", err, req.Method, req.URL.String())
		return nil, err
	}

	return resp, nil
}

func SendRequestWithOptions(url, method string, body any, headers map[string]string) (*http.Response, error) {
	req, err := prepareRequest(url, method, body)
	if err != nil {
		utils.LogRequestError("failed to prepare request", err, method, url)
		return nil, err
	}

	setRequestHeaders(req, headers)
	return SendRequest(req)
}

func setRequestHeaders(req *http.Request, headers map[string]string) {
	req.Header.Set(HeaderContentType, ContentTypeJSON)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
}

func prepareRequest(url, method string, data any) (*http.Request, error) {
	requestBody := &bytes.Buffer{}

	if data != nil {
		body, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		requestBody = bytes.NewBuffer(body)
	}
	return http.NewRequest(method, url, requestBody)
}

func Do(req *CustomRequest) (*http.Response, error) {
	headers := make(map[string]string)
	if req.HeaderType != "" {
		headers[req.HeaderType] = req.HeaderValue
	}

	resp, err := SendRequestWithOptions(req.URL, req.Method, req.Body, headers)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != req.ExpectedStatusCode {
		defer utils.CloseBody(resp.Body)
		if req.WithoutLog {
			return nil, fmt.Errorf("failed response")
		}
		return nil, utils.LogResponseError(req.URL, req.Method, resp.StatusCode, resp.Status)
	}

	return resp, nil
}

func ParseErrorStatusCode(err error) int {
	if !strings.HasPrefix(err.Error(), "failed response") {
		return -1
	}

	parts := strings.Split(err.Error(), "code")
	if len(parts) < 2 {
		return -1
	}

	code, err := strconv.Atoi(strings.TrimSpace(parts[len(parts)-1]))
	if err != nil {
		return -1
	}
	return code
}
