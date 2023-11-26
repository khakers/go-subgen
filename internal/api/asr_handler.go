package api

import "go-subgen/internal/asr_job"

type AsrHandler struct {
	QueueRepository asr_job.AsrJobQueueRepository
}

func NewAsrHandler(queueRepository asr_job.AsrJobQueueRepository) AsrHandler {
	return AsrHandler{
		QueueRepository: queueRepository,
	}
}
