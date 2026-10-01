package model

import "time"

type Photo struct {
	ID             uint   `gorm:"primaryKey"`
	Dest           string `gorm:"not null;uniqueIndex:idx_photo_dest_filename"`
	Filename       string `gorm:"not null;uniqueIndex:idx_photo_dest_filename"`
	Checksum       string `gorm:"not null;index"`
	Description    string `gorm:"type:text"`
	Status         string `gorm:"not null;default:active;index"`
	AnalysisStatus string `gorm:"not null;default:pending;index"`
	CreatedAt      time.Time
}
