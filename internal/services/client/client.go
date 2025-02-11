package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"log/slog"
	"net/http"
)

func SendRequest(request *http.Request) (*http.Response, error) {
	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		sl.Log.Error("failed to send request", slog.Any("error", err), slog.Any("request", request))
		return nil, err
	}

	return resp, nil
}

func SendRequestWithOptions(requestURL, method string, body any, headers map[string]string) (*http.Response, error) {
	req, err := prepareRequest(requestURL, method, body)
	if err != nil {
		LogRequestError("failed to prepare request", err, method, requestURL, body, headers)
		return nil, err
	}

	AddRequestHeaders(req, headers)

	resp, err := SendRequest(req)
	if err != nil {
		LogRequestError("failed to send request with options", err, method, requestURL, body, headers)
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

func LogRequestError(message string, err error, method, requestURL string, body any, headers map[string]string) {
	sl.Log.Error(
		message,
		slog.Any("error", err),
		slog.String("method", method),
		slog.String("url", requestURL),
		slog.Any("body", body),
		slog.Any("headers", headers),
	)
}

func LogResponseError(url string, code int, status string) error {
	sl.Log.Error(
		"failed response",
		slog.String("url", url),
		slog.Int("code", code),
		slog.String("status", status),
	)
	return fmt.Errorf("failed response from %s with code %d", url, code)
}
