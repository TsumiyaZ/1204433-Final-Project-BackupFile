package backup

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"FinalProject/dto"
	"FinalProject/internal/pathutil"
	"FinalProject/repository"
)

const defaultWorkerCount = 4

// ProgressFunc receives one update after each file finishes.
// It keeps the service independent from Wails events.
type ProgressFunc func(progress dto.MoveProgress)

// PhotoEmbeddingEnqueuer queues a photo after its file and database record
// have been saved. The backup package does not depend on the embedding service.
type PhotoEmbeddingEnqueuer interface {
	Enqueue(photoID uint, dest string)
}

type BackupService interface {
	MovePhotos(
		ctx context.Context,
		request dto.MovePhotosRequest,
		onProgress ProgressFunc,
	) (*dto.MovePhotosResult, error)
}

type backupService struct {
	repo           repository.BackupRepository
	settingRepo    repository.SettingRepository
	embeddingQueue PhotoEmbeddingEnqueuer
	workerCount    int
}

func NewBackupService(
	repo repository.BackupRepository,
	settingRepo repository.SettingRepository,
	embeddingQueues ...PhotoEmbeddingEnqueuer,
) BackupService {
	var embeddingQueue PhotoEmbeddingEnqueuer
	if len(embeddingQueues) > 0 {
		embeddingQueue = embeddingQueues[0]
	}

	return &backupService{
		repo:           repo,
		settingRepo:    settingRepo,
		embeddingQueue: embeddingQueue,
		workerCount:    defaultWorkerCount,
	}
}

type moveJob struct {
	index int
	path  string
}

type moveDone struct {
	index int
	item  dto.MoveItemResult
}

func (s *backupService) MovePhotos(
	ctx context.Context,
	request dto.MovePhotosRequest,
	onProgress ProgressFunc,
) (*dto.MovePhotosResult, error) {
	startedAt := time.Now()

	source, dest, err := s.loadFolders(ctx)
	if err != nil {
		return nil, err
	}

	paths := uniquePaths(request.Paths)
	if len(paths) == 0 {
		return nil, errors.New("no photos selected")
	}

	batch := &moveBatch{
		repo:           s.repo,
		source:         source,
		dest:           dest,
		embeddingQueue: s.embeddingQueue,
		checksums:      make(map[string]bool),
		filenames:      make(map[string]bool),
	}

	jobs := make(chan moveJob)
	results := make(chan moveDone)

	workerCount := min(s.workerCount, len(paths))

	var wg sync.WaitGroup
	for range workerCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				results <- moveDone{
					index: job.index,
					item:  batch.moveOne(ctx, job.path),
				}
			}
		}()
	}

	go func() {
		defer close(jobs)
		for index, path := range paths {
			select {
			case <-ctx.Done():
				return
			case jobs <- moveJob{index: index, path: path}:
			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	result := &dto.MovePhotosResult{
		Requested: len(paths),
		Items:     make([]dto.MoveItemResult, len(paths)),
	}
	finished := make([]bool, len(paths))
	current := 0

	for done := range results {
		current++
		finished[done.index] = true
		result.Items[done.index] = done.item

		if onProgress != nil {
			onProgress(dto.MoveProgress{
				Current: current,
				Total:   len(paths),
				Percent: current * 100 / len(paths),
				Item:    done.item,
			})
		}
	}

	// Files never picked up because the move was cancelled stay in Source.
	for index, path := range paths {
		if !finished[index] {
			result.Items[index] = newItem(path, dto.MoveStatusSkipped, "move was cancelled")
		}
	}

	for _, item := range result.Items {
		switch item.Status {
		case dto.MoveStatusMoved:
			result.Moved++
		case dto.MoveStatusSkipped:
			result.Skipped++
		default:
			result.Failed++
		}
	}

	result.DurationMs = time.Since(startedAt).Milliseconds()

	return result, nil
}

func (s *backupService) loadFolders(ctx context.Context) (string, string, error) {
	setting, err := s.settingRepo.GetSetting(ctx)
	if err != nil {
		return "", "", fmt.Errorf("cannot read folder setting: %w", err)
	}

	if strings.TrimSpace(setting.Source) == "" || strings.TrimSpace(setting.Dest) == "" {
		return "", "", errors.New("source and destination must be set before moving photos")
	}

	source, err := filepath.Abs(filepath.Clean(setting.Source))
	if err != nil {
		return "", "", fmt.Errorf("invalid source folder: %w", err)
	}

	dest, err := filepath.Abs(filepath.Clean(setting.Dest))
	if err != nil {
		return "", "", fmt.Errorf("invalid destination folder: %w", err)
	}

	if pathutil.Same(source, dest) {
		return "", "", errors.New("source and destination must be different")
	}

	for _, folder := range []string{source, dest} {
		info, err := os.Stat(folder)
		if err != nil {
			return "", "", fmt.Errorf("cannot access folder %s: %w", folder, err)
		}
		if !info.IsDir() {
			return "", "", fmt.Errorf("%s is not a folder", folder)
		}
	}

	return source, dest, nil
}

func uniquePaths(paths []string) []string {
	seen := make(map[string]bool, len(paths))
	unique := make([]string, 0, len(paths))

	for _, path := range paths {
		path = strings.TrimSpace(path)
		if path == "" {
			continue
		}

		key := pathutil.Key(path)
		if seen[key] {
			continue
		}

		seen[key] = true
		unique = append(unique, path)
	}

	return unique
}

func newItem(path string, status string, message string) dto.MoveItemResult {
	return dto.MoveItemResult{
		SourcePath: path,
		Filename:   filepath.Base(path),
		Status:     status,
		Message:    message,
	}
}
