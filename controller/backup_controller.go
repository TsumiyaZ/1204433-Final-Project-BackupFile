package controller

import (
	"context"

	"FinalProject/dto"
	"FinalProject/service/backup"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	EventBackupProgress  = "backup:progress"
	EventBackupCompleted = "backup:completed"
	EventBackupError     = "backup:error"
)

type BackupController struct {
	ctx     context.Context
	service backup.BackupService
}

func NewBackupController(service backup.BackupService) *BackupController {
	return &BackupController{
		service: service,
	}
}

// SetBackupContext stores the Wails context so the controller can emit events.
// It is a plain function, not a method, so Wails does not expose it to the frontend.
func SetBackupContext(c *BackupController, ctx context.Context) {
	c.ctx = ctx
}

func (c *BackupController) MovePhotos(
	request dto.MovePhotosRequest,
) (*dto.MovePhotosResult, error) {
	result, err := c.service.MovePhotos(c.ctx, request, func(progress dto.MoveProgress) {
		runtime.EventsEmit(c.ctx, EventBackupProgress, progress)
	})

	if err != nil {
		runtime.EventsEmit(c.ctx, EventBackupError, dto.BackupError{Message: err.Error()})
		return nil, err
	}

	runtime.EventsEmit(c.ctx, EventBackupCompleted, result)

	return result, nil
}
