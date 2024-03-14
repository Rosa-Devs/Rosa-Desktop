package app

import (
	"context"
	"log"
	"os"

	"github.com/Rosa-Devs/Database/src/manifest"
	"github.com/Rosa-Devs/core/core"
	"github.com/Rosa-Devs/core/store"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	C        core.Core
	Addres   string
	wailsctx context.Context
}

func (a *App) GetAddres() string {
	log.Println("Addr:", a.Addres)
	return a.Addres
}

func (a *App) OnWailsInit(ctx context.Context) {
	a.wailsctx = ctx
}

func (a *App) Init(path string) {
	Store, err := store.NewStore(path)
	if err != nil {
		log.Println("Loaded new store")
	}
	a.C = core.Core{
		Store: *Store,
	}

	a.Addres = a.StartApi()
	a.StartManager()

	// go a.EventHandler()
}

func (a *App) StartApi() string {
	return a.C.StartApi("")
}

func (a *App) StartManager() {
	a.C.StartManager()
}

// func (a *App) EventHandler() {
// 	go func(A *App) {
// 		for {
// 			select {
// 			case <-a.C.EventCh:
// 				// Handle the event
// 				log.Println("EVENT")
// 				runtime.EventsEmit(A.wailsctx, "update")
// 				// Implement your event handling logic here
// 			}
// 		}
// 	}(a)
// }

func (a *App) GetProfile() string {
	return a.C.GetProfile()
}

func (a *App) ExportManifest(m manifest.Manifest) {
	path, err := runtime.SaveFileDialog(a.wailsctx, runtime.SaveDialogOptions{
		Title:                "Save chat file...",
		CanCreateDirectories: true,
		DefaultFilename:      m.Name + ".json",
	})
	if err != nil {
		log.Println("Fail to chosee export path")
		return
	}

	data, err := m.Serialize()
	if err != nil {
		log.Println("Fail to serialize")
		return
	}

	err = os.WriteFile(path, data, os.FileMode(0775))
	if err != nil {
		log.Println("Fail to write a file!!")
		return
	}

	log.Println(path)
}
