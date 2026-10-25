package llama

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type CaptionClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewCaptionClient(baseURL string) *CaptionClient {
	return &CaptionClient{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{Timeout: 5 * time.Minute},
	}
}

func (c *CaptionClient) WaitUntilReady(ctx context.Context) error {
	return (&Client{baseURL: c.baseURL, httpClient: c.httpClient}).WaitUntilReady(ctx)
}

func (c *CaptionClient) DescribeImage(ctx context.Context, imagePath string) (string, error) {
	data, err := os.ReadFile(imagePath)
	if err != nil {
		return "", fmt.Errorf("cannot read image for description: %w", err)
	}
	contentType := http.DetectContentType(data)
	if !strings.HasPrefix(contentType, "image/") {
		return "", fmt.Errorf("file content is not an image")
	}
	if contentType == "image/webp" {
		data, err = convertWebPToJPEG(data)
		if err != nil {
			return "", fmt.Errorf("cannot convert image for description: %w", err)
		}
		contentType = "image/jpeg"
	}
	imageDataURL := "data:" + contentType + ";base64," + base64.StdEncoding.EncodeToString(data)

	payload := map[string]any{
		"model": "qwen2.5-vl-3b",
		"messages": []any{map[string]any{
			"role": "user",
			"content": []any{
				map[string]any{"type": "text", "text": "อธิบายภาพนี้เป็นภาษาไทย 1-2 ประโยค บอกเฉพาะสิ่งที่มองเห็นชัดเจน ไม่เดาข้อมูลที่ภาพยืนยันไม่ได้ ตอบเป็นคำบรรยายโดยตรง ไม่ต้องใส่หัวข้อหรือรายการแท็ก"},
				map[string]any{"type": "image_url", "image_url": map[string]any{"url": imageDataURL}},
			},
		}},
		"temperature": 0.2,
		"max_tokens":  120,
		"stream":      false,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("cannot encode image description request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("cannot request image description: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("image description request returned HTTP %s", resp.Status)
	}
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("cannot decode image description: %w", err)
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("image description response has no choices")
	}
	description := strings.TrimSpace(result.Choices[0].Message.Content)
	if description == "" {
		return "", fmt.Errorf("image description response is empty")
	}
	if len([]rune(description)) > 800 {
		return "", fmt.Errorf("image description is longer than expected")
	}
	return description, nil
}

func StartCaptionServer(ctx context.Context, onStatus func(SetupProgress)) error {
	executablePath, err := resolveOrInstallServer(ctx, onStatus)
	if err != nil {
		return err
	}
	modelPath, mmprojPath, err := ensureCaptionModel(ctx, onStatus)
	if err != nil {
		return fmt.Errorf("cannot prepare image description model: %w", err)
	}
	cmd := exec.CommandContext(ctx, executablePath,
		"-m", modelPath,
		"--mmproj", mmprojPath,
		"--host", "127.0.0.1",
		"--port", "8081",
		"--ctx-size", "2048",
		"--parallel", "1",
	)
	logFile := openCaptionServerLog()
	if logFile != nil {
		cmd.Stdout = logFile
		cmd.Stderr = logFile
	}
	configureServerProcess(cmd)
	if err := cmd.Start(); err != nil {
		if logFile != nil {
			logFile.Close()
		}
		return fmt.Errorf("cannot start image description model: %w", err)
	}
	go func() {
		defer func() {
			if logFile != nil {
				logFile.Close()
			}
		}()
		if err := cmd.Wait(); err != nil && ctx.Err() == nil && logFile != nil {
			fmt.Fprintf(logFile, "caption server stopped: %v\n", err)
		}
	}()
	return nil
}

func openCaptionServerLog() *os.File {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil
	}
	logDir := filepath.Join(configDir, "PhotoBackup", "logs")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return nil
	}
	logFile, err := os.OpenFile(filepath.Join(logDir, "llama-caption-server.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return nil
	}
	return logFile
}
