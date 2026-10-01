package service

import (
	"FinalProject/model"
	"FinalProject/repository"
	"context"
	"errors"
	"path/filepath"
	"strings"
)

type settingService struct {
	repo repository.SettingRepository
}

type SettingService interface {
	GetSetting(ctx context.Context) (*model.Setting, error)
	SaveSetting(ctx context.Context, setting *model.Setting) error
}

func NewSettingService(repo repository.SettingRepository) SettingService {
	return &settingService{
		repo: repo,
	}
}

func (s *settingService) GetSetting(ctx context.Context) (*model.Setting, error) {
	return s.repo.GetSetting(ctx)
}

func (s *settingService) SaveSetting(ctx context.Context, setting *model.Setting) error {
	if setting == nil {
		return errors.New("setting is required")
	}

	source := strings.TrimSpace(setting.Source)
	dest := strings.TrimSpace(setting.Dest)

	if source == "" {
		return errors.New("source is required")
	}

	if dest == "" {
		return errors.New("destination is required")
	}

	source = filepath.Clean(source)
	dest = filepath.Clean(dest)

	if strings.EqualFold(source, dest) {
		return errors.New("source and destination must be different")
	}

	setting.Source = source
	setting.Dest = dest

	return s.repo.SaveSetting(ctx, setting)
}
