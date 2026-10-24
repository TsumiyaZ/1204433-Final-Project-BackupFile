package repository

import (
	"io"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func DatabasePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	projectDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return databasePathFor(configDir, projectDir)
}

func databasePathFor(configDir string, projectDir string) (string, error) {
	databaseDir := filepath.Join(configDir, "PhotoBackup", "database")

	if err := os.MkdirAll(databaseDir, 0755); err != nil {
		return "", err
	}

	dbPath := filepath.Join(databaseDir, "database.db")
	legacyPath := filepath.Join(projectDir, "database", "database.db")

	if _, err := os.Stat(dbPath); err == nil {
		return dbPath, nil
	} else if !os.IsNotExist(err) {
		return "", err
	}

	// Keep the existing project-local database available after changing the
	// app's storage location. The original file is intentionally left intact.
	if _, err := os.Stat(legacyPath); err == nil {
		if err := copyDatabase(legacyPath, dbPath); err != nil {
			return "", err
		}
	} else if !os.IsNotExist(err) {
		return "", err
	}

	return dbPath, nil
}

func copyDatabase(sourcePath string, destinationPath string) error {
	source, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer source.Close()

	temporaryFile, err := os.CreateTemp(filepath.Dir(destinationPath), "database-*.tmp")
	if err != nil {
		return err
	}
	temporaryPath := temporaryFile.Name()
	defer os.Remove(temporaryPath)

	if _, err := io.Copy(temporaryFile, source); err != nil {
		temporaryFile.Close()
		return err
	}
	if err := temporaryFile.Close(); err != nil {
		return err
	}

	return os.Rename(temporaryPath, destinationPath)
}

func NewDBConnection() (*gorm.DB, error) {
	dbPath, err := DatabasePath()
	if err != nil {
		return nil, err
	}

	database, err := gorm.Open(
		sqlite.Open(dbPath+"?_pragma=foreign_keys(1)"),
		&gorm.Config{},
	)
	if err != nil {
		return nil, err
	}

	return database, nil
}
