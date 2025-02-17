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
	req, err := prepareImageRequest(app, data)
	if err != nil {
		sl.Log.Error("failed to prepare request", slog.Any("error", err))
		return "", err
	}

	resp, err := client.SendRequest(req)
	if err != nil {
		return "", err
	}
	defer utils.CloseBody(resp.Body)

	if resp.StatusCode != 201 {
		return "", utils.LogResponseError(req.URL.String(), req.Method, resp.StatusCode, resp.Status)
	}

	imageURL, err := parseImageURL(resp.Body)
	if err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return "", err
	}

	return imageURL, nil
}

func prepareImageRequest(app models.App, data []byte) (*http.Request, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("image", "image.jpg")
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %v", err)
	}

	_, err = io.Copy(part, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to copy image: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close writer: %v", err)
	}

	request, err := http.NewRequest(http.MethodPost, app.Vars.Host+"/upload", body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	return body, writer.FormDataContentType(), nil
}
