package main

import (
	"log"

	"translators/internal/cache"
	"translators/internal/dicts"
	"translators/internal/settings"

	webview "github.com/hhankj2u/webview_go"
)

// "translators/internal/settings"

func main() {
	// Call the function from dicts package and log to console
	con := cache.InitDB(settings.DICTS[0])
	_, soup, err := dicts.SearchCambridge(con, "banana", false)
	if err != nil {
		log.Println(err)
	}

	w := webview.New(false)
	defer w.Destroy()
	w.SetTitle("Basic Example")
	w.SetSize(1000, 600, webview.HintNone)
	html, err := soup.Html()
	if err != nil {
		log.Println(err)
		return
	}
	w.SetHtml(html)
	w.Run()
}
