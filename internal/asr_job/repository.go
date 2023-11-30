package asr_job

import (
	"context"
)

// AsrJobRepository is the interface for the job store
type AsrJobRepository interface {
	AddJob(ctx context.Context, job FileAsrJob) error
	AddJobs(ctx context.Context, jobs []FileAsrJob)
	RemoveJob(ctx context.Context, JobId uint) (FileAsrJob, error)
	GetJob(ctx context.Context, JobId uint) (FileAsrJob, error)
	GetNextJob() (FileAsrJob, error)
	JobCount() (int, error)
	GetAllJobs() ([]FileAsrJob, error)
}
