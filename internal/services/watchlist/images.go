package watchlist

import (
	"bytes"
	"fmt"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/client"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
)

func UploadImage(app models.App, data []byte) (string, error) {
	body, headerValue, err := prepareImageParams(data)
	if err != nil {
		sl.Log.Error("failed to prepare image params", slog.Any("error", err))
		return "", err
	}

	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderContentType,
			HeaderValue:        headerValue,
			Method:             http.MethodPost,
			URL:                fmt.Sprintf("%s/upload", app.Vars.Host),
			Body:               body,
			ExpectedStatusCode: http.StatusCreated,
		},
	)
	if err != nil {
		return "", err
	}
	defer utils.CloseBody(resp.Body)

	imageURL, err := parseImageURL(resp.Body)
	if err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return "", err
	}

	return imageURL, nil
}

func prepareImageParams(data []byte) (*bytes.Buffer, string, error) {
	var body *bytes.Buffer
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("image", "image.jpg")
	if err != nil {
		return nil, "", fmt.Errorf("failed to create form file: %v", err)
	}

	_, err = io.Copy(part, bytes.NewReader(data))
	if err != nil {
		return nil, "", fmt.Errorf("failed to copy image: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return nil, "", fmt.Errorf("failed to close writer: %v", err)
	}

	return body, writer.FormDataContentType(), nil
}
