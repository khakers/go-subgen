package interceptors

import (
	"context"

	"go-subgen/internal/asr_job"
)

// Intercepts calls to setProgress and passes them to the progress channel
// all other calls in the asr_job repository are passed through to the underlying repository
type AsrProgressInterceptor struct {
	Repository      asr_job.AsrJobRepository
	ProgressChannel chan AsrProgressEvent
}

type AsrProgressEvent struct {
	JobId    uint
	Progress float32
}

func NewAsrProgressInterceptor(repository asr_job.AsrJobRepository, progressChannel chan AsrProgressEvent) *AsrProgressInterceptor {
	return &AsrProgressInterceptor{
		Repository:      repository,
		ProgressChannel: progressChannel,
	}
}

func (i AsrProgressInterceptor) SetJobProgress(ctx context.Context, jobId uint, progress float32) error {
	i.ProgressChannel <- AsrProgressEvent{
		JobId:    jobId,
		Progress: progress,
	}
	return i.Repository.SetJobProgress(ctx, jobId, progress)
}

func (i AsrProgressInterceptor) AddJob(ctx context.Context, job asr_job.FileAsrJob) (asr_job.FileAsrJob, error) {
	return i.Repository.AddJob(ctx, job)
}

func (i AsrProgressInterceptor) AddJobs(ctx context.Context, jobs []asr_job.FileAsrJob) ([]asr_job.FileAsrJob, error) {
	i.Repository.AddJobs(ctx, jobs)
	return nil, nil
}

func (i AsrProgressInterceptor) RemoveJob(ctx context.Context, JobId uint) (asr_job.FileAsrJob, error) {
	return i.Repository.RemoveJob(ctx, JobId)
}

func (i AsrProgressInterceptor) GetJob(ctx context.Context, JobId uint) (asr_job.FileAsrJob, error) {
	return i.Repository.GetJob(ctx, JobId)
}

func (i AsrProgressInterceptor) GetNextJob(ctx context.Context) (asr_job.FileAsrJob, error) {
	return i.Repository.GetNextJob(ctx)
}

func (i AsrProgressInterceptor) JobCount(ctx context.Context) (int, error) {
	return i.Repository.JobCount(ctx)
}

func (i AsrProgressInterceptor) GetAllJobs(ctx context.Context) ([]asr_job.FileAsrJob, error) {
	return i.Repository.GetAllJobs(ctx)
}

func (i AsrProgressInterceptor) UpdateJob(ctx context.Context, JobId uint, job asr_job.FileAsrJob) (asr_job.FileAsrJob, error) {
	return i.Repository.UpdateJob(ctx, JobId, job)
}

func (i AsrProgressInterceptor) SetJobStatus(ctx context.Context, JobId uint, status asr_job.JobStatus) error {
	return i.Repository.SetJobStatus(ctx, JobId, status)
}

func (i AsrProgressInterceptor) GetNewJobChannel() chan uint {
	return i.Repository.GetNewJobChannel()
}
