package repository

import (
	"context"

	"FinalProject/model"
	"FinalProject/model/query"

	"gorm.io/gorm"
)

type photoRepository struct {
	q  *query.Query
	db *gorm.DB
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
	GetPhotosNeedingAnalysis(ctx context.Context, dest string) ([]*model.Photo, error)
	SavePhotoAnalysis(ctx context.Context, id uint, dest string, embedding []byte, modelName string, description string, tags []string) error
}

func NewPhotoRepository(db *gorm.DB) PhotoRepository {
	return &photoRepository{
		q:  query.Use(db),
		db: db,
	}
}

func (r *photoRepository) UpdatePhotoAnalysisStatus(
	ctx context.Context,
	id uint,
	dest string,
	status string,
	message string,
) error {
	return r.db.WithContext(ctx).Model(&model.Photo{}).
		Where("id = ? AND dest = ?", id, dest).
		Updates(map[string]interface{}{"analysis_status": status, "analysis_error": message}).Error
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

func (r *photoRepository) GetPhotosNeedingAnalysis(ctx context.Context, dest string) ([]*model.Photo, error) {
	photos := make([]*model.Photo, 0)
	err := r.db.WithContext(ctx).Where("dest = ? AND status = ? AND (analysis_status IN ? OR description IS NULL OR trim(description) = '' OR NOT EXISTS (SELECT 1 FROM photo_tags WHERE photo_tags.photo_id = photos.id))", dest, "active", []string{model.AnalysisStatusPending, model.AnalysisStatusProcessing}).Find(&photos).Error
	return photos, err
}

func (r *photoRepository) SavePhotoAnalysis(ctx context.Context, id uint, dest string, embedding []byte, modelName string, description string, tags []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM photo_tags WHERE photo_id = ?", id).Error; err != nil {
			return err
		}
		for _, name := range tags {
			if err := tx.Exec("INSERT INTO tags(name) VALUES (?) ON CONFLICT(name) DO NOTHING", name).Error; err != nil {
				return err
			}
			var tagID uint
			if err := tx.Raw("SELECT id FROM tags WHERE name = ?", name).Scan(&tagID).Error; err != nil {
				return err
			}
			if err := tx.Exec("INSERT INTO photo_tags(photo_id, tag_id) VALUES (?, ?) ON CONFLICT(photo_id, tag_id) DO NOTHING", id, tagID).Error; err != nil {
				return err
			}
		}
		return tx.Model(&model.Photo{}).Where("id = ? AND dest = ?", id, dest).Updates(map[string]interface{}{
			"embedding":       embedding,
			"embedding_model": modelName,
			"description":     description,
			"analysis_status": model.AnalysisStatusCompleted,
			"analysis_error":  "",
		}).Error
	})
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
	photo := new(model.Photo)
	err := r.db.WithContext(ctx).Preload("Tags").Where("id = ? AND dest = ?", id, dest).First(photo).Error
	return photo, err
}

func (r *photoRepository) CreatePhoto(ctx context.Context, photo *model.Photo) error {
	return r.q.Photo.WithContext(ctx).Create(photo)
}

func (r *photoRepository) GetPhotosByDestination(ctx context.Context, dest string) ([]*model.Photo, error) {
	photos := make([]*model.Photo, 0)
	err := r.db.WithContext(ctx).Preload("Tags").Where("dest = ?", dest).Order("created_at DESC").Find(&photos).Error
	return photos, err
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
