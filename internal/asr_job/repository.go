package asr_job

import "context"

// AsrJobQueueRepository is the interface for the job queue
type AsrJobQueueRepository interface {
	EnqueueJob(ctx context.Context, job FileAsrJob) error
	DequeueJob(ctx context.Context) (FileAsrJob, error)
	PeekJob() (FileAsrJob, error)
	JobCount() (int, error)
	PeekJobs() ([]FileAsrJob, error)
}
