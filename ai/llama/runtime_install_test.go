package llama

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestSafeRuntimeArchivePath(t *testing.T) {
	root := t.TempDir()

	path, err := safeRuntimeArchivePath(root, "bin/llama-server")
	if err != nil {
		t.Fatalf("valid archive path returned an error: %v", err)
	}
	if want := filepath.Join(root, "bin", "llama-server"); path != want {
		t.Fatalf("path = %q, want %q", path, want)
	}

	for _, name := range []string{"../outside", "bin/../../outside"} {
		if _, err := safeRuntimeArchivePath(root, name); err == nil {
			t.Errorf("expected archive path %q to be rejected", name)
		}
	}

}

func TestInstallRuntimeSmoke(t *testing.T) {
	if os.Getenv("PHOTO_BACKUP_TEST_RUNTIME_INSTALL") != "1" {
		t.Skip("set PHOTO_BACKUP_TEST_RUNTIME_INSTALL=1 to download and verify the pinned runtime")
	}

	executable, err := installServer(context.Background(), func(message string) {
		t.Log(message)
	})
	if err != nil {
		t.Fatal(err)
	}

	output, err := exec.Command(executable, "--version").CombinedOutput()
	if err != nil {
		t.Fatalf("llama-server --version failed: %v\n%s", err, output)
	}
	t.Logf("installed runtime: %s", strings.TrimSpace(string(output)))
}

func TestPinnedRuntimeAssetsHaveChecksums(t *testing.T) {
	for platform, asset := range llamaRuntimeAssets {
		if !strings.HasPrefix(asset.url, "https://github.com/ggml-org/llama.cpp/releases/download/") {
			t.Errorf("%s URL is not an official HTTPS release: %q", platform, asset.url)
		}
		if len(asset.sha256) != 64 {
			t.Errorf("%s has invalid SHA-256 length %d", platform, len(asset.sha256))
		}
		if asset.format != "zip" && asset.format != "tar.gz" {
			t.Errorf("%s has unsupported archive format %q", platform, asset.format)
		}
	}
}
