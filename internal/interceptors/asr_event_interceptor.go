package interceptors

import (
	"context"
	"time"

	"go-subgen/internal/asr_job"
)

// Intercepts calls to setProgress and passes them to the progress channel
// all other calls in the asr_job repository are passed through to the underlying repository
type AsrEventInterceptor struct {
	Repository    asr_job.AsrJobRepository
	EventChannels EventChannels
}

type AsrJobEvent struct {
	JobId      uint
	ChangeType ChangeType
	Job        asr_job.FileAsrJob
}

type ChangeType uint8

const (
	New ChangeType = iota
	Update
	Delete
)

type EventChannels struct {
	ProgressChannel       chan AsrProgressEvent
	StatusChannel         chan AsrStatusEvent
	JobChangeEventChannel chan AsrJobEvent
}

func NewAsrEventInterceptor(repository asr_job.AsrJobRepository) *AsrEventInterceptor {
	return &AsrEventInterceptor{
		Repository: repository,
		EventChannels: EventChannels{
			ProgressChannel:       make(chan AsrProgressEvent, 20),
			StatusChannel:         make(chan AsrStatusEvent, 20),
			JobChangeEventChannel: make(chan AsrJobEvent, 20),
		},
	}
}

func (i AsrEventInterceptor) SetJobProgress(ctx context.Context, jobId uint, progress float32) error {
	i.EventChannels.ProgressChannel <- AsrProgressEvent{
		JobId:    jobId,
		Progress: progress,
	}
	return i.Repository.SetJobProgress(ctx, jobId, progress)
}

func (i AsrEventInterceptor) AddJob(ctx context.Context, inboundAsrJob asr_job.FileAsrJob) (asr_job.FileAsrJob, error) {

	outboundAsrJob, err := i.Repository.AddJob(ctx, inboundAsrJob)

	i.EventChannels.JobChangeEventChannel <- AsrJobEvent{
		JobId:      outboundAsrJob.ID,
		ChangeType: New,
		Job:        outboundAsrJob,
	}

	return outboundAsrJob, err
}

func (i AsrEventInterceptor) AddJobs(ctx context.Context, jobs []asr_job.FileAsrJob) ([]asr_job.FileAsrJob, error) {

	addedJobs, err := i.Repository.AddJobs(ctx, jobs)

	for _, job := range addedJobs {
		i.EventChannels.JobChangeEventChannel <- AsrJobEvent{
			JobId:      job.ID,
			ChangeType: New,
			Job:        job,
		}
	}

	return addedJobs, err
}

func (i AsrEventInterceptor) RemoveJob(ctx context.Context, JobId uint) (asr_job.FileAsrJob, error) {
	i.EventChannels.JobChangeEventChannel <- AsrJobEvent{
		JobId:      JobId,
		ChangeType: Delete,
	}
	return i.Repository.RemoveJob(ctx, JobId)
}

func (i AsrEventInterceptor) GetJob(ctx context.Context, JobId uint) (asr_job.FileAsrJob, error) {
	return i.Repository.GetJob(ctx, JobId)
}

func (i AsrEventInterceptor) GetNextJob(ctx context.Context) (asr_job.FileAsrJob, error) {
	return i.Repository.GetNextJob(ctx)
}

func (i AsrEventInterceptor) JobCount(ctx context.Context) (int, error) {
	return i.Repository.JobCount(ctx)
}

func (i AsrEventInterceptor) GetAllJobs(ctx context.Context) ([]asr_job.FileAsrJob, error) {
	return i.Repository.GetAllJobs(ctx)
}

func (i AsrEventInterceptor) UpdateJob(ctx context.Context, JobId uint, job asr_job.FileAsrJob) (asr_job.FileAsrJob, error) {
	i.EventChannels.JobChangeEventChannel <- AsrJobEvent{
		JobId:      JobId,
		ChangeType: Update,
		Job:        job,
	}
	return i.Repository.UpdateJob(ctx, JobId, job)
}

func (i AsrEventInterceptor) SetJobStatus(ctx context.Context, JobId uint, status asr_job.JobStatus) error {
	i.EventChannels.StatusChannel <- AsrStatusEvent{
		JobId:  JobId,
		Status: status,
		Time:   time.Now(),
	}
	return i.Repository.SetJobStatus(ctx, JobId, status)
}

func (i AsrEventInterceptor) GetNewJobChannel() chan uint {
	return i.Repository.GetNewJobChannel()
}
