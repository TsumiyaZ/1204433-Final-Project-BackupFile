package main

import (
	"FinalProject/model"

	"gorm.io/gen"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "./model/query",
		Mode:    gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.ApplyBasic(
		model.Setting{},
		model.Photo{},
		model.Tag{},
		model.PhotoTag{},
	)

	g.Execute()
}
