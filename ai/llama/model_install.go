package llama

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const embeddingModelRevision = "bfcd298762cc34d0357ece5ebdd31791a3a374d8"

type modelAsset struct {
	name   string
	size   int64
	sha256 string
	url    string
}

var embeddingModelAssets = []modelAsset{
	{
		name:   "embeddinggemma-2-Q8_0.gguf",
		size:   309855456,
		sha256: "2188ac1deca4b77dffefd603c2776a9d76d9d74ec01841392982ebb840b09135",
		url:    "https://huggingface.co/ggml-org/embeddinggemma-2-GGUF/resolve/bfcd298762cc34d0357ece5ebdd31791a3a374d8/embeddinggemma-2-Q8_0.gguf",
	},
	{
		name:   "mmproj-embeddinggemma-2-Q8_0.gguf",
		size:   554821024,
		sha256: "c4a8a52691ecef40618438928bdf9e68379b854e24166f292592353db0aab64f",
		url:    "https://huggingface.co/ggml-org/embeddinggemma-2-GGUF/resolve/bfcd298762cc34d0357ece5ebdd31791a3a374d8/mmproj-embeddinggemma-2-Q8_0.gguf",
	},
}

func ensureEmbeddingModel(ctx context.Context, onStatus func(SetupProgress)) (string, string, error) {
	modelDir, err := ModelDirectory()
	if err != nil {
		return "", "", fmt.Errorf("cannot locate local application data: %w", err)
	}
	if err := os.MkdirAll(modelDir, 0o755); err != nil {
		return "", "", fmt.Errorf("cannot create model folder: %w", err)
	}

	paths := make([]string, len(embeddingModelAssets))
	var totalSize int64
	for _, asset := range embeddingModelAssets {
		totalSize += asset.size
	}
	var completedSize int64
	for i, asset := range embeddingModelAssets {
		path := filepath.Join(modelDir, asset.name)
		paths[i] = path
		if validModelFile(path, asset.sha256) {
			completedSize += asset.size
			continue
		}
		assetProgress := func(progress SetupProgress) {
			if progress.Determinate && totalSize > 0 {
				currentSize := asset.size * int64(progress.Percent) / 100
				progress.Percent = int((completedSize + currentSize) * 100 / totalSize)
			}
			reportSetupProgress(onStatus, progress)
		}
		if err := ensureModelFile(ctx, path, asset, assetProgress); err != nil {
			return "", "", err
		}
		completedSize += asset.size
	}
	return paths[0], paths[1], nil
}

func ModelDirectory() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("cannot locate local application data: %w", err)
	}
	return filepath.Join(cacheDir, "PhotoBackup", "Models"), nil
}

func ensureModelFile(ctx context.Context, destination string, asset modelAsset, onStatus func(SetupProgress)) error {
	if validModelFile(destination, asset.sha256) {
		return nil
	}
	if err := os.Remove(destination); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("cannot replace invalid model file %s: %w", asset.name, err)
	}

	partialPath := destination + ".download"
	if validModelFile(partialPath, asset.sha256) {
		if err := os.Rename(partialPath, destination); err != nil {
			return fmt.Errorf("cannot install downloaded model file %s: %w", asset.name, err)
		}
		return nil
	}
	if err := downloadModelFile(ctx, partialPath, asset, onStatus); err != nil {
		return err
	}
	if err := os.Rename(partialPath, destination); err != nil {
		return fmt.Errorf("cannot install model file %s: %w", asset.name, err)
	}
	return nil
}

func validModelFile(path, expectedSHA256 string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return false
	}
	return hex.EncodeToString(hasher.Sum(nil)) == expectedSHA256
}

func downloadModelFile(ctx context.Context, partialPath string, asset modelAsset, onStatus func(SetupProgress)) error {
	url := asset.url
	client := &http.Client{Timeout: 45 * time.Minute}

	var offset int64
	if info, err := os.Stat(partialPath); err == nil {
		offset = info.Size()
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("cannot inspect partial model file %s: %w", asset.name, err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("cannot create model download request: %w", err)
	}
	request.Header.Set("User-Agent", "PhotoBackup")
	if offset > 0 {
		request.Header.Set("Range", fmt.Sprintf("bytes=%d-", offset))
	}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("cannot download model file %s: %w", asset.name, err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusPartialContent {
		if response.StatusCode == http.StatusRequestedRangeNotSatisfiable && offset > 0 {
			response.Body.Close()
			if err := os.Remove(partialPath); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("cannot reset partial model download %s: %w", asset.name, err)
			}
			return downloadModelFile(ctx, partialPath, asset, onStatus)
		}
		return fmt.Errorf("model download for %s returned HTTP %s", asset.name, response.Status)
	}
	modelName := asset.name[:len(asset.name)-len(filepath.Ext(asset.name))]

	// A server may ignore Range. In that case restart the partial file rather than
	// appending a second full copy to it.
	appendMode := offset > 0 && response.StatusCode == http.StatusPartialContent
	if !appendMode {
		offset = 0
	}
	flags := os.O_CREATE | os.O_WRONLY
	if appendMode {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}
	file, err := os.OpenFile(partialPath, flags, 0o600)
	if err != nil {
		return fmt.Errorf("cannot create model file %s: %w", asset.name, err)
	}

	total := response.ContentLength
	if response.StatusCode == http.StatusPartialContent && response.ContentLength >= 0 {
		total += offset
	}
	percent := 0
	if total > 0 {
		percent = int(offset * 100 / total)
	}
	reportSetupProgress(onStatus, SetupProgress{
		Message:     "กำลังดาวน์โหลดโมเดล " + modelName,
		Percent:     percent,
		Determinate: total > 0,
	})
	writer := &modelDownloadWriter{
		file: file, name: modelName, downloaded: offset, total: total, onStatus: onStatus,
	}
	_, copyErr := io.Copy(writer, response.Body)
	closeErr := file.Close()
	if copyErr != nil {
		return fmt.Errorf("cannot save model file %s: %w", asset.name, copyErr)
	}
	if closeErr != nil {
		return fmt.Errorf("cannot close model file %s: %w", asset.name, closeErr)
	}
	if !validModelFile(partialPath, asset.sha256) {
		return fmt.Errorf("model file %s failed SHA-256 verification; partial download was kept for retry", asset.name)
	}
	return nil
}

type modelDownloadWriter struct {
	file       io.Writer
	name       string
	downloaded int64
	total      int64
	lastUpdate time.Time
	onStatus   func(SetupProgress)
}

func (w *modelDownloadWriter) Write(data []byte) (int, error) {
	n, err := w.file.Write(data)
	if n == 0 || w.onStatus == nil {
		return n, err
	}
	w.downloaded += int64(n)
	if time.Since(w.lastUpdate) < 250*time.Millisecond && (w.total <= 0 || w.downloaded < w.total) {
		return n, err
	}
	w.lastUpdate = time.Now()
	if w.total > 0 {
		percent := int(w.downloaded * 100 / w.total)
		w.onStatus(SetupProgress{Message: "กำลังดาวน์โหลดโมเดล " + w.name, Percent: percent, Determinate: true})
	} else {
		w.onStatus(SetupProgress{Message: fmt.Sprintf("กำลังดาวน์โหลดโมเดล %s (%d MB)", w.name, w.downloaded/(1024*1024))})
	}
	return n, err
}

func reportSetupProgress(onStatus func(SetupProgress), progress SetupProgress) {
	if onStatus != nil {
		onStatus(progress)
	}
}
