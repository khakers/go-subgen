package asr_job

import "time"

type FileAsrJob struct {
	FilePath string
	Lang     string
	addTime  time.Time
}
