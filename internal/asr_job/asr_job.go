package asr_job

import (
	"time"
)

type FileAsrJob struct {
	ID         uint
	FilePath   string
	Lang       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Status     JobStatus
	Progress   float32
	DurationMS float32
}
