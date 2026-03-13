package queue

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"go-auth/internal/model"
	"go-auth/internal/repository"
	"go-auth/internal/service"
)

// JobQueue manages the video generation job queue
type JobQueue interface {
	// Enqueue adds a job to the queue
	Enqueue(ctx context.Context, job *model.GenerationJob) error

	// Dequeue retrieves the next job from the queue
	Dequeue(ctx context.Context, limit int) ([]model.GenerationJob, error)

	// MarkProcessing marks a job as processing
	MarkProcessing(ctx context.Context, jobID string) error

	// MarkComplete marks a job as complete
	MarkComplete(ctx context.Context, jobID string) error

	// MarkFailed marks a job as failed
	MarkFailed(ctx context.Context, jobID string, errorMsg string) error

	// GetStats get queue statistics
	GetStats(ctx context.Context) map[string]interface{}

	// Start the worker process
	Start(ctx context.Context, workerCount int) error

	// Stop the worker process
	Stop() error
}

// SimpleJobQueue is a memory-backed job queue implementation
type SimpleJobQueue struct {
	jobRepo          repository.GenerationJobRepository
	videoGenService  service.VideoGenerationService
	maxRetries       int
	pollInterval     time.Duration
	processingWorkers int
	stopChan         chan bool
	mu               sync.RWMutex
	isRunning        bool
}

// NewJobQueue creates a new job queue
func NewJobQueue(
	jobRepo repository.GenerationJobRepository,
	videoGenService service.VideoGenerationService,
) JobQueue {
	return &SimpleJobQueue{
		jobRepo:         jobRepo,
		videoGenService: videoGenService,
		maxRetries:      3,
		pollInterval:    10 * time.Second,
		stopChan:        make(chan bool),
		isRunning:       false,
	}
}

func (q *SimpleJobQueue) Enqueue(ctx context.Context, job *model.GenerationJob) error {
	// Job is already created in the repository, just ensure it's queued
	job.Status = "queued"
	return q.jobRepo.Update(ctx, job)
}

func (q *SimpleJobQueue) Dequeue(ctx context.Context, limit int) ([]model.GenerationJob, error) {
	// Get pending jobs ordered by priority and creation time
	return q.jobRepo.GetPendingJobs(ctx, limit)
}

func (q *SimpleJobQueue) MarkProcessing(ctx context.Context, jobID string) error {
	// Parse jobID as UUID
	id, err := parseUUID(jobID)
	if err != nil {
		return err
	}
	return q.jobRepo.UpdateStatus(ctx, id, "processing", "")
}

func (q *SimpleJobQueue) MarkComplete(ctx context.Context, jobID string) error {
	id, err := parseUUID(jobID)
	if err != nil {
		return err
	}
	return q.jobRepo.UpdateStatus(ctx, id, "completed", "")
}

func (q *SimpleJobQueue) MarkFailed(ctx context.Context, jobID string, errorMsg string) error {
	id, err := parseUUID(jobID)
	if err != nil {
		return err
	}
	return q.jobRepo.UpdateStatus(ctx, id, "failed", errorMsg)
}

func (q *SimpleJobQueue) GetStats(ctx context.Context) map[string]interface{} {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return map[string]interface{}{
		"is_running":         q.isRunning,
		"processing_workers": q.processingWorkers,
		"poll_interval":      q.pollInterval.String(),
	}
}

func (q *SimpleJobQueue) Start(ctx context.Context, workerCount int) error {
	q.mu.Lock()
	if q.isRunning {
		q.mu.Unlock()
		return fmt.Errorf("queue is already running")
	}
	q.isRunning = true
	q.processingWorkers = workerCount
	q.mu.Unlock()

	// Start worker goroutines
	for i := 0; i < workerCount; i++ {
		go q.worker(ctx, i)
	}

	log.Printf("Started job queue with %d workers", workerCount)
	return nil
}

func (q *SimpleJobQueue) Stop() error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if !q.isRunning {
		return fmt.Errorf("queue is not running")
	}

	q.isRunning = false
	close(q.stopChan)
	log.Println("Stopped job queue")
	return nil
}

func (q *SimpleJobQueue) worker(ctx context.Context, workerID int) {
	log.Printf("Worker %d started", workerID)
	ticker := time.NewTicker(q.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-q.stopChan:
			log.Printf("Worker %d stopped", workerID)
			return

		case <-ticker.C:
			// Get pending jobs
			jobs, err := q.Dequeue(ctx, 1)
			if err != nil || len(jobs) == 0 {
				continue
			}

			job := jobs[0]
			log.Printf("Worker %d processing job %s (type: %s)", workerID, job.ID, job.JobType)

			// Mark as processing
			q.MarkProcessing(ctx, job.ID.String())

			// Process the job
			if err := q.videoGenService.ProcessGenerationJob(ctx, job.ID); err != nil {
				log.Printf("Worker %d error processing job %s: %v", workerID, job.ID, err)
				
				// Retry logic
				if job.RetryCount < job.MaxRetries {
					job.RetryCount++
					job.Status = "queued" // Re-queue for retry
					q.jobRepo.Update(ctx, &job)
				} else {
					q.MarkFailed(ctx, job.ID.String(), fmt.Sprintf("Max retries exceeded: %v", err))
				}
				continue
			}

			// Start polling job status
			go q.pollJobStatus(ctx, job.ID)
		}
	}
}

func (q *SimpleJobQueue) pollJobStatus(ctx context.Context, jobID uuid.UUID) {
	maxAttempts := 120 // Poll for up to 2 hours with 60-second intervals
	attempt := 0
	pollTicker := time.NewTicker(60 * time.Second)
	defer pollTicker.Stop()

	for {
		select {
		case <-q.stopChan:
			return

		case <-pollTicker.C:
			attempt++

			// Check job status
			job, err := q.jobRepo.GetByID(ctx, jobID)
			if err != nil {
				log.Printf("Error getting job status for %s: %v", jobID, err)
				return
			}

			if job.Status == "completed" || job.Status == "failed" {
				log.Printf("Job %s finished with status: %s", jobID, job.Status)
				return
			}

			// Poll video provider for updates
			if err := q.videoGenService.PollJobStatus(ctx, jobID); err != nil {
				log.Printf("Error polling job status for %s: %v", jobID, err)
			}

			if attempt >= maxAttempts {
				log.Printf("Job %s polling timeout after %d attempts", jobID, attempt)
				q.MarkFailed(ctx, jobID.String(), "Polling timeout")
				return
			}
		}
	}
}

// Helper function to parse UUID string
func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
