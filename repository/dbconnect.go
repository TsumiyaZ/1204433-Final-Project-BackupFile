package repository

import (
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func getDatabasePath() (string, error) {
	projectDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	databasePath := filepath.Join(projectDir, "database")

	err = os.MkdirAll(databasePath, 0755)
	if err != nil {
		return "", err
	}

	dbPath := filepath.Join(databasePath, "database.db")

	return dbPath, nil
}

func NewDBConnection() (*gorm.DB, error) {
	dbPath, err := getDatabasePath()
	if err != nil {
		return nil, err
	}

	database, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return database, nil
}
