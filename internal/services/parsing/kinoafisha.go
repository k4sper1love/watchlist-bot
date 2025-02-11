package parsing

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/client"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

func GetFilmFromKinoafisha(url string) (*apiModels.Film, error) {
	body, err := getDataFromKinoafisha(url)
	defer body.Close()
	if err != nil {
		return nil, err
	}

	var film apiModels.Film
	if err = parseFilmFromKinoafisha(&film, body); err != nil {
		sl.Log.Error("failed to parse film from Kinoafisha", slog.Any("error", err), slog.String("url", url))
		return nil, err
	}

	return &film, nil
}

func GetSeriesFromKinoafisha(url string) (*apiModels.Film, error) {
	body, err := getDataFromKinoafisha(url)
	defer body.Close()
	if err != nil {
		return nil, err
	}

	var film apiModels.Film
	if err = parseSeriesFromKinoafisha(&film, body); err != nil {
		slog.Error("failed to parse series from Kinoafisha", slog.Any("error", err), slog.String("url", url))
		return nil, err
	}

	return &film, err
}

func getDataFromKinoafisha(url string) (io.ReadCloser, error) {
	resp, err := client.SendRequestWithOptions(url, http.MethodGet, nil, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, client.LogResponseError(url, resp.StatusCode, resp.Status)
	}

	return resp.Body, nil
}

func parseFilmFromKinoafisha(dest *apiModels.Film, data io.Reader) error {
	doc, err := goquery.NewDocumentFromReader(data)
	if err != nil {
		return err
	}

	title := doc.Find(".trailer_title").Text()
	title = strings.Split(title, ",")[0]
	dest.Title = strings.TrimSpace(title)

	year := strings.TrimSpace(doc.Find(".trailer_year").Text())
	year = strings.Split(year, "/")[0]
	dest.Year, err = strconv.Atoi(strings.TrimSpace(year))
	if err != nil {
		return fmt.Errorf("failed to parse year: %v", err)
	}

	genre := doc.Find(".filmInfo_genreItem").First().Text()
	dest.Genre = strings.TrimSpace(genre)

	description := doc.Find(".filmDesc_editor").First().Text()
	dest.Description = strings.TrimSpace(description)

	rating := doc.Find(".rating_num").Text()
	dest.Rating, err = strconv.ParseFloat(strings.TrimSpace(rating), 64)
	if err != nil {
		return fmt.Errorf("failed to parse rating: %v", err)
	}

	dest.ImageURL = doc.Find(".filmInfo_posterLink").First().AttrOr("href", "")

	return nil
}

func parseSeriesFromKinoafisha(dest *apiModels.Film, data io.Reader) error {
	doc, err := goquery.NewDocumentFromReader(data)
	if err != nil {
		return nil
	}

	title := doc.Find(".newFilmInfo_title").Text()
	title = strings.Split(title, "(")[0]
	dest.Title = strings.TrimSpace(title)

	year := strings.TrimSpace(doc.Find(".newFilmInfo_infoData").First().Text())
	dest.Year, err = strconv.Atoi(strings.TrimSpace(year))
	if err != nil {
		return fmt.Errorf("failed to parse year: %v", err)
	}

	genre := doc.Find(".newFilmInfo_genreItem ").First().Text()
	dest.Genre = strings.TrimSpace(genre)

	description := doc.Find(".newFilmInfo_descEditor").First().Text()
	dest.Description = strings.TrimSpace(description)

	rating := doc.Find(".ratingBlockCard_externalVal").First().Text()
	dest.Rating, err = strconv.ParseFloat(strings.TrimSpace(rating), 64)
	if err != nil {
		return fmt.Errorf("failed to parse rating: %v", err)
	}

	imageData := doc.Find(".newFilmInfo_posterSlide").AttrOr("data-fullscreengallery-item", "")
	imageData = strings.Replace(imageData, "\\", "", -1)

	var jsonData map[string]string
	if err = json.Unmarshal([]byte(imageData), &jsonData); err != nil {
		return fmt.Errorf("failed to parse image data: %v", err)
	}
	dest.ImageURL = strings.TrimSpace(jsonData["image"])

	return nil
}
