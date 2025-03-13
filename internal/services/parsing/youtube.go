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

type externalVideoData struct {
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

	video, err := fetchYoutubeVideo(service, videoID)
	if err != nil {
		sl.Log.Error("failed to fetch youtube video", slog.Any("error", err), slog.String("url", url))
		return nil, err
	}

	externalData, err := getExternalVideoData(videoID)
	if err != nil {
		sl.Log.Warn("failed to get external video data", slog.Any("error", err), slog.String("videoID", videoID))
		externalData = &externalVideoData{}
	}

	return parseVideoFromYoutube(session, video, externalData), nil
}

func fetchYoutubeVideo(service *youtube.Service, videoID string) (*youtube.Video, error) {
	resp, err := service.Videos.List([]string{"snippet", "statistics", "contentDetails"}).Id(videoID).Do()
	if err != nil || len(resp.Items) == 0 {
		return nil, fmt.Errorf("video not found or error occured: %v", err)
	}
	return resp.Items[0], nil
}

func getExternalVideoData(videoID string) (*externalVideoData, error) {
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

	var data externalVideoData
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return &data, err
}

func parseVideoFromYoutube(session *models.Session, video *youtube.Video, externalData *externalVideoData) *apiModels.Film {
	return &apiModels.Film{
		Title:       video.Snippet.Title,
		Genre:       "YouTube Video",
		ImageURL:    parseThumbnailFromYoutube(video),
		Year:        parseYearFromYoutube(video.Snippet.PublishedAt),
		Rating:      utils.Round(externalData.Rating * 2),
		Description: formatDescription(session, video, externalData),
	}
}

func parseThumbnailFromYoutube(video *youtube.Video) string {
	if video.Snippet.Thumbnails.Maxres != nil {
		return video.Snippet.Thumbnails.Maxres.Url
	} else if video.Snippet.Thumbnails.High != nil {
		return video.Snippet.Thumbnails.High.Url
	}
	return ""
}

func parseYearFromYoutube(date string) int {
	if len(date) >= 4 {
		year, _ := strconv.Atoi(date[:4])
		return year
	}
	return 0
}

func formatDescription(session *models.Session, video *youtube.Video, externalData *externalVideoData) string {
	return fmt.Sprintf(
		"👨‍💼 %s: %s\n⏳ %s: %s\n👁️‍🗨️ %s: %d\n❤️ %s: %d / %d\n💬 %s: %d\n📆 %s: %s",
		translator.Translate(session.Lang, "author", nil, nil), video.Snippet.ChannelTitle,
		translator.Translate(session.Lang, "duration", nil, nil), utils.ParseISO8601Duration(video.ContentDetails.Duration),
		translator.Translate(session.Lang, "views", nil, nil), externalData.ViewCount,
		translator.Translate(session.Lang, "grades", nil, nil), externalData.Likes, externalData.Dislikes,
		translator.Translate(session.Lang, "comments", nil, nil), video.Statistics.CommentCount,
		translator.Translate(session.Lang, "dateOfRelease", nil, nil), utils.FormatTextDate(video.Snippet.PublishedAt),
	)
}
