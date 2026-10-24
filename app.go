package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context

	embeddingStatusMu sync.RWMutex
	embeddingStatus   EmbeddingStatus
}

type EmbeddingStatus struct {
	State           string `json:"state"`
	Message         string `json:"message"`
	Progress        int    `json:"progress"`
	ProgressVisible bool   `json:"progressVisible"`
}

func (a *App) SelectSourceDirectory() (string, error) {
	return runtime.OpenDirectoryDialog(
		a.ctx,
		runtime.OpenDialogOptions{
			Title: "เลือกโฟลเดอร์ต้นทาง",
		},
	)
}

func (a *App) SelectDestinationDirectory() (string, error) {
	return runtime.OpenDirectoryDialog(
		a.ctx,
		runtime.OpenDialogOptions{
			Title: "เลือกโฟลเดอร์ปลายทาง",
		},
	)
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		embeddingStatus: EmbeddingStatus{
			State:   "loading",
			Message: "กำลังเริ่มเตรียมโมเดลค้นหารูป...",
		},
	}
}

func (a *App) GetEmbeddingStatus() EmbeddingStatus {
	a.embeddingStatusMu.RLock()
	defer a.embeddingStatusMu.RUnlock()
	return a.embeddingStatus
}

func (a *App) setEmbeddingStatus(state string, message string) {
	a.embeddingStatusMu.Lock()
	defer a.embeddingStatusMu.Unlock()
	a.embeddingStatus = EmbeddingStatus{
		State:           state,
		Message:         message,
		ProgressVisible: false,
	}
}

func (a *App) setEmbeddingProgress(message string, progress int, visible bool) {
	if progress < 0 {
		progress = 0
	}
	if progress > 100 {
		progress = 100
	}
	a.embeddingStatusMu.Lock()
	defer a.embeddingStatusMu.Unlock()
	a.embeddingStatus = EmbeddingStatus{
		State:           "loading",
		Message:         message,
		Progress:        progress,
		ProgressVisible: visible,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
