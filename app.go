package main

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

func (a *App) SelectSourceDirectory() (string, error) {
	selectedFile, err := runtime.OpenFileDialog(
		a.ctx,
		runtime.OpenDialogOptions{
			Title: "เลือกรูปภาพจากโฟลเดอร์ต้นทาง",
			Filters: []runtime.FileFilter{
				{
					DisplayName: "Image Files",
					Pattern:     "*.jpg;*.jpeg;*.png;*.webp;*.gif",
				},
			},
		},
	)

	if err != nil || selectedFile == "" {
		return "", err
	}

	return filepath.Dir(selectedFile), nil
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
	return &App{}
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
