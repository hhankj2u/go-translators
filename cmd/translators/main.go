package main

import (
	"log"
	"translators/internal/cache"
	"translators/internal/dicts"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/webengine"
	"github.com/therecipe/qt/widgets"
)

// "translators/internal/settings"

func main() {
	// Call the function from dicts package and log to console
	con := cache.DB
	_, soup, err := dicts.SearchCambridge(con, "banana", false)
	if err != nil {
		log.Println(err)
	}
	// log.Println(soup.Text())

	// w := webview.New(false)
	// defer w.Destroy()
	// w.SetTitle("Basic Example")
	// w.SetSize(480, 320, webview.HintNone)
	html, err := soup.Html()
	if err != nil {
		log.Println(err)
		return
	}
	// w.Navigate(html)
	// w.Run()

	// Initialize the Qt application
	app := widgets.NewQApplication(len([]string{}), []string{})

	// Create a new window
	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("QWebEngineView Example with Base URL")
	window.SetMinimumSize2(800, 600)

	// Create a new QWebEngineView (for loading and displaying web content)
	webView := webengine.NewQWebEngineView(nil)

	// The HTML content to load
	// htmlContent := `
	// <html>
	// <head>
	// 		<title>Test Page</title>
	// 		<link rel="stylesheet" href="styles.css">
	// </head>
	// <body>
	// 		<h1>Hello from QWebEngineView</h1>
	// 		<img src="image.png">
	// </body>
	// </html>`

	// Base URL to resolve relative URLs (this is optional, but important to resolve relative paths)
	baseURL := "https://www.example.com/"
	qBaseURL := core.NewQUrl3(baseURL, core.QUrl__TolerantMode)

	// Load the HTML content with the base URL
	webView.SetHtml(html, qBaseURL)

	// Set the webView as the central widget of the window
	window.SetCentralWidget(webView)

	// Show the window
	window.Show()

	// Execute the Qt application
	app.Exec()
}
