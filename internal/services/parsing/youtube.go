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

// externalVideoData represents additional video data fetched from an external API.
type externalVideoData struct {
	ID          string  `json:"id"`          // Video ID.
	DateCreated string  `json:"dateCreated"` // Date the video was created.
	Likes       int64   `json:"likes"`       // Number of likes.
	RawDislikes int64   `json:"rawDislikes"` // Raw number of dislikes.
	RawLikes    int64   `json:"rawLikes"`    // Raw number of likes.
	Dislikes    int64   `json:"dislikes"`    // Number of dislikes.
	Rating      float64 `json:"rating"`      // Video rating.
	ViewCount   int64   `json:"viewCount"`   // Number of views.
	Deleted     bool    `json:"deleted"`     // Indicates if the video has been deleted.
}

// GetFilmFromYoutube fetches a YouTube video and parses it into an `models.Film` object.
// It extracts the video ID from the URL, fetches video details from the YouTube API,
// and retrieves additional data from an external API.
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

// fetchYoutubeVideo fetches video details from the YouTube API using the provided video ID.
func fetchYoutubeVideo(service *youtube.Service, videoID string) (*youtube.Video, error) {
	resp, err := service.Videos.List([]string{"snippet", "statistics", "contentDetails"}).Id(videoID).Do()
	if err != nil || len(resp.Items) == 0 {
		return nil, fmt.Errorf("video not found or error occured: %v", err)
	}
	return resp.Items[0], nil
}

// getExternalVideoData fetches additional video data (e.g., likes, dislikes) from an external API.
func getExternalVideoData(videoID string) (*externalVideoData, error) {
	resp, err := client.Do(
		&client.CustomRequest{
			Method:             http.MethodGet, // HTTP GET method for fetching data.
			URL:                fmt.Sprintf("https://returnyoutubedislikeapi.com/votes?videoId=%s", videoID),
			ExpectedStatusCode: http.StatusOK, // Expecting a 200 OK response.
		},
	)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body) // Ensure the response body is closed after use.

	var data externalVideoData
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return &data, err
}

// parseVideoFromYoutube parses a YouTube video into an `models.Film` object.
func parseVideoFromYoutube(session *models.Session, video *youtube.Video, externalData *externalVideoData) *apiModels.Film {
	return &apiModels.Film{
		Title:       video.Snippet.Title,
		Genre:       "YouTube Video", // Default genre for YouTube videos.
		ImageURL:    parseThumbnailFromYoutube(video),
		Year:        parseYearFromYoutube(video.Snippet.PublishedAt),
		Rating:      utils.Round(externalData.Rating * 2), // Scale the rating to match the application's format.
		Description: formatDescription(session, video, externalData),
	}
}

// parseThumbnailFromYoutube extracts the highest-resolution thumbnail URL from the YouTube video.
func parseThumbnailFromYoutube(video *youtube.Video) string {
	if video.Snippet.Thumbnails.Maxres != nil {
		return video.Snippet.Thumbnails.Maxres.Url
	} else if video.Snippet.Thumbnails.High != nil {
		return video.Snippet.Thumbnails.High.Url
	}
	return ""
}

// parseYearFromYoutube extracts the year from the video's publication date.
func parseYearFromYoutube(date string) int {
	if len(date) >= 4 {
		year, _ := strconv.Atoi(date[:4]) // Extract the first 4 characters as the year.
		return year
	}
	return 0
}

// formatDescription formats the video description with localized labels and statistics.
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
