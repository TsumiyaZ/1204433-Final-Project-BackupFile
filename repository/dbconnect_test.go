package repository

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDatabasePathMigratesLegacyDatabaseWithoutRemovingIt(t *testing.T) {
	configDir := t.TempDir()
	projectDir := t.TempDir()
	legacyDir := filepath.Join(projectDir, "database")
	if err := os.MkdirAll(legacyDir, 0755); err != nil {
		t.Fatal(err)
	}

	legacyPath := filepath.Join(legacyDir, "database.db")
	legacyContents := []byte("existing database")
	if err := os.WriteFile(legacyPath, legacyContents, 0600); err != nil {
		t.Fatal(err)
	}

	got, err := databasePathFor(configDir, projectDir)
	if err != nil {
		t.Fatal(err)
	}

	want := filepath.Join(configDir, "PhotoBackup", "database", "database.db")
	if got != want {
		t.Fatalf("databasePathFor() = %q, want %q", got, want)
	}

	newContents, err := os.ReadFile(got)
	if err != nil {
		t.Fatal(err)
	}
	if string(newContents) != string(legacyContents) {
		t.Fatalf("copied database = %q, want %q", newContents, legacyContents)
	}

	if _, err := os.Stat(legacyPath); err != nil {
		t.Fatalf("legacy database should be preserved: %v", err)
	}
}

func TestDatabasePathKeepsExistingUserDatabase(t *testing.T) {
	configDir := t.TempDir()
	projectDir := t.TempDir()
	userDB := filepath.Join(configDir, "PhotoBackup", "database", "database.db")
	if err := os.MkdirAll(filepath.Dir(userDB), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(userDB, []byte("user database"), 0600); err != nil {
		t.Fatal(err)
	}

	legacyDir := filepath.Join(projectDir, "database")
	if err := os.MkdirAll(legacyDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(legacyDir, "database.db"), []byte("legacy database"), 0600); err != nil {
		t.Fatal(err)
	}

	got, err := databasePathFor(configDir, projectDir)
	if err != nil {
		t.Fatal(err)
	}
	if got != userDB {
		t.Fatalf("databasePathFor() = %q, want existing user database %q", got, userDB)
	}

	contents, err := os.ReadFile(got)
	if err != nil {
		t.Fatal(err)
	}
	if string(contents) != "user database" {
		t.Fatalf("existing user database was changed: %q", contents)
	}
}
