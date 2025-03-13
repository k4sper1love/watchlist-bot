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

const (
	categoryMovies = "movies"
	categorySeries = "series"
)

func GetFilmFromKinoafisha(url string) (*apiModels.Film, error) {
	return getMediaFromKinoafisha(url, categoryMovies, parseFilmFromKinoafisha)
}

func GetSeriesFromKinoafisha(url string) (*apiModels.Film, error) {
	return getMediaFromKinoafisha(url, categorySeries, parseSeriesFromKinoafisha)
}

func getMediaFromKinoafisha(url, category string, parser func(*apiModels.Film, io.Reader) error) (*apiModels.Film, error) {
	id := parseKinoafishaID(url)
	if id == "" {
		return nil, fmt.Errorf("invalid Kinoafisha URL: %s", url)
	}

	resp, err := client.Do(
		&client.CustomRequest{
			Method:             http.MethodGet,
			URL:                fmt.Sprintf("https://www.kinoafisha.info/%s/%s", category, id),
			ExpectedStatusCode: http.StatusOK,
		},
	)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body)

	var film apiModels.Film
	if err = parser(&film, resp.Body); err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return &film, err
}

func parseFilmFromKinoafisha(dest *apiModels.Film, data io.Reader) error {
	doc, err := goquery.NewDocumentFromReader(data)
	if err != nil {
		return err
	}

	dest.Title = strings.Split(getTextOrDefault(doc, ".newFilmInfo_title", "Unknown"), ",")[0]
	dest.Year = getKinoafishaYear(doc)
	dest.Genre = getTextOrDefault(doc, ".newFilmInfo_genreItem", "")
	dest.Description = getTextOrDefault(doc, ".more_content p", "")
	dest.Rating = parseKinoafishaRating(getTextOrDefault(doc, ".rating_imdb", "0"))
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
	dest.Genre = getTextOrDefault(doc, ".newFilmInfo_genreItem", "")
	dest.Description = getTextOrDefault(doc, ".more_content p", "")
	dest.Rating = parseKinoafishaRating(getTextOrDefault(doc, ".ratingBlockCard_externalVal", "0"))
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

func parseKinoafishaRating(ratingStr string) float64 {
	parts := strings.Split(ratingStr, ":")
	if len(parts) > 1 {
		ratingStr = parts[1]
	}
	rating, _ := strconv.ParseFloat(strings.TrimSpace(ratingStr), 64)
	return rating
}

func getKinoafishaYear(doc *goquery.Document) int {
	var year int
	doc.Find(".newFilmInfo_infoItem").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Find(".newFilmInfo_infoName").Text(), "Год выпуска") {
			year, _ = strconv.Atoi(strings.TrimSpace(s.Find(".newFilmInfo_infoData").Text()))
		}
	})

	return year
}

func parseKinoafishaID(url string) string {
	parts := strings.Split(strings.TrimPrefix(url, "https://www.kinoafisha.info/"), "/")
	if len(parts) > 0 {
		return parts[1]
	}
	return ""
}
