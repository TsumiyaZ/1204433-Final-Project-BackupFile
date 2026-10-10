package backup

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// signatureSize is enough bytes to identify JPEG, PNG, GIF and WebP headers.
const signatureSize = 12

// validateImage checks that the extension is supported and that the file
// header matches that extension, so renamed non-image files are rejected.
func validateImage(path string) error {
	extension := strings.ToLower(filepath.Ext(path))

	expected, ok := map[string]string{
		".jpg":  "jpeg",
		".jpeg": "jpeg",
		".png":  "png",
		".gif":  "gif",
		".webp": "webp",
	}[extension]
	if !ok {
		return fmt.Errorf("unsupported file extension %q", extension)
	}

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()

	header := make([]byte, signatureSize)
	n, err := io.ReadFull(file, header)
	if err != nil && err != io.ErrUnexpectedEOF {
		return fmt.Errorf("cannot read file header: %w", err)
	}

	if detectImageType(header[:n]) != expected {
		return fmt.Errorf("file content does not match %s extension", extension)
	}

	return nil
}

func detectImageType(header []byte) string {
	switch {
	case bytes.HasPrefix(header, []byte{0xFF, 0xD8, 0xFF}):
		return "jpeg"
	case bytes.HasPrefix(header, []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A}):
		return "png"
	case bytes.HasPrefix(header, []byte("GIF87a")), bytes.HasPrefix(header, []byte("GIF89a")):
		return "gif"
	case len(header) >= 12 &&
		bytes.Equal(header[0:4], []byte("RIFF")) &&
		bytes.Equal(header[8:12], []byte("WEBP")):
		return "webp"
	default:
		return ""
	}
}

func fileChecksum(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// isInside reports whether path is root itself or located below root.
func isInside(root string, path string) bool {
	relativePath, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}

	return relativePath != ".." &&
		!strings.HasPrefix(relativePath, ".."+string(os.PathSeparator)) &&
		!filepath.IsAbs(relativePath)
}
