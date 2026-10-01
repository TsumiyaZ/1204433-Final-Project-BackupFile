package model

type Setting struct {
	ID     uint   `gorm:"primaryKey"`
	Source string `gorm:"not null"`
	Dest   string `gorm:"not null"`
}
