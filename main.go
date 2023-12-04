package main

import (
	"changeme/src"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {
	// Create an instance of the app structure
	app := NewApp()
	DbManager := src.DbManager{}

	// Create application with options
	err := wails.Run(&options.App{
		Title:            "Rosa",
		Width:            1650,
		Height:           1030,
		Assets:           assets,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		//Frameless:        true,
		CSSDragProperty: "--wails-draggable",
		CSSDragValue:    "drag",

		Bind: []interface{}{
			app,
			&DbManager,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
