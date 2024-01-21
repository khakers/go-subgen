package interceptors

import (
	"context"
	"time"

	"go-subgen/internal/asr_job"
)

// Intercepts calls to setProgress and passes them to the progress channel
// all other calls in the asr_job repository are passed through to the underlying repository
type AsrStatusInterceptor struct {
	Repository   asr_job.AsrJobRepository
	EventChannel chan AsrStatusEvent
}

type AsrStatusEvent struct {
	JobId  uint
	Status asr_job.JobStatus
	Time   time.Time
}

func NewAsrStatusInterceptor(repository asr_job.AsrJobRepository, progressChannel chan AsrStatusEvent) *AsrStatusInterceptor {
	return &AsrStatusInterceptor{
		Repository:   repository,
		EventChannel: progressChannel,
	}
}

func (i AsrStatusInterceptor) SetJobProgress(ctx context.Context, jobId uint, progress float32) error {
	return i.Repository.SetJobProgress(ctx, jobId, progress)
}

func (i AsrStatusInterceptor) AddJob(ctx context.Context, job asr_job.FileAsrJob) (asr_job.FileAsrJob, error) {
	return i.Repository.AddJob(ctx, job)
}

func (i AsrStatusInterceptor) AddJobs(ctx context.Context, jobs []asr_job.FileAsrJob) ([]asr_job.FileAsrJob, error) {
	i.Repository.AddJobs(ctx, jobs)
	return nil, nil
}

func (i AsrStatusInterceptor) RemoveJob(ctx context.Context, JobId uint) (asr_job.FileAsrJob, error) {
	return i.Repository.RemoveJob(ctx, JobId)
}

func (i AsrStatusInterceptor) GetJob(ctx context.Context, JobId uint) (asr_job.FileAsrJob, error) {
	return i.Repository.GetJob(ctx, JobId)
}

func (i AsrStatusInterceptor) GetNextJob(ctx context.Context) (asr_job.FileAsrJob, error) {
	return i.Repository.GetNextJob(ctx)
}

func (i AsrStatusInterceptor) JobCount(ctx context.Context) (int, error) {
	return i.Repository.JobCount(ctx)
}

func (i AsrStatusInterceptor) GetAllJobs(ctx context.Context) ([]asr_job.FileAsrJob, error) {
	return i.Repository.GetAllJobs(ctx)
}

func (i AsrStatusInterceptor) UpdateJob(ctx context.Context, JobId uint, job asr_job.FileAsrJob) (asr_job.FileAsrJob, error) {
	return i.Repository.UpdateJob(ctx, JobId, job)
}

func (i AsrStatusInterceptor) SetJobStatus(ctx context.Context, JobId uint, status asr_job.JobStatus) error {
	i.EventChannel <- AsrStatusEvent{
		JobId:  JobId,
		Status: status,
		Time:   time.Now(),
	}
	return i.Repository.SetJobStatus(ctx, JobId, status)
}

func (i AsrStatusInterceptor) GetNewJobChannel() chan uint {
	return i.Repository.GetNewJobChannel()
}
