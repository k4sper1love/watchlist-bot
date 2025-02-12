package parsing

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/client"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"log/slog"
	"net/http"
	"strconv"
)

type ExternalVideoData struct {
	ID          string  `json:"id"`
	DateCreated string  `json:"dateCreated"`
	Likes       int64   `json:"likes"`
	RawDislikes int64   `json:"rawDislikes"`
	RawLikes    int64   `json:"rawLikes"`
	Dislikes    int64   `json:"dislikes"`
	Rating      float64 `json:"rating"`
	ViewCount   int64   `json:"viewCount"`
	Deleted     bool    `json:"deleted"`
}

func GetFilmFromYoutube(app models.App, session *models.Session, url string) (*apiModels.Film, error) {
	videoID, err := utils.ExtractYoutubeVideoID(url)
	if err != nil {
		sl.Log.Error("failed to extract video ID", slog.Any("error", err), slog.String("url", url))
		return nil, err
	}

	service, err := youtube.NewService(context.Background(), option.WithAPIKey(app.Vars.YoutubeAPIToken))
	if err != nil {
		sl.Log.Error("failed to create youtube service", slog.Any("error", err))
		return nil, err
	}

	request := service.Videos.List([]string{"snippet", "statistics", "contentDetails"}).Id(videoID)
	resp, err := request.Do()
	if err != nil {
		sl.Log.Error("failed to do request", slog.Any("error", err), slog.String("url", url))
		return nil, err
	}

	if len(resp.Items) == 0 {
		return nil, fmt.Errorf("video not found")
	}

	video := resp.Items[0]

	externalData, err := getExternalVideoData(videoID)
	if err != nil {
		return nil, err
	}

	var film apiModels.Film
	if err = parseVideoFromYoutube(&film, session, video, externalData); err != nil {
		utils.LogParseJSONError(err, http.MethodGet, url)
		return nil, err
	}

	return &film, err
}

func getExternalVideoData(videoID string) (*ExternalVideoData, error) {
	resp, err := client.Do(
		&client.CustomRequest{
			Method:             http.MethodGet,
			URL:                fmt.Sprintf("https://returnyoutubedislikeapi.com/votes?videoId=%s", videoID),
			ExpectedStatusCode: http.StatusOK,
		},
	)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body)

	var data ExternalVideoData
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return &data, err
}

func parseVideoFromYoutube(dest *apiModels.Film, session *models.Session, video *youtube.Video, externalData *ExternalVideoData) error {
	var err error

	dest.Title = video.Snippet.Title
	dest.Genre = "YouTube Video"
	dest.ImageURL = video.Snippet.Thumbnails.Maxres.Url

	dest.Year, err = strconv.Atoi(video.Snippet.PublishedAt[:4])
	if err != nil {
		return fmt.Errorf("failed to parse year: %v", err)
	}

	dest.Rating, err = utils.Round(externalData.Rating * 2)
	if err != nil {
		return fmt.Errorf("failed to parse rating: %v", err)
	}

	duration := video.ContentDetails.Duration
	parsedDuration, err := utils.ParseISO8601Duration(duration)
	if err != nil {
		return fmt.Errorf("failed to parse duration: %v", err)
	}

	part1 := translator.Translate(session.Lang, "author", nil, nil)
	part2 := translator.Translate(session.Lang, "duration", nil, nil)
	part3 := translator.Translate(session.Lang, "views", nil, nil)
	part4 := translator.Translate(session.Lang, "grades", nil, nil)
	part5 := translator.Translate(session.Lang, "comments", nil, nil)
	part6 := translator.Translate(session.Lang, "dateOfRelease", nil, nil)

	dest.Description = fmt.Sprintf(
		"👨‍💼 %s: %s\n⏳ %s: %s\n👁️‍🗨️ %s: %d\n❤️ %s: %d / %d\n💬 %s: %d\n📆 %s: %s",
		part1, video.Snippet.ChannelTitle,
		part2, parsedDuration,
		part3, externalData.ViewCount,
		part4, externalData.Likes, externalData.Dislikes,
		part5, video.Statistics.CommentCount,
		part6, utils.FormatTextDate(video.Snippet.PublishedAt),
	)

	return nil
}
