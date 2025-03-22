package parsing

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/client"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	categoryMovies = "movies" // Category for movies on Kinoafisha.
	categorySeries = "series" // Category for series on Kinoafisha.
)

// GetFilmFromKinoafisha fetches film details from the Kinoafisha website using the provided URL.
func GetFilmFromKinoafisha(session *models.Session, url string) (*apiModels.Film, error) {
	return getMediaFromKinoafisha(session, url, categoryMovies, parseFilmFromKinoafisha)
}

// GetSeriesFromKinoafisha fetches series details from the Kinoafisha website using the provided URL.
func GetSeriesFromKinoafisha(session *models.Session, url string) (*apiModels.Film, error) {
	return getMediaFromKinoafisha(session, url, categorySeries, parseSeriesFromKinoafisha)
}

// getMediaFromKinoafisha is a helper function to fetch media (film or series) details from Kinoafisha.
// It uses the provided URL, category, and parser function to extract and parse the data.
func getMediaFromKinoafisha(session *models.Session, url, category string, parser func(*apiModels.Film, io.Reader) error) (*apiModels.Film, error) {
	id, err := parseKinoafishaID(url)
	if err != nil {
		utils.LogParseFromURLError(session.TelegramID, "failed to parse ID", err, url)
		return nil, err
	}

	resp, err := client.Do(
		&client.CustomRequest{
			Method:             http.MethodGet, // HTTP GET method for fetching data.
			URL:                fmt.Sprintf("https://www.kinoafisha.info/%s/%s", category, id),
			ExpectedStatusCode: http.StatusOK, // Expecting a 200 OK response.
			TelegramID:         session.TelegramID,
		},
	)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body) // Ensure the response body is closed after use.

	var film apiModels.Film
	if err = parser(&film, resp.Body); err != nil {
		utils.LogParseJSONError(session.TelegramID, err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return &film, err
}

// parseFilmFromKinoafisha parses film details from the Kinoafisha HTML document into an `models.Film` object.
func parseFilmFromKinoafisha(dest *apiModels.Film, data io.Reader) error {
	doc, err := goquery.NewDocumentFromReader(data)
	if err != nil {
		return err
	}

	// Extract film details from the HTML document.
	dest.Title = strings.Split(getTextOrDefault(doc, ".newFilmInfo_title", "Unknown"), ",")[0] // Extract the title.
	dest.Year = getKinoafishaYear(doc)                                                         // Extract the release year.
	dest.Genre = getTextOrDefault(doc, ".newFilmInfo_genreItem", "")                           // Extract the genre.
	dest.Description = getTextOrDefault(doc, ".more_content p", "")                            // Extract the description.
	dest.Rating = parseKinoafishaRating(getTextOrDefault(doc, ".rating_imdb", "0"))            // Extract and parse the rating.
	dest.ImageURL = getKinoafishaImageURL(doc)                                                 // Extract the image URL.

	return nil
}

// parseSeriesFromKinoafisha parses series details from the Kinoafisha HTML document into an `models.Film` object.
func parseSeriesFromKinoafisha(dest *apiModels.Film, data io.Reader) error {
	doc, err := goquery.NewDocumentFromReader(data)
	if err != nil {
		return err
	}

	// Extract series title from breadcrumbs.
	doc.Find(".newFilmInfo_breadcrumbs .breadcrumbs_item").Each(func(i int, s *goquery.Selection) {
		if i == 2 { // The third breadcrumb item typically contains the series title.
			dest.Title = strings.TrimSpace(s.Text())
		}
	})

	// Extract other series details.
	dest.Year = getKinoafishaYear(doc)                                                              // Extract the release year.
	dest.Genre = getTextOrDefault(doc, ".newFilmInfo_genreItem", "")                                // Extract the genre.
	dest.Description = getTextOrDefault(doc, ".more_content p", "")                                 // Extract the description.
	dest.Rating = parseKinoafishaRating(getTextOrDefault(doc, ".ratingBlockCard_externalVal", "0")) // Extract and parse the rating.
	dest.ImageURL = getKinoafishaImageURL(doc)                                                      // Extract the image URL.

	return nil
}

// getKinoafishaImageURL extracts the image URL from the Kinoafisha HTML document.
func getKinoafishaImageURL(doc *goquery.Document) string {
	imageData := doc.Find(".newFilmInfo_posterSlide").AttrOr("data-fullscreengallery-item", "")
	imageData = strings.Replace(imageData, "\\", "", -1) // Remove escape characters.

	var jsonData map[string]string
	if err := json.Unmarshal([]byte(imageData), &jsonData); err == nil {
		return strings.TrimSpace(jsonData["image"])
	}

	return ""
}

// parseKinoafishaRating parses the rating value from a string.
func parseKinoafishaRating(ratingStr string) float64 {
	parts := strings.Split(ratingStr, ":")
	if len(parts) > 1 {
		ratingStr = parts[1]
	}
	rating, _ := strconv.ParseFloat(strings.TrimSpace(ratingStr), 64)
	return rating
}

// getKinoafishaYear extracts the release year from the Kinoafisha HTML document.
func getKinoafishaYear(doc *goquery.Document) int {
	var year int
	doc.Find(".newFilmInfo_infoItem").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Find(".newFilmInfo_infoName").Text(), "Год выпуска") {
			year, _ = strconv.Atoi(strings.TrimSpace(s.Find(".newFilmInfo_infoData").Text()))
		}
	})

	return year
}

// parseKinoafishaID extracts the media ID from the Kinoafisha URL.
func parseKinoafishaID(url string) (string, error) {
	parts := strings.Split(strings.TrimPrefix(url, "https://www.kinoafisha.info/"), "/")
	if len(parts) > 0 {
		return parts[1], nil // The second part of the path is the media ID.
	}
	return "", fmt.Errorf("invalid Kinoafisha URL: %s", url)
}
