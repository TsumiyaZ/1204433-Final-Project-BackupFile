package repository

import (
	"FinalProject/model"

	"gorm.io/gorm"
)

func MigrateDatabase(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Setting{},
		&model.Photo{},
		&model.Tag{},
		&model.PhotoTag{},
	)
}
