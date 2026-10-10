package main

import (
	"context"
	"embed"

	"FinalProject/controller"
	"FinalProject/repository"
	"FinalProject/service"
	"FinalProject/service/backup"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	db, err := repository.NewDBConnection()
	if err != nil {
		println("Database connection error:", err.Error())
		return
	}

	err = repository.MigrateDatabase(db)
	if err != nil {
		println("Database migration error:", err.Error())
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		println("Database instance error:", err.Error())
		return
	}
	defer sqlDB.Close()

	println("Connected Database")

	settingRepo := repository.NewSettingRepository(db)
	settingService := service.NewSettingService(settingRepo)
	settingController := controller.NewSettingController(settingService)

	photoRepo := repository.NewPhotoRepository(db)
	photoService := service.NewPhotoService(photoRepo)
	photoController := controller.NewPhotoController(photoService)

	backupRepo := repository.NewBackupRepository(db)
	backupService := backup.NewBackupService(backupRepo, settingRepo)
	backupController := controller.NewBackupController(backupService)

	app := NewApp()

	err = wails.Run(&options.App{
		Title:  "FinalProject",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			controller.SetBackupContext(backupController, ctx)
		},
		Bind: []interface{}{
			app,
			settingController,
			photoController,
			backupController,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
