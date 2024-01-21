package job_queue

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/emirpasic/gods/lists/arraylist"
	"go-subgen/internal/asr_job"
)

type ArrayRepository struct {
	jobs       arraylist.List
	lock       *sync.RWMutex
	newJobChan chan uint
	idCounter  atomic.Uint32
}

func (repo *ArrayRepository) GetNewJobChannel() chan uint {
	return repo.newJobChan
}

func NewArrayRepository() *ArrayRepository {
	return &ArrayRepository{
		jobs:       *arraylist.New(),
		lock:       &sync.RWMutex{},
		newJobChan: make(chan uint, 100),
	}
}

func (repo *ArrayRepository) getNextId() uint {
	return uint(repo.idCounter.Add(1))
}

func (repo *ArrayRepository) AddJob(ctx context.Context, job asr_job.FileAsrJob) (asr_job.FileAsrJob, error) {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	job.ID = repo.getNextId()
	job.UpdatedAt = time.Now()

	repo.jobs.Add(job)
	repo.newJobChan <- job.ID

	return job, nil
}

func (repo *ArrayRepository) AddJobs(ctx context.Context, jobs []asr_job.FileAsrJob) ([]asr_job.FileAsrJob, error) {
	repo.lock.Lock()
	defer repo.lock.Unlock()
	var addedJobs []asr_job.FileAsrJob
	for _, job := range jobs {

		job.ID = repo.getNextId()
		job.UpdatedAt = time.Now()

		repo.jobs.Add(job)
		repo.newJobChan <- job.ID
		addedJobs = append(addedJobs, job)
	}
	return addedJobs, nil
}

func (repo *ArrayRepository) RemoveJob(ctx context.Context, jobId uint) (asr_job.FileAsrJob, error) {
	repo.lock.Lock()
	defer repo.lock.Unlock()
	for i := 0; i < repo.jobs.Size(); i++ {
		if jobInterface, ok := repo.jobs.Get(i); ok {
			if job, ok := jobInterface.(asr_job.FileAsrJob); ok && job.ID == jobId {
				repo.jobs.Remove(i)
				return job, nil
			}
		}
	}

	return asr_job.FileAsrJob{}, fmt.Errorf("job not found")
}

func (repo *ArrayRepository) GetJob(ctx context.Context, jobId uint) (asr_job.FileAsrJob, error) {
	repo.lock.RLock()
	defer repo.lock.RUnlock()
	for i := 0; i < repo.jobs.Size(); i++ {
		if jobInterface, ok := repo.jobs.Get(i); ok {
			if job, ok := jobInterface.(asr_job.FileAsrJob); ok && job.ID == jobId {
				return job, nil
			}
		}
	}

	return asr_job.FileAsrJob{}, fmt.Errorf("job not found")
}

// GetNextJob returns the next job in the repository.
// It acquires a read lock before accessing the repository, and releases the lock before returning.
// If the job list is not empty, it returns the first job in the repository.
// Otherwise, it returns an empty FileAsrJob and an error indicating the job list is empty.
func (repo *ArrayRepository) GetNextJob(ctx context.Context) (asr_job.FileAsrJob, error) {
	repo.lock.RLock()
	defer repo.lock.RUnlock()
	// todo we don't want the first job, we want the oldest job that's Queued (0)
	if repo.jobs.Size() > 0 {
		if jobInterface, ok := repo.jobs.Get(0); ok {
			if job, ok := jobInterface.(asr_job.FileAsrJob); ok {
				return job, nil
			}
		}
	}

	return asr_job.FileAsrJob{}, fmt.Errorf("job list is empty")
}

func (repo *ArrayRepository) JobCount(ctx context.Context) (int, error) {
	repo.lock.RLock()
	defer repo.lock.RUnlock()
	return repo.jobs.Size(), nil
}

func (repo *ArrayRepository) GetAllJobs(ctx context.Context) ([]asr_job.FileAsrJob, error) {
	repo.lock.RLock()
	defer repo.lock.RUnlock()
	allJobs := make([]asr_job.FileAsrJob, repo.jobs.Size())
	for i := 0; i < repo.jobs.Size(); i++ {
		if jobInterface, ok := repo.jobs.Get(i); ok {
			if job, ok := jobInterface.(asr_job.FileAsrJob); ok {
				allJobs[i] = job
			}
		}
	}

	return allJobs, nil
}

// UpdateJob updates a job in the repository based on JobId, replacing it with the provided FileAsrJob.
func (repo *ArrayRepository) UpdateJob(ctx context.Context, Id uint, newJob asr_job.FileAsrJob) (asr_job.FileAsrJob, error) {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	for i := 0; i < repo.jobs.Size(); i++ {
		if jobInterface, ok := repo.jobs.Get(i); ok {
			if fileJob, ok := jobInterface.(asr_job.FileAsrJob); ok && fileJob.ID == Id {
				newJob.UpdatedAt = time.Now()
				repo.jobs.Set(i, newJob)
				return fileJob, nil
			}
		}
	}

	return asr_job.FileAsrJob{}, fmt.Errorf("job not found")
}

func (repo *ArrayRepository) SetJobStatus(ctx context.Context, JobId uint, status asr_job.JobStatus) error {
	repo.lock.Lock()
	defer repo.lock.Unlock()
	for i := 0; i < repo.jobs.Size(); i++ {
		if jobInterface, ok := repo.jobs.Get(i); ok {
			if job, ok := jobInterface.(asr_job.FileAsrJob); ok && job.ID == JobId {
				job.Status = status
				job.UpdatedAt = time.Now()
				repo.jobs.Set(i, job)
				return nil
			}
		}
	}
	return fmt.Errorf("job not found")
}
func (repo *ArrayRepository) SetJobProgress(ctx context.Context, jobId uint, progress float32) error {
	repo.lock.Lock()
	defer repo.lock.Unlock()
	_, jobInterface := repo.jobs.Find(func(index int, value interface{}) bool {
		return value.(asr_job.FileAsrJob).ID == jobId
	})

	p, _ := jobInterface.(asr_job.FileAsrJob)
	p.Progress = progress

	for i := 0; i < repo.jobs.Size(); i++ {
		if jobInterface, ok := repo.jobs.Get(i); ok {
			if job, ok := jobInterface.(asr_job.FileAsrJob); ok && job.ID == jobId {
				job.Progress = progress
				job.UpdatedAt = time.Now()
				repo.jobs.Set(i, job)
				return nil
			}
		}
	}
	return fmt.Errorf("job not found")
}
