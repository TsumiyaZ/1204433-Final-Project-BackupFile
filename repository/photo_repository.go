package repository

import (
	"context"

	"FinalProject/model"
	"FinalProject/model/query"

	"gorm.io/gorm"
)

type photoRepository struct {
	q *query.Query
}

type PhotoRepository interface {
	CreatePhoto(ctx context.Context, photo *model.Photo) error
	GetPhotosByDestination(ctx context.Context, dest string) ([]*model.Photo, error)
	PhotoExists(ctx context.Context, dest string, filename string) (bool, error)
}

func NewPhotoRepository(db *gorm.DB) PhotoRepository {
	return &photoRepository{
		q: query.Use(db),
	}
}

func (r *photoRepository) CreatePhoto(ctx context.Context, photo *model.Photo) error {
	return r.q.Photo.WithContext(ctx).Create(photo)
}

func (r *photoRepository) GetPhotosByDestination(ctx context.Context, dest string) ([]*model.Photo, error) {
	return r.q.Photo.WithContext(ctx).Where(r.q.Photo.Dest.Eq(dest)).Order(r.q.Photo.CreatedAt.Desc()).Find()
}

func (r *photoRepository) PhotoExists(
	ctx context.Context,
	dest string,
	filename string,
) (bool, error) {
	count, err := r.q.Photo.
		WithContext(ctx).
		Where(
			r.q.Photo.Dest.Eq(dest),
			r.q.Photo.Filename.Eq(filename),
		).
		Count()

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
