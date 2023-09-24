package job_queue

import (
	"container/list"
	"context"
	"errors"
	"go-subgen/internal/asr_job"
	"sync"
)

// memoryQueue is a simple in-memory FIFO queue that implements the AsrJobQueueRepository interface

type MemoryQueueRepository struct {
	jobs list.List
	lock *sync.RWMutex
}

// NewMemoryQueueRepository creates a new instance of the memoryQueue struct
func NewMemoryQueueRepository() *MemoryQueueRepository {
	return &MemoryQueueRepository{
		jobs: *list.New(),
		lock: &sync.RWMutex{},
	}
}

// EnqueueJob adds a job to the queue
func (q *MemoryQueueRepository) EnqueueJob(ctx context.Context, job asr_job.FileAsrJob) error {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.enqueueJob(job)
}

func (q *MemoryQueueRepository) enqueueJob(job asr_job.FileAsrJob) error {
	q.jobs.PushFront(job)
	return nil
}

// DequeueJob removes a job from the queue
func (q *MemoryQueueRepository) DequeueJob(ctx context.Context) (asr_job.FileAsrJob, error) {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.dequeueJob()
}

// dequeueJob removes a job from the queue without locking
func (q *MemoryQueueRepository) dequeueJob() (asr_job.FileAsrJob, error) {
	if q.jobs.Len() == 0 {
		return asr_job.FileAsrJob{}, errors.New("no jobs in queue")
	}
	return q.jobs.Remove(q.jobs.Back()).(asr_job.FileAsrJob), nil
}

func (q *MemoryQueueRepository) PeekJob() (asr_job.FileAsrJob, error) {
	q.lock.RLock()
	defer q.lock.RUnlock()

	if q.jobs.Len() == 0 {
		return asr_job.FileAsrJob{}, errors.New("no jobs in queue")
	}

	return q.jobs.Back().Value.(asr_job.FileAsrJob), nil

}

func (q *MemoryQueueRepository) JobCount() (int, error) {
	q.lock.RLock()
	defer q.lock.RUnlock()
	return q.jobs.Len(), nil
}

func (q *MemoryQueueRepository) PeekJobs() ([]asr_job.FileAsrJob, error) {
	q.lock.RLock()
	defer q.lock.RUnlock()

	if q.jobs.Len() == 0 {
		return make([]asr_job.FileAsrJob, 0), errors.New("no jobs in queue")
	}

	jobs := make([]asr_job.FileAsrJob, q.jobs.Len())
	for i, e := 0, q.jobs.Front(); e != nil; i, e = i+1, e.Next() {
		jobs[i] = e.Value.(asr_job.FileAsrJob)
	}

	return jobs, nil
}
