package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go-subgen/internal/asr_job"

	log "github.com/sirupsen/logrus"
)

type GenericWebhookData struct {
	Files []string `json:"files"`
}

type GenericFileHandler interface {
	Serve(w http.ResponseWriter, r *http.Request)
}
type genericFileHandler struct {
	QueueRepository asr_job.AsrJobQueueRepository
}

func NewGenericFileHandler(repository asr_job.AsrJobQueueRepository) GenericFileHandler {
	return &genericFileHandler{
		QueueRepository: repository,
	}
}

func (h genericFileHandler) Serve(w http.ResponseWriter, r *http.Request) {
	log.Debugln("Received generic webhook")

	var data GenericWebhookData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.WithError(err).Errorln("Failed to decode webhook JSON data")

		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	log.WithField("data", fmt.Sprintf("%+v", data)).Debugln("Decoded webhook json data")

	for _, file := range data.Files {
		log.Debugf("Queued %v (API)", file)
		err := h.QueueRepository.EnqueueJob(
			r.Context(),
			asr_job.FileAsrJob{
				FilePath: file,
				Lang:     "en",
			},
		)
		if err != nil {
			return
		}
	}
	w.WriteHeader(200)
}
