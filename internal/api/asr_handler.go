package api

import "go-subgen/internal/asr_job"

type AsrHandler struct {
	QueueRepository asr_job.AsrJobRepository
}

func NewAsrHandler(queueRepository asr_job.AsrJobRepository) AsrHandler {
	return AsrHandler{
		QueueRepository: queueRepository,
	}
}
