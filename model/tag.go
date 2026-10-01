package model

type Tag struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null;uniqueIndex"`
}
