package webhooks

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type WebhookSeries struct {
	Id        int32
	Title     string
	TitleSlug string
	Path      string
	TvdbId    int32
	TvMazeId  int32
	ImdbId    string
	Type      SeriesTypes
}

//go:generate go-enum -type=SeriesTypes -all=false -string=false -new=true -text=true -json=true -yaml=false

type SeriesTypes uint8

const (
	standard SeriesTypes = iota
	daily
	anime
)

type WebhookEpisode struct {
	Id            int32
	EpisodeNumber int32
	SeasonNumber  int32
	Title         string
	OverView      string
	AirDate       string
	AirDateUtc    string
	SeriesId      int32
}

type WebhookEpisodeFile struct {
	Id             int32
	RelativePath   string
	Path           string
	Quality        string
	QualityVersion int32
	ReleaseGroup   string
	SceneName      string
	Size           int64
	DateAdded      string
	MediaInfo      WebhookEpisodeFileMediaInfo
}

type WebhookEpisodeFileMediaInfo struct {
	AudioChannels         big.Float
	AudioCodec            string
	AudioLanguages        []string
	Height                int32
	Width                 int32
	Subtitles             []string
	VideoCodec            string
	VideoDynamicRange     string
	VideoDynamicRangeType string
}

type SonarrWebhookPayload struct {
	Series             WebhookSeries
	Episodes           []WebhookEpisode
	EpisodeFile        WebhookEpisodeFile
	IsUpgrade          bool
	DownloadClient     string
	DownloadClientType string
	DownloadId         string
	DeletedFiles       []WebhookMovieFile
	CustomFormatInfo   WebhookCustomFormatInfo
}

func ServeSonarr(w http.ResponseWriter, r *http.Request) {
	log.Debugln("Received Sonarr webhook")

	var data SonarrWebhookPayload
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.WithError(err).Errorln("Failed to decode Sonarr webhook JSON data")
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	path := filepath.Join(data.Series.Path, data.EpisodeFile.RelativePath)

	log.WithField("data", fmt.Sprintf("%+v", data)).Debugln("Decoded Sonarr webhook json data")
	// internal.EnqueueSub(path)
	log.Debugf("Queued %v from sonarr", path)
	w.WriteHeader(200)
}
