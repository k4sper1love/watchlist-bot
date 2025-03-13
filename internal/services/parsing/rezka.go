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

func GetFilmFromRezka(url string) (*apiModels.Film, error) {
	resp, err := client.Do(
		&client.CustomRequest{
			Method:             http.MethodGet,
			URL:                url,
			ExpectedStatusCode: http.StatusOK,
		},
	)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body)

	var film apiModels.Film
	if err = parseFilmFromRezka(&film, resp.Body); err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return &film, err
}

func parseFilmFromRezka(dest *apiModels.Film, data io.Reader) error {
	doc, err := goquery.NewDocumentFromReader(data)
	if err != nil {
		return err
	}

	dest.Title = getTextOrDefault(doc, ".b-post__title", "Unknown")
	dest.Year = parseYearFromRezka(getTextOrDefault(doc, "a[href*='/year/']", "0"))
	dest.Genre = getGenreFromRezka(doc)
	dest.Description = getTextOrDefault(doc, ".b-post__description_text", "")
	dest.Rating = parseRatingFromRezka(getTextOrDefault(doc, ".b-post__info_rates.imdb .bold", "0"))
	dest.ImageURL = doc.Find(".b-sidecover a").AttrOr("href", "")

	return nil
}

func parseYearFromRezka(yearStr string) int {
	yearStr = strings.Replace(yearStr, " года", "", 1)
	year, _ := strconv.Atoi(yearStr)
	return year
}

func parseRatingFromRezka(ratingStr string) float64 {
	rating, _ := strconv.ParseFloat(ratingStr, 64)
	return rating
}

func getGenreFromRezka(doc *goquery.Document) string {
	var genre string
	doc.Find(".b-post__info tr").Each(func(i int, s *goquery.Selection) {
		label := strings.TrimSpace(s.Find("td.l").Text())
		if strings.Contains(label, "Жанр") {
			genre = strings.Split(strings.TrimSpace(s.Find("td").Last().Text()), ",")[0]
			return
		}
	})
	return genre
}
