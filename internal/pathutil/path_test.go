package pathutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSameReturnsTrueForSameDirectory(t *testing.T) {
	directory := t.TempDir()
	if !Same(directory, filepath.Join(directory, ".")) {
		t.Fatal("Same() should recognize identical directory paths")
	}
}

func TestSameReturnsFalseForDifferentDirectories(t *testing.T) {
	first := t.TempDir()
	second := t.TempDir()
	if Same(first, second) {
		t.Fatal("Same() should distinguish different directories")
	}
}

func TestSameRecognizesSameFilesystemObject(t *testing.T) {
	directory := t.TempDir()
	alias := filepath.Join(t.TempDir(), "alias")
	if err := os.Symlink(directory, alias); err != nil {
		t.Skipf("symlinks are unavailable: %v", err)
	}
	if !Same(directory, alias) {
		t.Fatal("Same() should recognize paths that resolve to the same directory")
	}
}
