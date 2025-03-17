package parsing

import (
	"github.com/PuerkitoBio/goquery"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/client"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// GetFilmFromRezka fetches film details from the Rezka website using the provided URL.
// It sends an HTTP GET request to the URL and parses the HTML response into an `models.Film` object.
func GetFilmFromRezka(url string) (*apiModels.Film, error) {
	resp, err := client.Do(
		&client.CustomRequest{
			Method:             http.MethodGet, // HTTP GET method for fetching data.
			URL:                url,
			ExpectedStatusCode: http.StatusOK, // Expecting a 200 OK response.
		},
	)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body) // Ensure the response body is closed after use.

	var film apiModels.Film
	if err = parseFilmFromRezka(&film, resp.Body); err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return &film, err
}

// parseFilmFromRezka parses film details from the Rezka HTML document into an `models.Film` object.
func parseFilmFromRezka(dest *apiModels.Film, data io.Reader) error {
	doc, err := goquery.NewDocumentFromReader(data)
	if err != nil {
		return err
	}

	// Extract film details from the HTML document.
	dest.Title = getTextOrDefault(doc, ".b-post__title", "Unknown")                                  // Default title is "Unknown".
	dest.Year = parseYearFromRezka(getTextOrDefault(doc, "a[href*='/year/']", "0"))                  // Extract and parse the year.
	dest.Genre = getGenreFromRezka(doc)                                                              // Extract the first genre from the list.
	dest.Description = getTextOrDefault(doc, ".b-post__description_text", "")                        // Default description is an empty string.
	dest.Rating = parseRatingFromRezka(getTextOrDefault(doc, ".b-post__info_rates.imdb .bold", "0")) // Parse the rating.
	dest.ImageURL = doc.Find(".b-sidecover a").AttrOr("href", "")                                    // Extract the image URL.

	return nil
}

// parseYearFromRezka parses the year value from a string.
func parseYearFromRezka(yearStr string) int {
	yearStr = strings.Replace(yearStr, " года", "", 1) // Remove " года" suffix.
	year, _ := strconv.Atoi(yearStr)                   // Convert to integer.
	return year
}

// parseRatingFromRezka parses the rating value from a string.
func parseRatingFromRezka(ratingStr string) float64 {
	rating, _ := strconv.ParseFloat(ratingStr, 64) // Convert to float64.
	return rating
}

// getGenreFromRezka extracts the first genre from the Rezka HTML document.
func getGenreFromRezka(doc *goquery.Document) string {
	var genre string
	doc.Find(".b-post__info tr").Each(func(i int, s *goquery.Selection) {
		label := strings.TrimSpace(s.Find("td.l").Text()) // Extract the label text.
		if strings.Contains(label, "Жанр") {              // Check if the label contains "Жанр".
			genre = strings.Split(strings.TrimSpace(s.Find("td").Last().Text()), ",")[0] // Extract the first genre.
			return
		}
	})
	return genre
}
