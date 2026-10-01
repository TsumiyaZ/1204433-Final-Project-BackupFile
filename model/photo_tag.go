package model

type PhotoTag struct {
	PhotoID uint `gorm:"primaryKey"`
	TagID   uint `gorm:"primaryKey"`

	Photo Photo `gorm:"foreignKey:PhotoID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Tag   Tag   `gorm:"foreignKey:TagID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}
