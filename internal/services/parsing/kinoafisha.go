package parsing

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/client"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func GetFilmFromKinoafisha(url string) (*apiModels.Film, error) {
	resp, err := getDataFromKinoafisha(url)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body)

	var film apiModels.Film
	if err = parseFilmFromKinoafisha(&film, resp.Body); err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return &film, nil
}

func GetSeriesFromKinoafisha(url string) (*apiModels.Film, error) {
	resp, err := getDataFromKinoafisha(url)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body)

	var film apiModels.Film
	if err = parseSeriesFromKinoafisha(&film, resp.Body); err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return &film, err
}

func getDataFromKinoafisha(url string) (*http.Response, error) {
	return client.Do(
		&client.CustomRequest{
			Method:             http.MethodGet,
			URL:                url,
			ExpectedStatusCode: http.StatusOK,
		},
	)
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
