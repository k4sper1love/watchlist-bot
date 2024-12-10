package watchlist

import (
	"bytes"
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/client"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

func UploadImage(app models.App, data []byte) (string, error) {
	request, err := prepareFileRequest(app, data)
	if err != nil {
		log.Printf("Error at 15: preparing file request: %v\n", err)
		return "", err
	}

	resp, err := client.SendRequest(request)
	if err != nil {
		log.Printf("Error at 21: sending request: %v\n", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		log.Printf("Error: server returned status %v. Response: %s\n", resp.StatusCode, string(respBody))
		return "", fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	return parseImageURL(resp.Body)
}

func prepareFileRequest(app models.App, data []byte) (*http.Request, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("image", "image.jpg")
	if err != nil {
		log.Printf("Error at 35: creating form file: %v\n", err)
		return nil, err
	}

	_, err = io.Copy(part, bytes.NewReader(data))
	if err != nil {
		log.Printf("Error at 41: copying image to form part: %v\n", err)
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		log.Printf("Error at 47: closing writer: %v\n", err)
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, app.Vars.Host+"/upload", &body)
	if err != nil {
		log.Printf("Error at 53: creating HTTP request: %v\n", err)
		return nil, err
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())

	return request, nil
}
