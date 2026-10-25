package repository

import (
	"context"

	"FinalProject/dto"

	"gorm.io/gorm"
)

type TagRepository interface {
	GetTagCloud(ctx context.Context, dest string) ([]dto.TagCloudItem, error)
}

type tagRepository struct{ db *gorm.DB }

func NewTagRepository(db *gorm.DB) TagRepository { return &tagRepository{db: db} }

func (r *tagRepository) GetTagCloud(ctx context.Context, dest string) ([]dto.TagCloudItem, error) {
	items := make([]dto.TagCloudItem, 0)
	err := r.db.WithContext(ctx).
		Table("tags").
		Select("tags.name AS name, COUNT(DISTINCT photos.id) AS photo_count").
		Joins("JOIN photo_tags ON photo_tags.tag_id = tags.id").
		Joins("JOIN photos ON photos.id = photo_tags.photo_id").
		Where("photos.dest = ? AND photos.status = ?", dest, "active").
		Group("tags.id, tags.name").
		Order("photo_count DESC, tags.name COLLATE NOCASE ASC").
		Scan(&items).Error
	return items, err
}
