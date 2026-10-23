package pathutil

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Same reports whether two paths identify the same location. It first checks
// the filesystem so case-sensitive and case-insensitive volumes both behave
// correctly, then applies Windows path comparison rules as a fallback.
func Same(first string, second string) bool {
	firstPath, firstErr := filepath.Abs(filepath.Clean(first))
	secondPath, secondErr := filepath.Abs(filepath.Clean(second))
	if firstErr != nil || secondErr != nil {
		return false
	}
	if firstPath == secondPath {
		return true
	}

	firstInfo, firstErr := os.Stat(firstPath)
	secondInfo, secondErr := os.Stat(secondPath)
	if firstErr == nil && secondErr == nil && os.SameFile(firstInfo, secondInfo) {
		return true
	}

	return runtime.GOOS == "windows" && strings.EqualFold(firstPath, secondPath)
}

// Key normalizes a path for duplicate tracking using the platform's common
// path comparison rules.
func Key(path string) string {
	cleanPath := filepath.Clean(path)
	if runtime.GOOS == "windows" {
		return strings.ToLower(cleanPath)
	}
	return cleanPath
}
