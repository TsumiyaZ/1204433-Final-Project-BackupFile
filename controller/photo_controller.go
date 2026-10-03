package controller

import (
	"context"

	"FinalProject/dto"
	"FinalProject/service"
)

type PhotoController struct {
	service service.PhotoService
}

func NewPhotoController(
	service service.PhotoService,
) *PhotoController {
	return &PhotoController{
		service: service,
	}
}

func (c *PhotoController) ScanPhotos(
	source string,
) ([]dto.ScannedPhoto, error) {
	return c.service.ScanPhotos(
		context.Background(),
		source,
	)
}

func (c *PhotoController) GetPhotoPreview(source string, path string) (string, error) {
	return c.service.GetPhotoPreview(context.Background(), source, path)
}
