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
