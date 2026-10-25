package service

import (
	"context"
	"log"
	"sync"

	"FinalProject/model"
	"FinalProject/repository"
)

type embeddingJob struct {
	photoID uint
	dest    string
}

type EmbeddingStatusUpdate struct {
	PhotoID     uint     `json:"photoId"`
	Status      string   `json:"status"`
	Error       string   `json:"error"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

// EmbeddingQueue creates image embeddings and similarity-based tags outside the photo-move call.
type EmbeddingQueue struct {
	photoRepo    repository.PhotoRepository
	settingRepo  repository.SettingRepository
	photoService PhotoService

	mu       sync.Mutex
	jobs     []embeddingJob
	queued   map[uint]bool
	wake     chan struct{}
	started  bool
	paused   bool
	onUpdate func(update EmbeddingStatusUpdate)
}

func NewEmbeddingQueue(
	photoRepo repository.PhotoRepository,
	settingRepo repository.SettingRepository,
	photoService PhotoService,
) *EmbeddingQueue {
	return &EmbeddingQueue{
		photoRepo:    photoRepo,
		settingRepo:  settingRepo,
		photoService: photoService,
		queued:       make(map[uint]bool),
		wake:         make(chan struct{}, 1),
	}
}

// Start recovers unfinished jobs from the database, then starts one worker.
func (q *EmbeddingQueue) Start(
	ctx context.Context,
	onUpdate func(update EmbeddingStatusUpdate),
) error {
	setting, err := q.settingRepo.GetSetting(ctx)
	if err != nil {
		return err
	}

	if setting.Dest != "" {
		photos, err := q.photoRepo.GetPhotosNeedingAnalysis(ctx, setting.Dest)
		if err != nil {
			return err
		}
		for _, photo := range photos {
			q.Enqueue(photo.ID, photo.Dest)
		}
	}

	q.mu.Lock()
	if q.started {
		q.mu.Unlock()
		return nil
	}
	q.started = true
	q.onUpdate = onUpdate
	q.mu.Unlock()

	go q.run(ctx)
	return nil
}

// Enqueue only records the job in memory and signals the worker; it does not
// wait for the model to process the image.
func (q *EmbeddingQueue) Enqueue(photoID uint, dest string) {
	q.mu.Lock()
	if q.queued[photoID] {
		q.mu.Unlock()
		return
	}
	q.queued[photoID] = true
	q.jobs = append(q.jobs, embeddingJob{photoID: photoID, dest: dest})
	q.mu.Unlock()

	select {
	case q.wake <- struct{}{}:
	default:
	}
}

// Pause stops the worker from starting another queued job. A job already in
// progress is allowed to finish so its database status remains consistent.
func (q *EmbeddingQueue) Pause() {
	q.mu.Lock()
	q.paused = true
	q.mu.Unlock()
}

func (q *EmbeddingQueue) Resume() {
	q.mu.Lock()
	q.paused = false
	q.mu.Unlock()
	select {
	case q.wake <- struct{}{}:
	default:
	}
}

func (q *EmbeddingQueue) IsPaused() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.paused
}

func (q *EmbeddingQueue) run(ctx context.Context) {
	for {
		job, ok := q.next(ctx)
		if !ok {
			return
		}
		q.process(ctx, job)
		q.mu.Lock()
		delete(q.queued, job.photoID)
		q.mu.Unlock()
	}
}

func (q *EmbeddingQueue) next(ctx context.Context) (embeddingJob, bool) {
	for {
		q.mu.Lock()
		if len(q.jobs) > 0 && !q.paused {
			job := q.jobs[0]
			q.jobs[0] = embeddingJob{}
			q.jobs = q.jobs[1:]
			q.mu.Unlock()
			return job, true
		}
		q.mu.Unlock()

		select {
		case <-ctx.Done():
			return embeddingJob{}, false
		case <-q.wake:
		}
	}
}

func (q *EmbeddingQueue) process(ctx context.Context, job embeddingJob) {
	if err := q.photoRepo.UpdatePhotoAnalysisStatus(
		ctx,
		job.photoID,
		job.dest,
		model.AnalysisStatusProcessing,
		"",
	); err != nil {
		log.Printf("cannot mark photo %d as processing: %v", job.photoID, err)
		return
	}
	q.notify(job, model.AnalysisStatusProcessing, "")

	if err := q.photoService.GeneratePhotoEmbedding(ctx, job.photoID, job.dest); err != nil {
		if ctx.Err() != nil {
			// Leave it as processing so Start can recover it next launch.
			return
		}
		if updateErr := q.photoRepo.UpdatePhotoAnalysisStatus(
			ctx,
			job.photoID,
			job.dest,
			model.AnalysisStatusFailed,
			err.Error(),
		); updateErr != nil {
			log.Printf("cannot mark photo %d as failed: %v (original error: %v)", job.photoID, updateErr, err)
		}
		q.notify(job, model.AnalysisStatusFailed, err.Error())
		return
	}
	photo, err := q.photoRepo.GetPhotoByID(ctx, job.photoID, job.dest)
	if err != nil {
		log.Printf("cannot load completed photo %d: %v", job.photoID, err)
		q.notify(job, model.AnalysisStatusCompleted, "")
		return
	}
	tags := make([]string, 0, len(photo.Tags))
	for _, tag := range photo.Tags {
		tags = append(tags, tag.Name)
	}
	q.notifyResult(job, model.AnalysisStatusCompleted, "", photo.Description, tags)
}

func (q *EmbeddingQueue) notify(job embeddingJob, status string, message string) {
	q.notifyResult(job, status, message, "", nil)
}

func (q *EmbeddingQueue) notifyResult(job embeddingJob, status string, message string, description string, tags []string) {
	if q.onUpdate == nil {
		return
	}
	q.onUpdate(EmbeddingStatusUpdate{
		PhotoID:     job.photoID,
		Status:      status,
		Error:       message,
		Description: description,
		Tags:        tags,
	})
}
