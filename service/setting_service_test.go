package service

import (
	"context"
	"path/filepath"
	"testing"

	"FinalProject/model"
)

type fakeSettingRepository struct {
	setting *model.Setting
}

func (f *fakeSettingRepository) GetSetting(
	ctx context.Context,
) (*model.Setting, error) {
	return f.setting, nil
}

func (f *fakeSettingRepository) SaveSetting(
	ctx context.Context,
	setting *model.Setting,
) error {
	f.setting = setting
	return nil
}

func TestSettingServiceSaveSetting(t *testing.T) {
	repo := &fakeSettingRepository{}
	service := NewSettingService(repo)

	setting := &model.Setting{
		Source: " D:/Source/ ",
		Dest:   " D:/Backup/ ",
	}

	err := service.SaveSetting(context.Background(), setting)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedSource := filepath.Clean("D:/Source/")
	expectedDest := filepath.Clean("D:/Backup/")

	if repo.setting.Source != expectedSource {
		t.Errorf(
			"expected source %q, got %q",
			expectedSource,
			repo.setting.Source,
		)
	}

	if repo.setting.Dest != expectedDest {
		t.Errorf(
			"expected destination %q, got %q",
			expectedDest,
			repo.setting.Dest,
		)
	}
}

func TestSettingServiceValidation(t *testing.T) {
	tests := []struct {
		name    string
		setting *model.Setting
	}{
		{
			name:    "nil setting",
			setting: nil,
		},
		{
			name: "empty source",
			setting: &model.Setting{
				Source: "",
				Dest:   "D:/Backup",
			},
		},
		{
			name: "empty destination",
			setting: &model.Setting{
				Source: "D:/Source",
				Dest:   "",
			},
		},
		{
			name: "same source and destination",
			setting: &model.Setting{
				Source: "D:/Photos",
				Dest:   "d:/photos",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &fakeSettingRepository{}
			service := NewSettingService(repo)

			err := service.SaveSetting(
				context.Background(),
				test.setting,
			)

			if err == nil {
				t.Error("expected validation error, got nil")
			}
		})
	}
}
