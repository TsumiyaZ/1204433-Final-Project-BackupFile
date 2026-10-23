package service

import "context"

type EmbeddingEngine interface {
	EmbedImage(ctx context.Context, imagePath string) ([]float32, error)
	EmbedText(ctx context.Context, text string) ([]float32, error)
}
