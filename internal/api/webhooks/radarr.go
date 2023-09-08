package webhooks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"go-subgen/pkg"
)

type WebhookMovie struct {
	Id          int32
	Title       string
	Year        int32
	FilePath    string
	ReleaseDate string
	FolderPath  string
	TmdbId      int32
	ImdbId      string
}

type WebhookRemoteMovie struct {
	TmdbId int32
	ImdbId string
	Title  string
	Year   int32
}

type WebhookMovieFile struct {
	Id             int32
	RelativePath   string
	Path           string
	Quality        string
	QualityVersion int32
	ReleaseGroup   string
	SceneName      string
	IndexerFlags   string
	Size           int64
	DateAdded      string
	MediaInfo      WebhookMovieFileMediaInfo
}

type WebhookMovieFileMediaInfo struct {
	// Original is a decimal, but go doesn't play nice with unmarshalling that.
	// Should make a custom unmarshaller at some point
	AudioChannels         float64
	AudioCodec            string
	AudioLanguages        []string
	Height                int32
	Width                 int32
	Subtitles             []string
	VideoCodec            string
	VideoDynamicRange     string
	VideoDynamicRangeType string
}

type WebhookCustomFormatInfo struct {
	CustomFormats     []WebhookCustomFormat
	CustomFormatScore int32
}

type WebhookCustomFormat struct {
	Id   int32
	Name string
}

type RadarrWebhookPayload struct {
	Movie              WebhookMovie
	RemoteMovie        WebhookRemoteMovie
	MovieFile          WebhookMovieFile
	IsUpgrade          bool
	DownloadClient     string
	DownloadClientType string
	DownloadId         string
	DeletedFiles       []WebhookMovieFile
	CustomFormatInfo   WebhookCustomFormatInfo
}

func ServeRadarr(w http.ResponseWriter, r *http.Request) {
	log.Debugln("Received Radarr webhook")

	var data RadarrWebhookPayload
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.WithError(err).Errorln("Failed to decode Radarr webhook JSON data")
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	path := filepath.Join(data.Movie.FolderPath, data.MovieFile.RelativePath)

	log.WithField("data", fmt.Sprintf("%+v", data)).Debugln("Decoded Radarr webhook json data")
	pkg.EnqueueSub(path)
	log.Debugf("Queued %v from radarr", data.MovieFile.Path)
	w.WriteHeader(200)
}
