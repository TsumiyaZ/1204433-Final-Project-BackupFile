package main

import (
	"context"
	"embed"

	"FinalProject/ai/llama"
	"FinalProject/controller"
	"FinalProject/repository"
	"FinalProject/service"
	"FinalProject/service/backup"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/.output/public
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

	embeddingEngine := llama.NewClient("http://127.0.0.1:8080")

	settingRepo := repository.NewSettingRepository(db)
	settingService := service.NewSettingService(settingRepo)
	settingController := controller.NewSettingController(settingService)

	photoRepo := repository.NewPhotoRepository(db)
	photoService := service.NewPhotoService(photoRepo, embeddingEngine)
	embeddingQueue := service.NewEmbeddingQueue(photoRepo, settingRepo, photoService)
	photoController := controller.NewPhotoController(photoService)

	backupRepo := repository.NewBackupRepository(db)
	backupService := backup.NewBackupService(backupRepo, settingRepo, embeddingQueue)
	backupController := controller.NewBackupController(backupService)

	tagRepo := repository.NewTagRepository(db)
	galleryService := service.NewGalleryService(photoRepo, settingRepo, photoService, tagRepo)
	galleryController := controller.NewGalleryController(galleryService, embeddingQueue)

	app := NewApp()

	err = wails.Run(&options.App{
		Title:  "FinalProject",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 247, G: 241, B: 231, A: 255},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			controller.SetBackupContext(backupController, ctx)
			app.setEmbeddingStatus("loading", "กำลังตรวจสอบส่วนประกอบสำหรับค้นหารูป...")
			go func() {
				if err := llama.StartServer(ctx, func(progress llama.SetupProgress) {
					if progress.Determinate {
						app.setEmbeddingProgress(progress.Message, progress.Percent, true)
						return
					}
					app.setEmbeddingStatus("loading", progress.Message)
				}); err != nil {
					app.setEmbeddingStatus(
						"error",
						"เตรียม llama.cpp ไม่สำเร็จ: "+err.Error(),
					)
					return
				}

				app.setEmbeddingStatus("loading", "กำลังตรวจสอบและเตรียมโมเดล อาจใช้เวลาสักครู่...")
				if err := embeddingEngine.WaitUntilReady(ctx); err != nil {
					if ctx.Err() == nil {
						app.setEmbeddingStatus(
							"error",
							"เตรียมโมเดลไม่สำเร็จ ตรวจสอบอินเทอร์เน็ต แล้วเปิดแอปใหม่: "+err.Error(),
						)
					}
					return
				}
				if err := embeddingQueue.Start(ctx, func(update service.EmbeddingStatusUpdate) {
					runtime.EventsEmit(ctx, "embedding:status", update)
				}); err != nil {
					app.setEmbeddingStatus(
						"error",
						"โมเดลพร้อมแล้ว แต่เริ่มคิววิเคราะห์รูปไม่สำเร็จ: "+err.Error(),
					)
					return
				}
				app.setEmbeddingStatus("ready", "โมเดลพร้อมสร้าง embedding และ tag รูปภาพแล้ว")
			}()
		},
		Bind: []interface{}{
			app,
			settingController,
			photoController,
			backupController,
			galleryController,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
