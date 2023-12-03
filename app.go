package main

import (
	"context"
)

type Contact struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"imageUrl"`
}

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

func (a *App) GetContacts() []Contact {
	contacts := []Contact{
		{ID: 1, Name: "John Doe", ImageURL: "https://cdn3.iconfinder.com/data/icons/shinysocialball/512/Technorati_512x512.png"},
		{ID: 2, Name: "Jane Smith", ImageURL: "https://ccia.ugr.es/cvg/CG/images/base/5.gif"},
		{ID: 3, Name: "Alice Johnson", ImageURL: "https://upload.wikimedia.org/wikipedia/commons/c/cc/Icon_Pinguin_1_512x512.png"},
		{ID: 4, Name: "Mihalic2040", ImageURL: "https://avatarfiles.alphacoders.com/762/thumb-76262.png"},
		// Add more contacts as needed
	}

	return contacts
}
