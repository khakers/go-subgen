package asr_job

import "time"

type FileAsrJob struct {
	FilePath     string
	Lang         string
	CreationTime time.Time
	Status       JobStatus
	Progress     float32
}
