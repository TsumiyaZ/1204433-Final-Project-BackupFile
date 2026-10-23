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

func (c *GalleryController) GenerateGalleryPhotoEmbedding(photoID uint) error {
	return c.service.GenerateGalleryPhotoEmbedding(context.Background(), photoID)
}

func (c *GalleryController) DeleteGalleryPhoto(photoID uint) error {
	return c.service.DeleteGalleryPhoto(context.Background(), photoID)
}

func (c *GalleryController) SearchGalleryPhotos(
	text string,
) ([]dto.PhotoSearchResult, error) {
	return c.service.SearchGalleryPhotos(
		context.Background(),
		text,
	)
}
