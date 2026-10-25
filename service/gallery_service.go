package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"FinalProject/dto"
	"FinalProject/internal/fileopener"
	"FinalProject/model"
	"FinalProject/repository"
)

type galleryService struct {
	photoRepo    repository.PhotoRepository
	tagRepo      repository.TagRepository
	settingRepo  repository.SettingRepository
	photoService PhotoService
}

type GalleryService interface {
	GetGalleryPhoto(ctx context.Context) ([]dto.GalleryPhoto, error)
	GetGalleryPhotoPreview(ctx context.Context, photoID uint) (string, error)
	OpenGalleryPhoto(ctx context.Context, photoID uint) error
	ShowGalleryPhotoInFolder(ctx context.Context, photoID uint) error
	GenerateGalleryPhotoEmbedding(ctx context.Context, photoID uint) error
	DeleteGalleryPhoto(ctx context.Context, photoID uint) error
	SearchGalleryPhotos(ctx context.Context, text string) ([]dto.PhotoSearchResult, error)
	GetGalleryTagCloud(ctx context.Context) ([]dto.TagCloudItem, error)
}

func NewGalleryService(
	photoRepo repository.PhotoRepository,
	settingRepo repository.SettingRepository,
	photoService PhotoService,
	tagRepos ...repository.TagRepository,
) GalleryService {
	instance := &galleryService{
		photoRepo:    photoRepo,
		settingRepo:  settingRepo,
		photoService: photoService,
	}
	if len(tagRepos) > 0 {
		instance.tagRepo = tagRepos[0]
	}
	return instance
}

func (s *galleryService) GetGalleryTagCloud(ctx context.Context) ([]dto.TagCloudItem, error) {
	setting, err := s.settingRepo.GetSetting(ctx)
	if err != nil {
		return nil, err
	}
	if setting.Dest == "" || s.tagRepo == nil {
		return []dto.TagCloudItem{}, nil
	}
	return s.tagRepo.GetTagCloud(ctx, setting.Dest)
}

func (s *galleryService) SearchGalleryPhotos(
	ctx context.Context,
	text string,
) ([]dto.PhotoSearchResult, error) {
	setting, err := s.settingRepo.GetSetting(ctx)
	if err != nil {
		return nil, err
	}

	if setting.Dest == "" {
		return nil, fmt.Errorf("destination is not configured")
	}

	return s.photoService.SearchPhotos(ctx, setting.Dest, text)
}

func (s *galleryService) DeleteGalleryPhoto(ctx context.Context, photoID uint) error {
	setting, err := s.settingRepo.GetSetting(ctx)
	if err != nil {
		return err
	}

	if setting.Dest == "" {
		return fmt.Errorf("Destination is not configured")
	}

	photo, err := s.photoRepo.GetPhotoByID(ctx, photoID, setting.Dest)
	if err != nil {
		return err
	}

	if photo.Filename == "" ||
		photo.Filename == "." ||
		photo.Filename == ".." ||
		filepath.Base(photo.Filename) != photo.Filename {
		return fmt.Errorf("invalid photo filename")
	}

	path := filepath.Join(photo.Dest, photo.Filename)

	info, err := os.Lstat(path)

	if os.IsNotExist(err) {
	} else if err != nil {
		return fmt.Errorf("cannot access photo: %w", err)
	} else if !info.Mode().IsRegular() {
		return fmt.Errorf("selected path is not a regular file")
	}

	err = os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("cannot delete photo file: %w", err)
	}

	err = s.photoRepo.DeletePhoto(ctx, photo.ID, photo.Dest)
	if err != nil {
		return fmt.Errorf(
			"file is absent but database deletion failed; retry deletion: %w",
			err,
		)
	}

	return nil
}

func (s *galleryService) GetGalleryPhotoPreview(
	ctx context.Context,
	photoID uint,
) (string, error) {
	setting, err := s.settingRepo.GetSetting(ctx)
	if err != nil {
		return "", err
	}

	if setting.Dest == "" {
		return "", fmt.Errorf("destination is not configured")
	}

	photo, err := s.photoRepo.GetPhotoByID(ctx, photoID, setting.Dest)
	if err != nil {
		return "", err
	}

	path := filepath.Join(photo.Dest, photo.Filename)

	return s.photoService.GetPhotoPreview(ctx, setting.Dest, path)
}

func (s *galleryService) OpenGalleryPhoto(ctx context.Context, photoID uint) error {
	path, err := s.getGalleryPhotoPath(ctx, photoID)
	if err != nil {
		return err
	}
	return fileopener.Open(path)
}

func (s *galleryService) ShowGalleryPhotoInFolder(ctx context.Context, photoID uint) error {
	path, err := s.getGalleryPhotoPath(ctx, photoID)
	if err != nil {
		return err
	}
	return fileopener.Reveal(path)
}

func (s *galleryService) getGalleryPhotoPath(ctx context.Context, photoID uint) (string, error) {
	setting, err := s.settingRepo.GetSetting(ctx)
	if err != nil {
		return "", err
	}
	if setting.Dest == "" {
		return "", fmt.Errorf("destination is not configured")
	}

	photo, err := s.photoRepo.GetPhotoByID(ctx, photoID, setting.Dest)
	if err != nil {
		return "", err
	}
	if photo.Filename == "" || photo.Filename == "." || photo.Filename == ".." || filepath.Base(photo.Filename) != photo.Filename {
		return "", fmt.Errorf("invalid photo filename")
	}

	path := filepath.Join(photo.Dest, photo.Filename)
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("photo file no longer exists")
		}
		return "", fmt.Errorf("cannot access photo file: %w", err)
	}
	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("selected photo is not a regular file")
	}
	return filepath.Abs(path)
}

func (s *galleryService) GenerateGalleryPhotoEmbedding(
	ctx context.Context,
	photoID uint,
) error {
	setting, err := s.settingRepo.GetSetting(ctx)
	if err != nil {
		return err
	}

	if setting.Dest == "" {
		return fmt.Errorf("destination is not configured")
	}

	if err := s.photoRepo.UpdatePhotoAnalysisStatus(
		ctx,
		photoID,
		setting.Dest,
		model.AnalysisStatusProcessing,
		"",
	); err != nil {
		return fmt.Errorf("cannot mark photo as processing: %w", err)
	}

	if err := s.photoService.GeneratePhotoEmbedding(ctx, photoID, setting.Dest); err != nil {
		if updateErr := s.photoRepo.UpdatePhotoAnalysisStatus(
			ctx,
			photoID,
			setting.Dest,
			model.AnalysisStatusFailed,
			err.Error(),
		); updateErr != nil {
			return fmt.Errorf("%w; cannot save failed status: %v", err, updateErr)
		}
		return err
	}

	return nil
}

func (s *galleryService) GetGalleryPhoto(ctx context.Context) ([]dto.GalleryPhoto, error) {
	setting, err := s.settingRepo.GetSetting(ctx)
	if err != nil {
		return nil, err
	}

	if setting.Dest == "" {
		return []dto.GalleryPhoto{}, nil
	}

	photos, err := s.photoRepo.GetPhotosByDestination(ctx, setting.Dest)
	if err != nil {
		return nil, err
	}

	result := make([]dto.GalleryPhoto, 0, len(photos))

	for _, photo := range photos {
		path := filepath.Join(photo.Dest, photo.Filename)

		status := "active"

		info, err := os.Stat(path)

		if os.IsNotExist(err) {
			status = "missing"
		} else if err != nil {
			return nil, fmt.Errorf(
				"cannot check photo %q: %w",
				photo.Filename,
				err,
			)
		} else if !info.Mode().IsRegular() {
			status = "missing"
		}

		if photo.Status != status {
			err := s.photoRepo.UpdatePhotoStatus(ctx, photo.ID, photo.Dest, status)
			if err != nil {
				return nil, err
			}
		}

		result = append(result, dto.GalleryPhoto{
			ID:             photo.ID,
			Filename:       photo.Filename,
			Path:           path,
			Size:           fileSize(info, status),
			ModifiedAt:     fileModifiedAt(info, status),
			Status:         status,
			AnalysisStatus: photo.AnalysisStatus,
			AnalysisError:  photo.AnalysisError,
			Tags:           tagNames(photo.Tags),
		})
	}

	return result, nil

}

func tagNames(tags []model.Tag) []string {
	names := make([]string, 0, len(tags))
	for _, tag := range tags {
		names = append(names, tag.Name)
	}
	return names
}

func fileSize(info os.FileInfo, status string) int64 {
	if status != "active" || info == nil {
		return 0
	}
	return info.Size()
}

func fileModifiedAt(info os.FileInfo, status string) string {
	if status != "active" || info == nil {
		return ""
	}
	return info.ModTime().Format(time.RFC3339)
}
