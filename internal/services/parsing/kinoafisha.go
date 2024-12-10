package parsing

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/client"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func GetFilmFromKinoafisha(url string) (*apiModels.Film, error) {
	body, err := getDataFromKinoafisha(url)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var film apiModels.Film
	if err := parseFilmFromKinoafisha(&film, body); err != nil {
		return nil, err
	}

	return &film, nil
}

func GetSeriesFromKinoafisha(url string) (*apiModels.Film, error) {
	body, err := getDataFromKinoafisha(url)
	if err != nil {
		return nil, err
	}

	var film apiModels.Film
	if err := parseSeriesFromKinoafisha(&film, body); err != nil {
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
		resp.Body.Close()
		return nil, fmt.Errorf("failed response. Status is %s", resp.Status)
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
		return err
	}

	genre := doc.Find(".filmInfo_genreItem").First().Text()
	dest.Genre = strings.TrimSpace(genre)

	description := doc.Find(".filmDesc_editor").First().Text()
	dest.Description = strings.TrimSpace(description)

	rating := doc.Find(".rating_num").Text()
	dest.Rating, err = strconv.ParseFloat(strings.TrimSpace(rating), 64)
	if err != nil {
		return err
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
		return err
	}

	genre := doc.Find(".newFilmInfo_genreItem ").First().Text()
	dest.Genre = strings.TrimSpace(genre)

	description := doc.Find(".newFilmInfo_descEditor").First().Text()
	dest.Description = strings.TrimSpace(description)

	rating := doc.Find(".ratingBlockCard_externalVal").First().Text()
	dest.Rating, err = strconv.ParseFloat(strings.TrimSpace(rating), 64)
	if err != nil {
		return err
	}

	imageData := doc.Find(".newFilmInfo_posterSlide").AttrOr("data-fullscreengallery-item", "")
	imageData = strings.Replace(imageData, "\\", "", -1)

	var jsonData map[string]string
	if err := json.Unmarshal([]byte(imageData), &jsonData); err != nil {
		dest.ImageURL = ""
	}
	dest.ImageURL = strings.TrimSpace(jsonData["image"])

	return nil
}
