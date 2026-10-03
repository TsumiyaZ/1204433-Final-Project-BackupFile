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

	"encoding/base64"
	"net/http"
)

type photoService struct {
	repo repository.PhotoRepository
}

type PhotoService interface {
	ScanPhotos(
		ctx context.Context,
		source string,
	) ([]dto.ScannedPhoto, error)
	GetPhotoPreview(ctx context.Context, source string, path string) (string, error)
}

func NewPhotoService(
	repo repository.PhotoRepository,
) PhotoService {
	return &photoService{
		repo: repo,
	}
}

func (s *photoService) GetPhotoPreview(ctx context.Context, source string, path string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	sourcePath, err := filepath.Abs(filepath.Clean(source))
	if err != nil {
		return "", fmt.Errorf("invalid folder path: %w", err)
	}

	photoPath, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return "", fmt.Errorf("invalid photo path: %w", err)
	}

	relativePath, err := filepath.Rel(sourcePath, photoPath)
	if err != nil {
		return "", fmt.Errorf("cannot validate photo path: %w", err)
	}

	if relativePath == ".." ||
		strings.HasPrefix(relativePath, ".."+string(os.PathSeparator)) ||
		filepath.IsAbs(relativePath) {
		return "", fmt.Errorf("photo is outside selected folder")
	}

	if !isSupportedImage(photoPath) {
		return "", fmt.Errorf("unsupported image type")
	}

	info, err := os.Stat(photoPath)
	if err != nil {
		return "", fmt.Errorf("cannot access photo: %w", err)
	}

	const maxPreviewSize = 20 * 1024 * 1024

	if info.Size() > maxPreviewSize {
		return "", fmt.Errorf("photo is too large for preview")
	}

	data, err := os.ReadFile(photoPath)
	if err != nil {
		return "", fmt.Errorf("cannot read photo: %w", err)
	}

	contentType := http.DetectContentType(data)

	if !strings.HasPrefix(contentType, "image/") {
		return "", fmt.Errorf("file content is not an image")
	}

	encoded := base64.StdEncoding.EncodeToString(data)

	return "data:" + contentType + ";base64," + encoded, nil
}

func (s *photoService) ScanPhotos(
	ctx context.Context,
	source string,
) ([]dto.ScannedPhoto, error) {
	source = strings.TrimSpace(source)

	if source == "" {
		return nil, fmt.Errorf("folder is required")
	}

	source, err := filepath.Abs(filepath.Clean(source))
	if err != nil {
		return nil, fmt.Errorf("cannot resolve folder path: %w", err)
	}

	sourceInfo, err := os.Stat(source)
	if err != nil {
		return nil, fmt.Errorf("cannot access folder: %w", err)
	}

	if !sourceInfo.IsDir() {
		return nil, fmt.Errorf("selected path is not a folder")
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
		return nil, fmt.Errorf("cannot scan folder: %w", err)
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
