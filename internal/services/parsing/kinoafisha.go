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
	url = fmt.Sprintf("https://www.kinoafisha.info/movies/%s", parseKinoafishaID(url))

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
	url = fmt.Sprintf("https://www.kinoafisha.info/series/%s", parseKinoafishaID(url))

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

	dest.Title = client.GetTextOrDefault(doc, ".newFilmInfo_title", "Unknown")
	dest.Title = strings.Split(dest.Title, ",")[0]

	dest.Year = getKinoafishaYear(doc)

	dest.Genre = client.GetTextOrDefault(doc, ".newFilmInfo_genreItem", "")
	dest.Description = client.GetTextOrDefault(doc, ".more_content p", "")

	ratingStr := client.GetTextOrDefault(doc, ".rating_imdb ", "0")
	parts := strings.Split(ratingStr, ":")
	if len(parts) > 1 {
		ratingStr = parts[1]
	} else {
		ratingStr = "0"
	}
	dest.Rating, _ = strconv.ParseFloat(strings.TrimSpace(ratingStr), 64)

	dest.ImageURL = getKinoafishaImageURL(doc)

	return nil
}

func parseSeriesFromKinoafisha(dest *apiModels.Film, data io.Reader) error {
	doc, err := goquery.NewDocumentFromReader(data)
	if err != nil {
		return err
	}

	doc.Find(".newFilmInfo_breadcrumbs .breadcrumbs_item").Each(func(i int, s *goquery.Selection) {
		if i == 2 {
			dest.Title = strings.TrimSpace(s.Text())
		}
	})

	dest.Year = getKinoafishaYear(doc)

	dest.Genre = client.GetTextOrDefault(doc, ".newFilmInfo_genreItem", "")
	dest.Description = client.GetTextOrDefault(doc, ".more_content p", "")

	ratingStr := client.GetTextOrDefault(doc, ".ratingBlockCard_externalVal", "0")
	dest.Rating, _ = strconv.ParseFloat(strings.TrimSpace(ratingStr), 64)

	dest.ImageURL = getKinoafishaImageURL(doc)

	return nil
}

func getKinoafishaImageURL(doc *goquery.Document) string {
	imageData := doc.Find(".newFilmInfo_posterSlide").AttrOr("data-fullscreengallery-item", "")
	imageData = strings.Replace(imageData, "\\", "", -1)

	var jsonData map[string]string
	if err := json.Unmarshal([]byte(imageData), &jsonData); err == nil {
		return strings.TrimSpace(jsonData["image"])
	}

	return ""
}

func getKinoafishaYear(doc *goquery.Document) int {
	var year int
	doc.Find(".newFilmInfo_infoItem").Each(func(i int, s *goquery.Selection) {
		name := s.Find(".newFilmInfo_infoName").Text()
		if strings.Contains(name, "Год выпуска") {
			yearStr := s.Find(".newFilmInfo_infoData").Text()
			yearStr = strings.TrimSpace(yearStr)
			year, _ = strconv.Atoi(yearStr)
		}
	})

	return year
}

func parseKinoafishaID(url string) string {
	shortUrl := strings.TrimPrefix(url, "https://www.kinoafisha.info/")
	parts := strings.Split(shortUrl, "/")

	if len(parts) > 0 {
		return parts[1]
	}

	return ""
}
