package backup

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"FinalProject/dto"
	"FinalProject/internal/pathutil"
	"FinalProject/model"
	"FinalProject/repository"
)

// moveBatch holds state shared by all workers of one MovePhotos call.
type moveBatch struct {
	repo           repository.BackupRepository
	source         string
	dest           string
	embeddingQueue PhotoEmbeddingEnqueuer

	// mu guards the reservations below so two workers never pick the same
	// destination name or move two copies of the same content.
	mu        sync.Mutex
	checksums map[string]bool
	filenames map[string]bool
}

func (b *moveBatch) moveOne(ctx context.Context, rawPath string) dto.MoveItemResult {
	item := newItem(rawPath, dto.MoveStatusFailed, "")

	if ctx.Err() != nil {
		item.Status = dto.MoveStatusSkipped
		item.Message = "move was cancelled"
		return item
	}

	sourcePath, err := filepath.Abs(filepath.Clean(rawPath))
	if err != nil {
		item.Message = "invalid photo path"
		return item
	}
	item.SourcePath = sourcePath

	if !isInside(b.source, sourcePath) || isInside(b.dest, sourcePath) {
		item.Message = "photo is outside the source folder"
		return item
	}

	info, err := os.Stat(sourcePath)
	if err != nil {
		item.Message = "cannot access photo"
		return item
	}
	if !info.Mode().IsRegular() {
		item.Message = "path is not a file"
		return item
	}

	if err := validateImage(sourcePath); err != nil {
		item.Message = err.Error()
		return item
	}

	checksum, err := fileChecksum(sourcePath)
	if err != nil {
		item.Message = "cannot calculate checksum"
		return item
	}

	filename, existingPhoto, duplicate, err := b.reserve(ctx, checksum, filepath.Base(sourcePath))
	if err != nil {
		item.Message = err.Error()
		return item
	}
	if duplicate {
		item.Status = dto.MoveStatusSkipped
		item.Message = "same photo already exists in destination"
		return item
	}

	destPath := filepath.Join(b.dest, filename)
	item.Filename = filename
	item.DestinationPath = destPath

	if err := b.transfer(sourcePath, destPath, info.Size(), checksum); err != nil {
		b.release(checksum, filename)
		item.Message = err.Error()
		return item
	}

	photo := existingPhoto
	needsAnalysis := true
	if photo == nil {
		photo = &model.Photo{
			Dest:     b.dest,
			Filename: filename,
			Checksum: checksum,
		}
		if err := b.repo.CreatePhoto(ctx, photo); err != nil {
			// The file is already safe in Destination. Report it so the
			// integrity check can reconcile the missing record later.
			item.Message = "photo moved but saving the record failed"
			return item
		}
	} else {
		needsAnalysis = photo.AnalysisStatus != model.AnalysisStatusCompleted || len(photo.Embedding) == 0
		if err := b.repo.RestorePhoto(ctx, photo.ID, photo.Dest, needsAnalysis); err != nil {
			item.Message = "photo restored but updating its database record failed"
			return item
		}
	}
	if b.embeddingQueue != nil && needsAnalysis {
		b.embeddingQueue.Enqueue(photo.ID, photo.Dest)
	}

	item.Status = dto.MoveStatusMoved
	item.Message = "moved"

	return item
}

// reserve checks for duplicates and claims a free destination filename.
// It returns duplicate=true when the same content already exists.
func (b *moveBatch) reserve(
	ctx context.Context,
	checksum string,
	filename string,
) (string, *model.Photo, bool, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.checksums[checksum] {
		return "", nil, true, nil
	}

	trackedPhotos, err := b.repo.PhotosByChecksum(ctx, b.dest, checksum)
	if err != nil {
		return "", nil, false, fmt.Errorf("cannot check duplicate photos")
	}

	var missingPhoto *model.Photo
	for _, photo := range trackedPhotos {
		if photo == nil {
			continue
		}
		photoPath := filepath.Join(b.dest, photo.Filename)
		info, statErr := os.Stat(photoPath)
		if os.IsNotExist(statErr) {
			if missingPhoto == nil {
				missingPhoto = photo
			}
			continue
		}
		if statErr != nil {
			return "", nil, false, fmt.Errorf("cannot check existing photo")
		}
		if !info.Mode().IsRegular() {
			continue
		}
		existingChecksum, checksumErr := fileChecksum(photoPath)
		if checksumErr != nil {
			return "", nil, false, fmt.Errorf("cannot check existing photo")
		}
		if existingChecksum == checksum {
			return "", nil, true, nil
		}
	}

	if missingPhoto != nil {
		key := pathutil.Key(missingPhoto.Filename)
		if b.filenames[key] {
			return "", nil, false, fmt.Errorf("cannot reserve filename for restored photo")
		}
		b.checksums[checksum] = true
		b.filenames[key] = true
		return missingPhoto.Filename, missingPhoto, false, nil
	}

	extension := filepath.Ext(filename)
	baseName := strings.TrimSuffix(filename, extension)

	for attempt := 0; attempt < 1000; attempt++ {
		candidate := filename
		if attempt > 0 {
			candidate = fmt.Sprintf("%s (%d)%s", baseName, attempt, extension)
		}

		key := pathutil.Key(candidate)
		if b.filenames[key] {
			continue
		}

		candidatePath := filepath.Join(b.dest, candidate)
		if _, err := os.Stat(candidatePath); err == nil {
			// A file not tracked in the database may still be the same photo.
			existingChecksum, err := fileChecksum(candidatePath)
			if err == nil && existingChecksum == checksum {
				return "", nil, true, nil
			}
			continue
		} else if !os.IsNotExist(err) {
			return "", nil, false, fmt.Errorf("cannot check destination file")
		}

		taken, err := b.repo.FilenameExists(ctx, b.dest, candidate)
		if err != nil {
			return "", nil, false, fmt.Errorf("cannot check destination filename")
		}
		if taken {
			continue
		}

		b.checksums[checksum] = true
		b.filenames[key] = true

		return candidate, nil, false, nil
	}

	return "", nil, false, fmt.Errorf("cannot find a free filename in destination")
}

func (b *moveBatch) release(checksum string, filename string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	delete(b.checksums, checksum)
	delete(b.filenames, pathutil.Key(filename))
}

// transfer renames the file when possible and falls back to
// copy + verify + delete when Source and Destination are on different drives.
func (b *moveBatch) transfer(
	sourcePath string,
	destPath string,
	size int64,
	checksum string,
) error {
	if sameVolume(sourcePath, destPath) {
		err := os.Rename(sourcePath, destPath)
		if err == nil {
			info, statErr := os.Stat(destPath)
			if statErr != nil || info.Size() != size {
				return fmt.Errorf("moved file could not be verified")
			}
			return nil
		}

		if !isCrossDeviceError(err) {
			return fmt.Errorf("cannot move photo")
		}
	}

	return copyVerifyDelete(sourcePath, destPath, size, checksum)
}

func copyVerifyDelete(
	sourcePath string,
	destPath string,
	size int64,
	checksum string,
) error {
	tempPath, err := copyToTemp(sourcePath, filepath.Dir(destPath))
	if err != nil {
		return err
	}

	// Remove the temporary file on every failure path below.
	keepTemp := false
	defer func() {
		if !keepTemp {
			os.Remove(tempPath)
		}
	}()

	info, err := os.Stat(tempPath)
	if err != nil || info.Size() != size {
		return fmt.Errorf("copied file size does not match")
	}

	copiedChecksum, err := fileChecksum(tempPath)
	if err != nil || copiedChecksum != checksum {
		return fmt.Errorf("copied file checksum does not match")
	}

	if _, err := os.Stat(destPath); err == nil {
		return fmt.Errorf("destination file already exists")
	}

	if err := os.Rename(tempPath, destPath); err != nil {
		return fmt.Errorf("cannot finalize copied file")
	}
	keepTemp = true

	// Source is deleted only after Destination passed verification.
	if err := os.Remove(sourcePath); err != nil {
		return fmt.Errorf("photo copied but source file could not be deleted")
	}

	return nil
}

func copyToTemp(sourcePath string, destDir string) (string, error) {
	source, err := os.Open(sourcePath)
	if err != nil {
		return "", fmt.Errorf("cannot open photo")
	}
	defer source.Close()

	temp, err := os.CreateTemp(destDir, ".backup-*.tmp")
	if err != nil {
		return "", fmt.Errorf("cannot create temporary file")
	}
	tempPath := temp.Name()

	_, copyErr := io.Copy(temp, source)
	syncErr := temp.Sync()
	closeErr := temp.Close()

	if copyErr != nil || syncErr != nil || closeErr != nil {
		os.Remove(tempPath)
		return "", fmt.Errorf("cannot copy photo")
	}

	return tempPath, nil
}

func sameVolume(sourcePath string, destPath string) bool {
	return strings.EqualFold(filepath.VolumeName(sourcePath), filepath.VolumeName(destPath))
}
