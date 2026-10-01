package controller

import (
	"context"

	"FinalProject/model"
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
