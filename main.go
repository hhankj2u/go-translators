package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
    // Initialize the App instance
    app := NewApp()

    err := wails.Run(&options.App{
        Title:     "Translators",
        Width:     1024,
        Height:    768,
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
        Bind:      []interface{}{app}, // Bind `App` instance for frontend access
    })
    if err != nil {
        log.Fatal(err)
    }
}
