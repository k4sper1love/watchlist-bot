package parsing

import (
	"github.com/PuerkitoBio/goquery"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/client"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GetFilmFromRezka(url string) (*apiModels.Film, error) {
	resp, err := client.SendRequestWithOptions(url, http.MethodGet, nil, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println(resp.Status)
		return nil, err
	}

	var film apiModels.Film
	if err := parseFilmFromRezka(&film, resp.Body); err != nil {
		log.Println(err)
		return nil, err
	}

	return &film, err
}

func parseFilmFromRezka(dest *apiModels.Film, data io.Reader) error {
	doc, err := goquery.NewDocumentFromReader(data)
	if err != nil {
		return nil
	}

	dest.Title = strings.TrimSpace(doc.Find(".b-post__title").Text())

	year := strings.TrimSpace(doc.Find("a[href*='/year/']").Text())
	year = strings.Replace(year, " года", "", 1)
	dest.Year, _ = strconv.Atoi(year)

	doc.Find(".b-post__info tr").Each(func(i int, s *goquery.Selection) {
		label := strings.TrimSpace(s.Find("td.l").Text())
		if strings.Contains(label, "Жанр") {
			genres := strings.TrimSpace(s.Find("td").Last().Text())
			dest.Genre = strings.Split(genres, ",")[0]
		}
	})

	dest.Description = strings.TrimSpace(doc.Find(".b-post__description_text").Text())

	rating := strings.TrimSpace(doc.Find(".imbd .bold").Text())

	dest.Rating, _ = strconv.ParseFloat(rating, 64)

	dest.ImageURL = strings.TrimSpace(doc.Find(".b-sidecover a").AttrOr("href", ""))

	log.Printf("Fetched film: %+v\n", dest)
	return nil
}
