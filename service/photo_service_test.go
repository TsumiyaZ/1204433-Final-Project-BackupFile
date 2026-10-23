package service

import "testing"

func TestImageExtensionSupport(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		scannable   bool
		previewable bool
	}{
		{name: "JPEG is scanned and previewed", path: "photo.jpg", scannable: true, previewable: true},
		{name: "uppercase JPEG is scanned and previewed", path: "photo.JPEG", scannable: true, previewable: true},
		{name: "PNG is scanned and previewed", path: "photo.png", scannable: true, previewable: true},
		{name: "WebP is previewed but not scanned", path: "photo.webp", scannable: false, previewable: true},
		{name: "GIF is previewed but not scanned", path: "photo.gif", scannable: false, previewable: true},
		{name: "AVIF is unsupported", path: "photo.avif", scannable: false, previewable: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isScannableImage(tt.path); got != tt.scannable {
				t.Errorf("isScannableImage(%q) = %t, want %t", tt.path, got, tt.scannable)
			}
			if got := isPreviewableImage(tt.path); got != tt.previewable {
				t.Errorf("isPreviewableImage(%q) = %t, want %t", tt.path, got, tt.previewable)
			}
		})
	}
}
