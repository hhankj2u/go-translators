package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"translators/internal/cache"
	"translators/internal/dicts"
	"translators/internal/settings"
)

type App struct {
	dbConnections map[string]*sql.DB
}

func NewApp() *App {
	// Initialize database or any other required setup
	dbConnections := make(map[string]*sql.DB)
	dbConnections[settings.WEBSTER] = cache.InitDB(settings.WEBSTER)
	dbConnections[settings.CAMBRIDGE] = cache.InitDB(settings.CAMBRIDGE)
	dbConnections[settings.SOHA] = cache.InitDB(settings.SOHA)
	return &App{dbConnections: dbConnections}
}

func (a *App) Startup(ctx context.Context) {
	log.Println("App is starting up!")
}

// SearchDictionary performs the dictionary search and returns the result HTML from all dictionaries
func (a *App) SearchDictionary(term string) (map[string]string, error) {
	dictionaries := map[string]dicts.Dictionary{
		settings.WEBSTER:   dicts.WebsterDictionary{},
		settings.CAMBRIDGE: dicts.CambridgeDictionary{},
		settings.SOHA:      dicts.SohaDictionary{},
	}

	results := make(map[string]string)
	var mu sync.Mutex
	var wg sync.WaitGroup
	errChan := make(chan error, len(dictionaries))

	for name, dictionary := range dictionaries {
		wg.Add(1)
		go func(name string, dictionary dicts.Dictionary) {
			defer wg.Done()
			con := a.dbConnections[name]
			_, soup, err := dictionary.Search(con, term, false)
			if err != nil {
				errChan <- fmt.Errorf("error searching %s dictionary: %w", name, err)
				return
			}
			html, err := soup.Html()
			if err != nil {
				errChan <- fmt.Errorf("error getting HTML from %s dictionary: %w", name, err)
				return
			}
			mu.Lock()
			results[name] = html
			mu.Unlock()
		}(name, dictionary)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return nil, <-errChan
	}

	return results, nil
}
