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

	dest.Title = client.GetTextOrDefault(doc, ".b-post__title", "Unknown")

	yearStr := client.GetTextOrDefault(doc, "a[href*='/year/']", "0")
	yearStr = strings.Replace(yearStr, " года", "", 1)
	dest.Year, _ = strconv.Atoi(yearStr)

	doc.Find(".b-post__info tr").Each(func(i int, s *goquery.Selection) {
		label := strings.TrimSpace(s.Find("td.l").Text())
		if strings.Contains(label, "Жанр") {
			genres := strings.TrimSpace(s.Find("td").Last().Text())
			dest.Genre = strings.Split(genres, ",")[0]
			return
		}
	})

	dest.Description = client.GetTextOrDefault(doc, ".b-post__description_text", "")

	ratingStr := client.GetTextOrDefault(doc, ".b-post__info_rates.imdb .bold", "0")
	dest.Rating, _ = strconv.ParseFloat(ratingStr, 64)

	dest.ImageURL = doc.Find(".b-sidecover a").AttrOr("href", "")

	return nil
}
