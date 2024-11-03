package dicts

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/hhankj2u/translators/pkg/cache"
	"github.com/hhankj2u/translators/pkg/settings"
	"github.com/hhankj2u/translators/pkg/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/corpix/uarand"
	"github.com/mattn/go-sqlite3"
)

// Dictionary defines the interface for dictionary operations.
type Dictionary interface {
	Search(con *sql.DB, inputWord string, isFresh bool) (string, *goquery.Document, error)
	Fetch(reqURL, inputWord string) (bool, string, string, error)
	FreshRun(con *sql.DB, reqURL, inputWord string) (string, *goquery.Document, error)
	ParseResponseWord(soup *goquery.Document) string
}

// Fetch makes a web request with retry mechanism.
func Fetch(url string, client *http.Client) (*http.Response, error) {
	ua := uarand.GetRandom()
	headers := map[string]string{"User-Agent": ua}
	for key, value := range headers {
		client.Transport = &transport{http.DefaultTransport, key, value}
	}
	attempt := 0

	for {
		log.Printf("%s %s", settings.OP[0], url)
		resp, err := client.Get(url)
		if err != nil {
			var retry bool
			attempt, retry = callOnError(err, url, attempt, settings.OP[2])
			if retry {
				continue
			}
			return nil, err
		}
		return resp, nil
	}
}

type transport struct {
	Transport http.RoundTripper
	Key       string
	Value     string
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set(t.Key, t.Value)
	return t.Transport.RoundTrip(req)
}

// CacheRun checks the cache is from Cambridge or Merriam-Webster.
func CacheRun(con *sql.DB, inputWord, reqURL string) (bool, *goquery.Document, error) {
	// data is a tuple (response_url, response_text) if any
	_, responseText, err := cache.GetCache(con, inputWord, reqURL)
	if err != nil {
		return false, nil, err
	}

	if responseText != "" {
		doc, err := utils.MakeASoup(string(responseText))
		if err != nil {
			return false, nil, err
		}
		return true, doc, nil
	}

	return false, nil, nil
}

// Save saves a word info into local DB for cache.
func Save(con *sql.DB, inputWord, responseWord, responseURL, responseText string) {
	err := cache.InsertIntoTable(con, inputWord, responseWord, responseURL, responseText)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.Code, sqlite3.ErrConstraint) {
			log.Printf("%s caching \"%s\" - [ERROR] - already cached before\n", settings.OP[8], inputWord)
		}
	} else {
		log.Printf("%s the search result of \"%s\"", settings.OP[7], inputWord)
	}
}

// callOnError handles errors and determines if a retry is needed.
func callOnError(err error, url string, attempt int, op string) (int, bool) {
	attempt++
	if attempt >= 3 {
		log.Printf("%s %s - [ERROR] - %v\n", op, url, err)
		return attempt, false
	}
	time.Sleep(2 * time.Second)
	return attempt, true
}
