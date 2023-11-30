package job_queue

import (
	"context"
	"fmt"
	"sync"

	"github.com/emirpasic/gods/lists/arraylist"
	"go-subgen/internal/asr_job"
)

type ArrayRepository struct {
	jobs arraylist.List
	lock *sync.RWMutex
}

func NewArrayRepository() *ArrayRepository {
	return &ArrayRepository{
		jobs: *arraylist.New(),
		lock: &sync.RWMutex{},
	}
}

func (repo *ArrayRepository) AddJob(ctx context.Context, job asr_job.FileAsrJob) error {
	repo.lock.Lock()
	defer repo.lock.Unlock()
	repo.jobs.Add(job)

	return nil
}

func (repo *ArrayRepository) AddJobs(ctx context.Context, jobs []asr_job.FileAsrJob) {
	repo.lock.Lock()
	defer repo.lock.Unlock()
	for _, job := range jobs {
		repo.jobs.Add(job)
	}
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
func (repo *ArrayRepository) GetNextJob() (asr_job.FileAsrJob, error) {
	repo.lock.RLock()
	defer repo.lock.RUnlock()
	// todo we don't want the first job, we want the oldest job thats Queued (0)
	if repo.jobs.Size() > 0 {
		if jobInterface, ok := repo.jobs.Get(0); ok {
			if job, ok := jobInterface.(asr_job.FileAsrJob); ok {
				return job, nil
			}
		}
	}

	return asr_job.FileAsrJob{}, fmt.Errorf("job list is empty")
}

func (repo *ArrayRepository) JobCount() (int, error) {
	return repo.jobs.Size(), nil
}

func (repo *ArrayRepository) GetAllJobs() ([]asr_job.FileAsrJob, error) {
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
