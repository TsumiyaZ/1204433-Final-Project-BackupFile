package service

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"FinalProject/dto"
	"FinalProject/repository"
)

type photoService struct {
	repo repository.PhotoRepository
}

type PhotoService interface {
	ScanPhotos(
		ctx context.Context,
		source string,
	) ([]dto.ScannedPhoto, error)
}

func NewPhotoService(
	repo repository.PhotoRepository,
) PhotoService {
	return &photoService{
		repo: repo,
	}
}

func (s *photoService) ScanPhotos(
	ctx context.Context,
	source string,
) ([]dto.ScannedPhoto, error) {
	source = strings.TrimSpace(source)

	if source == "" {
		return nil, fmt.Errorf("source folder is required")
	}

	source, err := filepath.Abs(filepath.Clean(source))
	if err != nil {
		return nil, fmt.Errorf("cannot resolve source path: %w", err)
	}

	sourceInfo, err := os.Stat(source)
	if err != nil {
		return nil, fmt.Errorf("cannot access source folder: %w", err)
	}

	if !sourceInfo.IsDir() {
		return nil, fmt.Errorf("source path is not a folder")
	}

	photos := make([]dto.ScannedPhoto, 0)

	err = filepath.WalkDir(
		source,
		func(path string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			if entry.IsDir() || !isSupportedImage(path) {
				return nil
			}

			info, err := entry.Info()
			if err != nil {
				return err
			}

			relativePath, err := filepath.Rel(source, path)
			if err != nil {
				return err
			}

			photos = append(photos, dto.ScannedPhoto{
				Filename:     entry.Name(),
				Path:         path,
				RelativePath: relativePath,
				Extension:    strings.ToLower(filepath.Ext(path)),
				Size:         info.Size(),
			})

			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("cannot scan source folder: %w", err)
	}

	sort.Slice(photos, func(i, j int) bool {
		return photos[i].RelativePath < photos[j].RelativePath
	})

	return photos, nil
}

func isSupportedImage(path string) bool {
	extension := strings.ToLower(filepath.Ext(path))

	switch extension {
	case ".jpg", ".jpeg", ".png", ".webp", ".gif":
		return true
	default:
		return false
	}
}
