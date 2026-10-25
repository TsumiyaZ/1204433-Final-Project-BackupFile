package controller

import (
	"context"

	"FinalProject/dto"
	"FinalProject/service"
)

type GalleryController struct {
	service        service.GalleryService
	embeddingQueue *service.EmbeddingQueue
}

func NewGalleryController(
	service service.GalleryService,
	queues ...*service.EmbeddingQueue,
) *GalleryController {
	controller := &GalleryController{service: service}
	if len(queues) > 0 {
		controller.embeddingQueue = queues[0]
	}
	return controller
}

func (c *GalleryController) PauseEmbeddingQueue() {
	if c.embeddingQueue != nil {
		c.embeddingQueue.Pause()
	}
}

func (c *GalleryController) ResumeEmbeddingQueue() {
	if c.embeddingQueue != nil {
		c.embeddingQueue.Resume()
	}
}

func (c *GalleryController) IsEmbeddingQueuePaused() bool {
	return c.embeddingQueue != nil && c.embeddingQueue.IsPaused()
}

func (c *GalleryController) GetGalleryPhotos() ([]dto.GalleryPhoto, error) {
	return c.service.GetGalleryPhoto(context.Background())
}

func (c *GalleryController) GetGalleryTagCloud() ([]dto.TagCloudItem, error) {
	return c.service.GetGalleryTagCloud(context.Background())
}

func (c *GalleryController) GetGalleryPhotoPreview(photoID uint) (string, error) {
	return c.service.GetGalleryPhotoPreview(context.Background(), photoID)
}

func (c *GalleryController) OpenGalleryPhoto(photoID uint) error {
	return c.service.OpenGalleryPhoto(context.Background(), photoID)
}

func (c *GalleryController) ShowGalleryPhotoInFolder(photoID uint) error {
	return c.service.ShowGalleryPhotoInFolder(context.Background(), photoID)
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
