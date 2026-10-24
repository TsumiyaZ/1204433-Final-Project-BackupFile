package llama

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const llamaRuntimeVersion = "b11856"

type llamaRuntimeAsset struct {
	url        string
	sha256     string
	format     string
	executable string
}

// Checksums and archive URLs are pinned to an official llama.cpp release so
// first-run setup never executes an unverified binary or silently changes it.
var llamaRuntimeAssets = map[string]llamaRuntimeAsset{
	"windows/amd64": {
		url:        "https://github.com/ggml-org/llama.cpp/releases/download/b11856/llama-b11856-bin-win-cpu-x64.zip",
		sha256:     "b0be608ec81c811029406430796dfcd16b8d83dc152705a6682f0d47834a2f3b",
		format:     "zip",
		executable: "llama-server.exe",
	},
	"windows/arm64": {
		url:        "https://github.com/ggml-org/llama.cpp/releases/download/b11856/llama-b11856-bin-win-cpu-arm64.zip",
		sha256:     "24a8750d9116d66843dc293c71c66564461e916f6b2d3204045576312175ebe9",
		format:     "zip",
		executable: "llama-server.exe",
	},
	"darwin/arm64": {
		url:        "https://github.com/ggml-org/llama.cpp/releases/download/b11856/llama-b11856-bin-macos-arm64.tar.gz",
		sha256:     "f961fd5a488a22177c4942765b7ba8371f7025b53320656d46c06aed6ca43738",
		format:     "tar.gz",
		executable: "llama-server",
	},
	"darwin/amd64": {
		url:        "https://github.com/ggml-org/llama.cpp/releases/download/b11856/llama-b11856-bin-macos-x64.tar.gz",
		sha256:     "a06974186c508a8cb229803b5561fc6f0540702197fdd2dafc243fcee1d44708",
		format:     "tar.gz",
		executable: "llama-server",
	},
	"linux/amd64": {
		url:        "https://github.com/ggml-org/llama.cpp/releases/download/b11856/llama-b11856-bin-ubuntu-x64.tar.gz",
		sha256:     "71b660dd4a90c13db899acdef59862226f7ec8d27907b7c107febc99904c00be",
		format:     "tar.gz",
		executable: "llama-server",
	},
	"linux/arm64": {
		url:        "https://github.com/ggml-org/llama.cpp/releases/download/b11856/llama-b11856-bin-ubuntu-arm64.tar.gz",
		sha256:     "540871d9a0aac2b9b92ff42dd6f14449e418a7638302e7c11ccaa9f8b9dc6c98",
		format:     "tar.gz",
		executable: "llama-server",
	},
}

func installServer(ctx context.Context, onStatus func(message string)) (string, error) {
	asset, ok := llamaRuntimeAssets[runtime.GOOS+"/"+runtime.GOARCH]
	if !ok {
		return "", fmt.Errorf("automatic llama.cpp setup is not supported on %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("cannot locate user cache directory: %w", err)
	}
	platform := runtime.GOOS + "-" + runtime.GOARCH
	installDir := filepath.Join(cacheDir, "PhotoBackup", "llama.cpp", llamaRuntimeVersion, platform+"-"+asset.sha256[:8])
	if executablePath, err := findRuntimeExecutable(installDir, asset.executable); err == nil {
		return executablePath, nil
	}

	parentDir := filepath.Dir(installDir)
	if err := os.MkdirAll(parentDir, 0o755); err != nil {
		return "", fmt.Errorf("cannot create runtime folder: %w", err)
	}
	tempDir, err := os.MkdirTemp(parentDir, ".llama-install-*")
	if err != nil {
		return "", fmt.Errorf("cannot create temporary runtime folder: %w", err)
	}
	defer os.RemoveAll(tempDir)

	archivePath := filepath.Join(tempDir, "llama-runtime."+strings.ReplaceAll(asset.format, ".", ""))
	if err := downloadRuntime(ctx, asset, archivePath, onStatus); err != nil {
		return "", err
	}

	extractDir := filepath.Join(tempDir, "extracted")
	if err := os.Mkdir(extractDir, 0o755); err != nil {
		return "", err
	}
	if onStatus != nil {
		onStatus("ดาวน์โหลดส่วนประกอบเสร็จแล้ว กำลังติดตั้ง...")
	}
	if err := extractRuntimeArchive(archivePath, asset.format, extractDir); err != nil {
		return "", fmt.Errorf("cannot unpack llama.cpp runtime: %w", err)
	}
	if _, err := findRuntimeExecutable(extractDir, asset.executable); err != nil {
		return "", err
	}

	if err := os.Rename(extractDir, installDir); err != nil {
		if executablePath, existingErr := findRuntimeExecutable(installDir, asset.executable); existingErr == nil {
			return executablePath, nil
		}
		return "", fmt.Errorf("cannot install llama.cpp runtime: %w", err)
	}
	return findRuntimeExecutable(installDir, asset.executable)
}

func downloadRuntime(
	ctx context.Context,
	asset llamaRuntimeAsset,
	destination string,
	onStatus func(message string),
) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, asset.url, nil)
	if err != nil {
		return fmt.Errorf("cannot create runtime download request: %w", err)
	}
	request.Header.Set("User-Agent", "PhotoBackup")
	client := &http.Client{Timeout: 15 * time.Minute}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("cannot download llama.cpp runtime: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("llama.cpp download returned HTTP %s", response.Status)
	}

	file, err := os.OpenFile(destination, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("cannot create runtime download file: %w", err)
	}
	defer file.Close()

	hasher := sha256.New()
	writer := &runtimeDownloadWriter{
		writer: file,
		hasher: hasher,
		total:  response.ContentLength,
		onStatus: func(percent int) {
			if onStatus != nil {
				onStatus(fmt.Sprintf("กำลังดาวน์โหลด llama.cpp %d%%", percent))
			}
		},
	}
	if _, err := io.Copy(writer, response.Body); err != nil {
		return fmt.Errorf("cannot save llama.cpp runtime: %w", err)
	}
	if got := hex.EncodeToString(hasher.Sum(nil)); got != asset.sha256 {
		return fmt.Errorf("llama.cpp runtime checksum mismatch (got %s)", got)
	}
	return nil
}

type runtimeDownloadWriter struct {
	writer      io.Writer
	hasher      io.Writer
	total       int64
	downloaded  int64
	lastPercent int
	onStatus    func(percent int)
}

func (w *runtimeDownloadWriter) Write(data []byte) (int, error) {
	n, err := w.writer.Write(data)
	if n > 0 {
		_, hashErr := w.hasher.Write(data[:n])
		if hashErr != nil {
			return n, hashErr
		}
		w.downloaded += int64(n)
		if w.total > 0 && w.onStatus != nil {
			percent := int(w.downloaded * 100 / w.total)
			if percent != w.lastPercent {
				w.lastPercent = percent
				w.onStatus(percent)
			}
		}
	}
	return n, err
}

func extractRuntimeArchive(archivePath, format, destination string) error {
	if format == "zip" {
		return extractRuntimeZip(archivePath, destination)
	}
	return extractRuntimeTarGz(archivePath, destination)
}

func extractRuntimeZip(archivePath, destination string) error {
	archive, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer archive.Close()

	for _, entry := range archive.File {
		target, err := safeRuntimeArchivePath(destination, entry.Name)
		if err != nil {
			return err
		}
		if entry.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
			continue
		}
		if err := ensureNoRuntimeSymlinkParents(destination, target); err != nil {
			return err
		}
		if entry.Mode()&fs.ModeSymlink != 0 {
			reader, err := entry.Open()
			if err != nil {
				return err
			}
			linkTarget, readErr := io.ReadAll(io.LimitReader(reader, 4096))
			reader.Close()
			if readErr != nil {
				return readErr
			}
			if err := createSafeRuntimeSymlink(destination, target, string(linkTarget)); err != nil {
				return err
			}
			continue
		}

		reader, err := entry.Open()
		if err != nil {
			return err
		}
		if err := writeRuntimeFile(target, reader, entry.Mode().Perm()); err != nil {
			reader.Close()
			return err
		}
		reader.Close()
	}
	return nil
}

func extractRuntimeTarGz(archivePath, destination string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()
	compressed, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer compressed.Close()
	reader := tar.NewReader(compressed)

	for {
		header, err := reader.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		target, err := safeRuntimeArchivePath(destination, header.Name)
		if err != nil {
			return err
		}
		if err := ensureNoRuntimeSymlinkParents(destination, target); err != nil {
			return err
		}
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
		case tar.TypeReg, tar.TypeRegA:
			if err := writeRuntimeFile(target, reader, fs.FileMode(header.Mode).Perm()); err != nil {
				return err
			}
		case tar.TypeSymlink:
			if err := createSafeRuntimeSymlink(destination, target, header.Linkname); err != nil {
				return err
			}
		case tar.TypeLink:
			linkTarget, err := safeRuntimeArchivePath(destination, header.Linkname)
			if err != nil {
				return err
			}
			if err := ensureNoRuntimeSymlinkParents(destination, linkTarget); err != nil {
				return err
			}
			if err := os.Link(linkTarget, target); err != nil {
				return err
			}
		}
	}
}

func writeRuntimeFile(destination string, source io.Reader, mode fs.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return err
	}
	file, err := os.OpenFile(destination, os.O_CREATE|os.O_EXCL|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(file, source)
	closeErr := file.Close()
	if copyErr != nil {
		return copyErr
	}
	if closeErr != nil {
		return closeErr
	}
	if runtime.GOOS != "windows" {
		return os.Chmod(destination, mode)
	}
	return nil
}

func safeRuntimeArchivePath(root, name string) (string, error) {
	if filepath.IsAbs(name) || filepath.VolumeName(name) != "" {
		return "", fmt.Errorf("archive contains an absolute path")
	}
	target := filepath.Join(root, filepath.Clean(filepath.FromSlash(name)))
	relative, err := filepath.Rel(root, target)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("archive path escapes installation directory")
	}
	return target, nil
}

func ensureNoRuntimeSymlinkParents(root, target string) error {
	current := filepath.Dir(target)
	for {
		if current == root || current == filepath.Dir(current) {
			return nil
		}
		info, err := os.Lstat(current)
		if err == nil && info.Mode()&fs.ModeSymlink != 0 {
			return fmt.Errorf("archive attempts to write through a symbolic link")
		}
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		current = filepath.Dir(current)
	}
}

func createSafeRuntimeSymlink(root, target string, linkName string) error {
	if filepath.IsAbs(linkName) || filepath.VolumeName(linkName) != "" {
		return fmt.Errorf("archive contains an absolute symbolic link")
	}
	linkTarget := filepath.Clean(filepath.Join(filepath.Dir(target), filepath.FromSlash(linkName)))
	if _, err := safeRuntimeArchivePath(root, mustRelativePath(root, linkTarget)); err != nil {
		return err
	}
	if err := ensureNoRuntimeSymlinkParents(root, target); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}
	return os.Symlink(linkName, target)
}

func mustRelativePath(root, target string) string {
	relative, err := filepath.Rel(root, target)
	if err != nil {
		return ".."
	}
	return relative
}

func findRuntimeExecutable(root, executable string) (string, error) {
	var found string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || entry.Name() != executable {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if info.Mode().IsRegular() {
			found = path
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if found == "" {
		return "", fmt.Errorf("llama-server executable was not found in runtime folder")
	}
	return filepath.Abs(found)
}

func (w *runtimeDownloadWriter) Close() error {
	if closer, ok := w.writer.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

var _ io.Writer = (*runtimeDownloadWriter)(nil)
