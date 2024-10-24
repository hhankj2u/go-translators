package dicts

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"translators/internal/cache"
	"translators/internal/settings"
	"translators/internal/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/corpix/uarand"
	"github.com/mattn/go-sqlite3"
)

// fetch makes a web request with retry mechanism.
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

// cacheRun checks the cache is from Cambridge or Merrian Webster.
func CacheRun(con *sql.DB, inputWord, reqURL string) (bool, *goquery.Document, error) {
	// data is a tuple (response_url, response_text) if any
	response_url, response_text, err := cache.GetCache(con, inputWord, reqURL)
	if err != nil {
		return false, nil, err
	}

	if response_text != "" {
		doc, err := utils.MakeASoup(string(response_text), utils.ExtractSchemeAndHost(response_url))
		if err != nil {
			return false, nil, err
		}
		return true, doc, nil
	}

	return false, nil, nil
}

// save saves a word info into local DB for cache.
func Save(con *sql.DB, inputWord, responseWord, responseURL, responseText string) {
	err := cache.InsertIntoTable(con, inputWord, responseWord, responseURL, responseText)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.Code == sqlite3.ErrConstraint {
			log.Printf("%s caching \"%s\" - [ERROR] - already cached before\n", settings.OP[8], inputWord)
		} else {
			log.Printf("%s caching \"%s\" - [ERROR] - %v\n", settings.OP[8], inputWord, err)
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
