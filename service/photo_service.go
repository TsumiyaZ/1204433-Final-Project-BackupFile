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

const embeddingModelName = "embeddinggemma-2"
const minimumSearchScore = 0.65

type photoService struct {
	repo   repository.PhotoRepository
	engine EmbeddingEngine
}

type PhotoService interface {
	ScanPhotos(
		ctx context.Context,
		source string,
	) ([]dto.ScannedPhoto, error)
	GetPhotoPreview(ctx context.Context, source string, path string) (string, error)
	GeneratePhotoEmbedding(ctx context.Context, photoID uint, dest string) error
	SearchPhotos(ctx context.Context, dest string, text string) ([]dto.PhotoSearchResult, error)
}

func NewPhotoService(
	repo repository.PhotoRepository,
	engine EmbeddingEngine,
) PhotoService {
	return &photoService{
		repo:   repo,
		engine: engine,
	}
}

func (s *photoService) SearchPhotos(
	ctx context.Context,
	dest string,
	text string,
) ([]dto.PhotoSearchResult, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, fmt.Errorf("search text is required")
	}

	queryVector, err := s.engine.EmbedText(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("cannot embed search text: %w", err)
	}

	photos, err := s.repo.GetPhotosWithEmbedding(ctx, dest)
	if err != nil {
		return nil, fmt.Errorf("cannot get photos with embeddings: %w", err)
	}

	results := make([]dto.PhotoSearchResult, 0, len(photos))
	for _, photo := range photos {
		photoVector, err := decodeEmbedding(photo.Embedding)
		if err != nil {
			return nil, fmt.Errorf("cannot decode embedding for %q: %w", photo.Filename, err)
		}

		score, err := cosineSimilarity(queryVector, photoVector)
		if err != nil {
			return nil, fmt.Errorf("cannot compare embedding for %q: %w", photo.Filename, err)
		}
		if score < minimumSearchScore {
			continue
		}

		results = append(results, dto.PhotoSearchResult{
			ID:       photo.ID,
			Filename: photo.Filename,
			Score:    score,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results, nil
}

func (s *photoService) GeneratePhotoEmbedding(
	ctx context.Context,
	photoID uint,
	dest string,
) error {
	photo, err := s.repo.GetPhotoByID(ctx, photoID, dest)
	if err != nil {
		return fmt.Errorf("cannot find photo: %w", err)
	}

	imagePath := filepath.Join(photo.Dest, photo.Filename)

	vector, err := s.engine.EmbedImage(ctx, imagePath)
	if err != nil {
		return fmt.Errorf("cannot create image embedding: %w", err)
	}

	embeddingBytes, err := encodeEmbedding(vector)
	if err != nil {
		return fmt.Errorf("cannot encode image embedding: %w", err)
	}

	if err := s.repo.UpdatePhotoEmbedding(
		ctx,
		photo.ID,
		photo.Dest,
		embeddingBytes,
		embeddingModelName,
	); err != nil {
		return fmt.Errorf("cannot save image embedding: %w", err)
	}

	return nil
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

	if !isPreviewableImage(photoPath) {
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

			if entry.IsDir() || !isScannableImage(path) {
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

func isScannableImage(path string) bool {
	extension := strings.ToLower(filepath.Ext(path))

	switch extension {
	case ".jpg", ".jpeg", ".png":
		return true
	default:
		return false
	}
}

func isPreviewableImage(path string) bool {
	extension := strings.ToLower(filepath.Ext(path))

	switch extension {
	case ".jpg", ".jpeg", ".png", ".webp", ".gif":
		return true
	default:
		return false
	}
}
