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
	UpdatePhotoEmbedding(ctx context.Context, id uint, dest string, embedding []byte, modelName string) error
	GetPhotosWithEmbedding(ctx context.Context, dest string) ([]*model.Photo, error)
	UpdatePhotoAnalysisStatus(
		ctx context.Context,
		id uint,
		dest string,
		status string,
		message string,
	) error
}

func NewPhotoRepository(db *gorm.DB) PhotoRepository {
	return &photoRepository{
		q: query.Use(db),
	}
}

func (r *photoRepository) UpdatePhotoAnalysisStatus(
	ctx context.Context,
	id uint,
	dest string,
	status string,
	message string,
) error {
	_, err := r.q.Photo.
		WithContext(ctx).
		Where(
			r.q.Photo.ID.Eq(id),
			r.q.Photo.Dest.Eq(dest),
		).
		UpdateSimple(
			r.q.Photo.AnalysisStatus.Value(status),
			r.q.Photo.AnalysisError.Value(""),
		)

	return err
}

func (r *photoRepository) GetPhotosWithEmbedding(
	ctx context.Context,
	dest string,
) ([]*model.Photo, error) {
	return r.q.Photo.
		WithContext(ctx).
		Where(
			r.q.Photo.Dest.Eq(dest),
			r.q.Photo.Status.Eq("active"),
			r.q.Photo.AnalysisStatus.Eq(model.AnalysisStatusCompleted),
			r.q.Photo.Embedding.IsNotNull(),
		).
		Find()
}

func (r *photoRepository) UpdatePhotoEmbedding(
	ctx context.Context,
	id uint,
	dest string,
	embedding []byte,
	modelName string,
) error {
	_, err := r.q.Photo.
		WithContext(ctx).
		Where(
			r.q.Photo.ID.Eq(id),
			r.q.Photo.Dest.Eq(dest),
		).
		UpdateSimple(
			r.q.Photo.Embedding.Value(embedding),
			r.q.Photo.EmbeddingModel.Value(modelName),
			r.q.Photo.AnalysisStatus.Value(model.AnalysisStatusCompleted),
			r.q.Photo.AnalysisError.Value(""),
		)

	return err
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
