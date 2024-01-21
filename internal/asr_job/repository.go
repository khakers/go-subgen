package asr_job

import (
	"context"
)

// AsrJobRepository is the interface for the job store
type AsrJobRepository interface {
	AddJob(ctx context.Context, job FileAsrJob) (FileAsrJob, error)
	AddJobs(ctx context.Context, jobs []FileAsrJob) ([]FileAsrJob, error)
	RemoveJob(ctx context.Context, JobId uint) (FileAsrJob, error)
	GetJob(ctx context.Context, JobId uint) (FileAsrJob, error)
	GetNextJob(ctx context.Context) (FileAsrJob, error)
	JobCount(ctx context.Context) (int, error)
	GetAllJobs(ctx context.Context) ([]FileAsrJob, error)
	UpdateJob(ctx context.Context, JobId uint, job FileAsrJob) (FileAsrJob, error)
	SetJobStatus(ctx context.Context, JobId uint, status JobStatus) error
	SetJobProgress(ctx context.Context, jobId uint, progress float32) error
	GetNewJobChannel() chan uint
}
