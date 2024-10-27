package main

import (
	"context"
	"database/sql"
	"log"
	"translators/internal/cache"
	"translators/internal/dicts"
	"translators/internal/settings"
)

type App struct {
    dbConnection *sql.DB
}

func NewApp() *App {
    // Initialize database or any other required setup
    con := cache.InitDB(settings.DICTS[0])
    return &App{dbConnection: con}
}

func (a *App) Startup(ctx context.Context) {
    log.Println("App is starting up!")
}

// SearchDictionary performs the dictionary search and returns the result HTML
func (a *App) SearchDictionary(term string) (string, error) {
    _, soup, err := dicts.SearchCambridge(a.dbConnection, term, false)
    if err != nil {
        return "", err
    }
    html, err := soup.Html()
    if err != nil {
        return "", err
    }
    return html, nil
}
