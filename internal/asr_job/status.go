package asr_job

type JobStatus uint8

const (
	Queued JobStatus = iota
	InProgress
	Canceled
	Complete
	Failed
)
