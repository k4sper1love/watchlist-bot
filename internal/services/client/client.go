package client

import (
	"bytes"
	"encoding/json"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"log"
	"net/http"
	"strconv"
	"strings"
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

	log.Println(req.Header)

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
	var headers map[string]string

	if request.HeaderType != "" {
		headers = map[string]string{
			request.HeaderType: request.HeaderValue,
		}
	}

	log.Println(request.URL, request.Body, headers)
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

func ParseErrorStatusCode(err error) int {
	errStr := err.Error()

	if strings.HasPrefix(errStr, "failed response") {
		parts := strings.Split(errStr, "code")
		codeStr := strings.TrimSpace(parts[len(parts)-1])
		code, err := strconv.Atoi(codeStr)
		if err != nil {
			return -1
		}
		return code
	}

	return -1
}
