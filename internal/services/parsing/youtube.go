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

	service, err := youtube.NewService(context.Background(), option.WithAPIKey(app.Config.YoutubeAPIToken))
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
		sl.Log.Warn("failed to get external video data", slog.Any("error", err), slog.String("videoID", videoID))
		externalData = &ExternalVideoData{}
	}

	var film apiModels.Film
	parseVideoFromYoutube(&film, session, video, externalData)

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

func parseVideoFromYoutube(dest *apiModels.Film, session *models.Session, video *youtube.Video, externalData *ExternalVideoData) {
	dest.Title = video.Snippet.Title
	dest.Genre = "YouTube Video"

	if video.Snippet.Thumbnails.Maxres != nil {
		dest.ImageURL = video.Snippet.Thumbnails.Maxres.Url
	} else if video.Snippet.Thumbnails.High != nil {
		dest.ImageURL = video.Snippet.Thumbnails.High.Url
	} else {
		dest.ImageURL = ""
	}

	if len(video.Snippet.PublishedAt) >= 4 {
		dest.Year, _ = strconv.Atoi(video.Snippet.PublishedAt[:4])
	} else {
		dest.Year = 0
	}

	dest.Rating, _ = utils.Round(externalData.Rating * 2)

	duration := video.ContentDetails.Duration
	parsedDuration, _ := utils.ParseISO8601Duration(duration)

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
}
