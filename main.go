package main

import (
	"embed"
	"log"
	"os"

	"github.com/Rosa-Devs/Rosa-Desktop/core"
	"github.com/Rosa-Devs/Rosa-Desktop/store"
	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {

	// Create an instance of the app
	path := ""
	if len(os.Args) > 1 {
		if os.Args[1] != "" {
			path = os.Args[1]
		}
	}
	Store, err := store.NewStore(path)
	if err != nil {
		log.Println("Loaded new store")
	}
	Core := core.Core{
		Icon:  icon,
		Store: *Store,
	}

	systray.Register(Core.TrayReady, Core.TrayExit)

	// Create application with options
	err = wails.Run(&options.App{
		Title:             "Rosa",
		Width:             1650,
		Height:            1030,
		Assets:            assets,
		BackgroundColour:  &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:         Core.OnWailsInit,
		HideWindowOnClose: true,
		//Frameless:        true,
		CSSDragProperty: "--wails-draggable",
		CSSDragValue:    "drag",

		Bind: []interface{}{
			&Core,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
