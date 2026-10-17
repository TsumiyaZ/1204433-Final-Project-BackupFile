package controller

import (
	"context"

	"FinalProject/dto"
	"FinalProject/service"
)

type GalleryController struct {
	service service.GalleryService
}

func NewGalleryController(
	service service.GalleryService,
) *GalleryController {
	return &GalleryController{
		service: service,
	}
}

func (c *GalleryController) GetGalleryPhotos() ([]dto.GalleryPhoto, error) {
	return c.service.GetGalleryPhoto(context.Background())
}

func (c *GalleryController) GetGalleryPhotoPreview(photoID uint) (string, error) {
	return c.service.GetGalleryPhotoPreview(context.Background(), photoID)
}
