package repository

import (
	"context"

	"FinalProject/model"
	"FinalProject/model/query"

	"gorm.io/gorm"
)

type backupRepository struct {
	q *query.Query
}

type BackupRepository interface {
	ChecksumExists(ctx context.Context, dest string, checksum string) (bool, error)
	FilenameExists(ctx context.Context, dest string, filename string) (bool, error)
	PhotosByChecksum(ctx context.Context, dest string, checksum string) ([]*model.Photo, error)
	RestorePhoto(ctx context.Context, id uint, dest string, needsAnalysis bool) error
	CreatePhoto(ctx context.Context, photo *model.Photo) error
}

func NewBackupRepository(db *gorm.DB) BackupRepository {
	return &backupRepository{
		q: query.Use(db),
	}
}

func (r *backupRepository) ChecksumExists(
	ctx context.Context,
	dest string,
	checksum string,
) (bool, error) {
	count, err := r.q.Photo.
		WithContext(ctx).
		Where(
			r.q.Photo.Dest.Eq(dest),
			r.q.Photo.Checksum.Eq(checksum),
		).
		Count()

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *backupRepository) FilenameExists(
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

func (r *backupRepository) PhotosByChecksum(
	ctx context.Context,
	dest string,
	checksum string,
) ([]*model.Photo, error) {
	return r.q.Photo.WithContext(ctx).
		Where(r.q.Photo.Dest.Eq(dest), r.q.Photo.Checksum.Eq(checksum)).
		Find()
}

func (r *backupRepository) RestorePhoto(
	ctx context.Context,
	id uint,
	dest string,
	needsAnalysis bool,
) error {
	updates := map[string]interface{}{"status": "active"}
	if needsAnalysis {
		updates["analysis_status"] = model.AnalysisStatusPending
		updates["analysis_error"] = ""
	}
	_, err := r.q.Photo.WithContext(ctx).
		Where(r.q.Photo.ID.Eq(id), r.q.Photo.Dest.Eq(dest)).
		Updates(updates)
	return err
}

func (r *backupRepository) CreatePhoto(ctx context.Context, photo *model.Photo) error {
	return r.q.Photo.WithContext(ctx).Create(photo)
}
