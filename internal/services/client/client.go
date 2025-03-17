// Package client provides utilities for making HTTP requests to external services.
//
// It includes functions for preparing, sending, and handling HTTP requests with customizable headers,
// body, and expected status codes.
//
// The package also supports error handling and logging for failed requests.
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
	HeaderAuthorization  = "Authorization"    // Header for authentication tokens.
	HeaderVerification   = "Verification"     // Header for verification tokens.
	HeaderExternalAPIKey = "X-API-KEY"        // Header for external API keys.
	HeaderContentType    = "Content-Type"     // Header for specifying content type.
	ContentTypeJSON      = "application/json" // Default content type for JSON requests.
)

// CustomRequest represents a custom HTTP request with additional options.
type CustomRequest struct {
	HeaderType         string // Type of header to set (e.g., Authorization, X-API-KEY).
	HeaderValue        string // Value for the specified header.
	Method             string // HTTP method (e.g., GET, POST, PUT, DELETE).
	URL                string // Target URL for the request.
	Body               any    // Request body (optional).
	ExpectedStatusCode int    // Expected HTTP status code for successful responses.
	WithoutLog         bool   // Whether to suppress logging for failed responses.
}

// SendRequest sends an HTTP request and returns the response.
func SendRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		utils.LogRequestError("failed to send request", err, req.Method, req.URL.String())
		return nil, err
	}

	return resp, nil
}

// SendRequestWithOptions sends an HTTP request with custom headers and body.
func SendRequestWithOptions(url, method string, body any, headers map[string]string) (*http.Response, error) {
	req, err := prepareRequest(url, method, body)
	if err != nil {
		utils.LogRequestError("failed to prepare request", err, method, url)
		return nil, err
	}

	setRequestHeaders(req, headers)
	return SendRequest(req)
}

// setRequestHeaders sets the headers for an HTTP request.
func setRequestHeaders(req *http.Request, headers map[string]string) {
	req.Header.Set(HeaderContentType, ContentTypeJSON)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
}

// prepareRequest prepares an HTTP request with the given URL, method, and data.
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

// Do sends a custom HTTP request and validates the response status code.
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

// ParseErrorStatusCode extracts the HTTP status code from an error message.
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
