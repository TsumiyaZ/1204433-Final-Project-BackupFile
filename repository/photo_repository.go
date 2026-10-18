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
	GetPhotoByID(ctx context.Context, id uint, dest string) (*model.Photo, error)
	UpdatePhotoStatus(ctx context.Context, id uint, dest string, status string) error
	DeletePhoto(ctx context.Context, id uint, dest string) error
}

func NewPhotoRepository(db *gorm.DB) PhotoRepository {
	return &photoRepository{
		q: query.Use(db),
	}
}

func (r *photoRepository) DeletePhoto(ctx context.Context, id uint, dest string) error {
	_, err := r.q.Photo.WithContext(ctx).Where(r.q.Photo.ID.Eq(id), r.q.Photo.Dest.Eq(dest)).Delete()
	return err
}

func (r *photoRepository) UpdatePhotoStatus(ctx context.Context, id uint, dest string, status string) error {
	_, err := r.q.Photo.WithContext(ctx).Where(r.q.Photo.ID.Eq(id), r.q.Photo.Dest.Eq(dest)).Update(r.q.Photo.Status, status)
	return err
}

func (r *photoRepository) GetPhotoByID(ctx context.Context, id uint, dest string) (*model.Photo, error) {
	return r.q.Photo.WithContext(ctx).Where(r.q.Photo.ID.Eq(id), r.q.Photo.Dest.Eq(dest)).First()
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
