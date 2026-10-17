package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"FinalProject/dto"
	"FinalProject/repository"
)

type galleryService struct {
	photoRepo    repository.PhotoRepository
	settingRepo  repository.SettingRepository
	photoService PhotoService
}

type GalleryService interface {
	GetGalleryPhoto(ctx context.Context) ([]dto.GalleryPhoto, error)
	GetGalleryPhotoPreview(ctx context.Context, photoID uint) (string, error)
}

func NewGalleryService(
	photoRepo repository.PhotoRepository,
	settingRepo repository.SettingRepository,
	photoService PhotoService,
) GalleryService {
	return &galleryService{
		photoRepo:    photoRepo,
		settingRepo:  settingRepo,
		photoService: photoService,
	}
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
			Description:    photo.Description,
			Status:         status,
			AnalysisStatus: photo.AnalysisStatus,
		})
	}

	return result, nil

}
