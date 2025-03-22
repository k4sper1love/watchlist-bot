package watchlist

import (
	"bytes"
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/client"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
)

// UploadImage uploads an image to the API and returns the URL of the uploaded image.
// It prepares a multipart/form-data request, sends it to the API, and parses the response.
func UploadImage(app models.App, data []byte) (string, error) {
	// Prepare the multipart/form-data request for uploading the image.
	req, err := prepareImageRequest(app, data)
	if err != nil {
		slog.Error("failed to prepare request", slog.Any("error", err))
		return "", err
	}

	// Send the request to the API.
	resp, err := client.SendRequest(int(app.GetChatID()), req)
	if err != nil {
		return "", err
	}
	defer utils.CloseBody(resp.Body) // Ensure the response body is closed after use.

	// Check if the response status code indicates success.
	if resp.StatusCode != 201 {
		return "", utils.LogResponseError(int(app.GetChatID()), req.URL.String(), req.Method, 201, resp.StatusCode, resp.Status)
	}

	// Parse the response to extract the image URL.
	imageURL, err := parseImageURL(resp.Body)
	if err != nil {
		utils.LogParseJSONError(int(app.GetChatID()), err, resp.Request.Method, resp.Request.URL.String())
		return "", err
	}

	return imageURL, nil
}

// prepareImageRequest prepares a multipart/form-data HTTP request for uploading an image.
// It creates a form file with the provided image data and sets the appropriate headers.
func prepareImageRequest(app models.App, data []byte) (*http.Request, error) {
	body := new(bytes.Buffer) // Buffer to hold the multipart/form-data body.
	writer := multipart.NewWriter(body)

	// Create a form file field named "image" with the filename "image.jpg".
	part, err := writer.CreateFormFile("image", "image.jpg")
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %v", err)
	}

	// Copy the image data into the form file.
	_, err = io.Copy(part, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to copy image: %v", err)
	}

	// Close the multipart writer to finalize the body.
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close writer: %v", err)
	}

	// Create the HTTP POST request with the multipart/form-data body.
	request, err := http.NewRequest(http.MethodPost, app.Config.APIHost+"/upload", body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set the Content-Type header to the multipart/form-data boundary.
	request.Header.Set("Content-Type", writer.FormDataContentType())

	return request, nil
}
