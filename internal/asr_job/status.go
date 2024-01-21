package asr_job

type JobStatus uint8

const (
	Queued JobStatus = iota
	AudioStripping
	InProgress
	Canceled
	Complete
	Failed
)
