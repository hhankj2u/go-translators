package main

import (
    "os"

    "github.com/therecipe/qt/core"
    "github.com/therecipe/qt/widgets"
    "github.com/therecipe/qt/webengine" // Correct import for webengine
)

func main() {
    // Initialize Qt application
    app := widgets.NewQApplication(len(os.Args), os.Args)

    // Create the main window
    window := widgets.NewQMainWindow(nil, 0)
    window.SetWindowTitle("Qt WebEngine Example")
    window.SetMinimumSize2(1024, 768)

    // Create a central widget and layout
    centralWidget := widgets.NewQWidget(nil, 0)
    layout := widgets.NewQVBoxLayout()
    centralWidget.SetLayout(layout)
    window.SetCentralWidget(centralWidget)

    // Create a WebEngineView for rendering web content
    webView := webengine.NewQWebEngineView(nil)
    if webView == nil {
        panic("Failed to create QWebEngineView") // Check if webView is initialized properly
    }

    // Load a URL in the WebEngineView
    webView.SetUrl(core.NewQUrl3("https://www.example.com", 0))

    // Add the WebEngineView to the layout
    layout.AddWidget(webView, 0, 0)

    // Show the window
    window.Show()

    // Run the application
    app.Exec()
}
