package api

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
	"go-subgen/internal/asr_job"
)

type JobHandler interface {
	ServeJobsRequest(w http.ResponseWriter, r *http.Request)
}
type jobHandler struct {
	QueueRepository asr_job.AsrJobRepository
}

func NewJobHandler(repository asr_job.AsrJobRepository) GenericFileHandler {
	return &jobHandler{
		QueueRepository: repository,
	}
}

func (h jobHandler) Serve(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.QueueRepository.GetAllJobs(nil)
	if err != nil {
		log.WithError(err).Error("Failed to serve jobs")
		http.Error(w, err.Error(), 500)
		return
	}
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(jobs)
	if err != nil {
		log.WithError(err).Error("Failed to serve jobs")
		http.Error(w, err.Error(), 500)
		return
	}
	return
}
