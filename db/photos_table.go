package db

import (
	"gorm.io/gorm"
)

type Photo struct {
	gorm.Model
	Name       string
	Preview    string
	PreviewURL string
	ImageURL   string
}

func InitMigration() {
	GetDBConnection().AutoMigrate(&Photo{})
}
