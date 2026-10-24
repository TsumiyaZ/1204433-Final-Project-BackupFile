package llama

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	_ "golang.org/x/image/webp"
)

const embeddingDimensions = 768

type SetupProgress struct {
	Message     string
	Percent     int
	Determinate bool
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 2 * time.Minute,
		},
	}
}

func (c *Client) WaitUntilReady(ctx context.Context) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	timeout := time.NewTimer(30 * time.Minute)
	defer timeout.Stop()

	for {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			c.baseURL+"/health",
			nil,
		)
		if err != nil {
			return err
		}

		resp, err := c.httpClient.Do(req)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeout.C:
			return fmt.Errorf("timed out waiting for llama.cpp model")
		case <-ticker.C:
		}
	}
}

type embeddingRequest struct {
	Model          string `json:"model"`
	Input          any    `json:"input"`
	EncodingFormat string `json:"encoding_format"`
}

type embeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
}

type imageURL struct {
	URL string `json:"url"`
}

type imageContent struct {
	Type     string    `json:"type"`
	ImageURL *imageURL `json:"image_url,omitempty"`
}

type imageInput struct {
	Content []imageContent `json:"content"`
}

func StartServer(ctx context.Context, onStatus func(SetupProgress)) error {
	executablePath, err := ResolveServerPath()
	if err != nil {
		if strings.TrimSpace(os.Getenv("LLAMA_SERVER_PATH")) != "" {
			return err
		}
		reportSetupProgress(onStatus, SetupProgress{Message: "ยังไม่พบ llama.cpp กำลังดาวน์โหลดส่วนประกอบที่จำเป็น..."})
		executablePath, err = installServer(ctx, onStatus)
		if err != nil {
			return fmt.Errorf("cannot prepare llama.cpp runtime: %w", err)
		}
	}

	modelPath, mmprojPath, err := ensureEmbeddingModel(ctx, onStatus)
	if err != nil {
		return fmt.Errorf("cannot prepare embedding model: %w", err)
	}

	cmd := exec.CommandContext(
		ctx,
		executablePath,
		"-m", modelPath,
		"--mmproj", mmprojPath,
		"--embeddings",
		"--host", "127.0.0.1",
		"--port", "8080",
	)
	logFile := openServerLog()
	if logFile != nil {
		cmd.Stdout = logFile
		cmd.Stderr = logFile
	}
	configureServerProcess(cmd)

	if err := cmd.Start(); err != nil {
		if logFile != nil {
			logFile.Close()
		}
		return fmt.Errorf("cannot start llama server: %w", err)
	}

	go func() {
		defer func() {
			if logFile != nil {
				logFile.Close()
			}
		}()
		if err := cmd.Wait(); err != nil && ctx.Err() == nil {
			if logFile != nil {
				fmt.Fprintf(logFile, "llama server stopped: %v\n", err)
			}
		}
	}()

	return nil
}

func openServerLog() *os.File {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil
	}
	logDir := filepath.Join(configDir, "PhotoBackup", "logs")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return nil
	}
	logFile, err := os.OpenFile(
		filepath.Join(logDir, "llama-server.log"),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0o600,
	)
	if err != nil {
		return nil
	}
	return logFile
}

// ResolveServerPath finds llama-server using an explicit environment override,
// next to the application executable, or on PATH.
func ResolveServerPath() (string, error) {
	if configuredPath := strings.TrimSpace(os.Getenv("LLAMA_SERVER_PATH")); configuredPath != "" {
		return validateServerPath(configuredPath)
	}

	executableName := "llama-server"
	if runtime.GOOS == "windows" {
		executableName += ".exe"
	}

	if appPath, err := os.Executable(); err == nil {
		appDirectory := filepath.Dir(appPath)
		platformDirectory := runtime.GOOS + "-" + runtime.GOARCH
		candidates := []string{
			filepath.Join(appDirectory, executableName),
			filepath.Join(appDirectory, "llama", platformDirectory, executableName),
			filepath.Join(appDirectory, "llama", executableName),
		}
		if runtime.GOOS == "darwin" {
			// In a macOS .app bundle, keep the OS-specific runtime in Contents/Resources.
			candidates = append(candidates,
				filepath.Join(appDirectory, "..", "Resources", "llama", platformDirectory, executableName),
				filepath.Join(appDirectory, "..", "Resources", "llama", executableName),
			)
		}
		for _, candidate := range candidates {
			if resolved, err := validateServerPath(candidate); err == nil {
				return resolved, nil
			}
		}
	}

	if path, err := exec.LookPath(executableName); err == nil {
		return validateServerPath(path)
	}

	return "", fmt.Errorf(
		"llama-server was not found; place it next to the app or in the llama folder, or set LLAMA_SERVER_PATH",
	)
}

func validateServerPath(path string) (string, error) {
	absolutePath, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return "", fmt.Errorf("invalid llama-server path: %w", err)
	}

	info, err := os.Stat(absolutePath)
	if err != nil {
		return "", fmt.Errorf("cannot access llama-server at %q: %w", absolutePath, err)
	}
	if info.IsDir() {
		return "", fmt.Errorf("llama-server path is a directory: %q", absolutePath)
	}
	return absolutePath, nil
}

func (c *Client) EmbedText(ctx context.Context, text string) ([]float32, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, fmt.Errorf("search text is required")
	}

	return c.requestEmbedding(ctx, embeddingRequest{
		Model:          "embeddinggemma-2",
		Input:          "task: search result | query: " + text,
		EncodingFormat: "float",
	})
}

func (c *Client) EmbedImage(ctx context.Context, imagePath string) ([]float32, error) {
	data, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read image: %w", err)
	}

	contentType := http.DetectContentType(data)
	if !strings.HasPrefix(contentType, "image/") {
		return nil, fmt.Errorf("file content is not a supported image")
	}
	if contentType == "image/webp" {
		data, err = convertWebPToJPEG(data)
		if err != nil {
			return nil, fmt.Errorf("cannot convert WebP image for llama.cpp: %w", err)
		}
		contentType = "image/jpeg"
	}

	dataURL := "data:" + contentType + ";base64," +
		base64.StdEncoding.EncodeToString(data)

	return c.requestEmbedding(ctx, embeddingRequest{
		Model: "embeddinggemma-2",
		Input: []imageInput{
			{
				Content: []imageContent{
					{
						Type:     "image_url",
						ImageURL: &imageURL{URL: dataURL},
					},
				},
			},
		},
		EncodingFormat: "float",
	})
}

func convertWebPToJPEG(data []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	var converted bytes.Buffer
	if err := jpeg.Encode(&converted, img, &jpeg.Options{Quality: 90}); err != nil {
		return nil, err
	}
	return converted.Bytes(), nil
}

func (c *Client) requestEmbedding(
	ctx context.Context,
	payload embeddingRequest,
) ([]float32, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("cannot encode embedding request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/v1/embeddings",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot call llama.cpp: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		responseBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("llama.cpp returned %s: %s", resp.Status, responseBody)
	}

	var result embeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("cannot decode embedding response: %w", err)
	}
	if len(result.Data) == 0 {
		return nil, fmt.Errorf("llama.cpp returned no embedding")
	}

	embedding := result.Data[0].Embedding
	if len(embedding) != embeddingDimensions {
		return nil, fmt.Errorf(
			"expected %d embedding values, got %d",
			embeddingDimensions,
			len(embedding),
		)
	}

	return embedding, nil
}
