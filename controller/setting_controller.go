package controller

import (
	"context"

	"FinalProject/internal/fileopener"
	"FinalProject/model"
	"FinalProject/repository"
	"FinalProject/service"
)

type SettingController struct {
	service service.SettingService
}

func NewSettingController(service service.SettingService) *SettingController {
	return &SettingController{
		service: service,
	}
}

func (c *SettingController) GetSetting() (*model.Setting, error) {
	return c.service.GetSetting(context.Background())
}

func (c *SettingController) SaveSetting(source string, dest string) error {
	setting := &model.Setting{
		Source: source,
		Dest:   dest,
	}

	return c.service.SaveSetting(
		context.Background(),
		setting,
	)
}

func (c *SettingController) ShowDatabaseInFolder() error {
	databasePath, err := repository.DatabasePath()
	if err != nil {
		return err
	}
	return fileopener.Reveal(databasePath)
}
