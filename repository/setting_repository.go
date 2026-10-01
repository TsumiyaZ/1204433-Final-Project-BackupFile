package repository

import (
	"context"

	"FinalProject/model"
	"FinalProject/model/query"

	"gorm.io/gorm"
)

type settingRepository struct {
	q *query.Query
}

type SettingRepository interface {
	GetSetting(ctx context.Context) (*model.Setting, error)
	SaveSetting(ctx context.Context, setting *model.Setting) error
}

func NewSettingRepository(db *gorm.DB) SettingRepository {
	return &settingRepository{
		q: query.Use(db),
	}
}

func (r *settingRepository) GetSetting(ctx context.Context) (*model.Setting, error) {
	return r.q.Setting.WithContext(ctx).Where(r.q.Setting.ID.Eq(1)).First()
}

func (r *settingRepository) SaveSetting(ctx context.Context, setting *model.Setting) error {
	setting.ID = 1
	return r.q.Setting.WithContext(ctx).Save(setting)
}
