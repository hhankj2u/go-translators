package dicts

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"

	"github.com/hhankj2u/translators/pkg/settings"
	"github.com/hhankj2u/translators/pkg/utils"

	"github.com/PuerkitoBio/goquery"
)

type WebsterDictionary struct{}

const (
	WEBSTER_BASE_URL      = "https://www.merriam-webster.com"
	WEBSTER_DICT_BASE_URL = WEBSTER_BASE_URL + "/dictionary/"
)

// SearchWebster requests web resource and returns the result.
func (w WebsterDictionary) Search(con *sql.DB, inputWord string, isFresh bool) (string, *goquery.Document, error) {
	reqURL := utils.GetRequestURL(WEBSTER_DICT_BASE_URL, inputWord, settings.WEBSTER)

	if !isFresh {
		cached, soup, err := CacheRun(con, inputWord, reqURL)
		if err != nil {
			return "", nil, err
		}
		if !cached {
			return w.FreshRun(con, reqURL, inputWord)
		}
		return reqURL, soup, nil
	} else {
		return w.FreshRun(con, reqURL, inputWord)
	}
}

// FetchWebster gets response URL and response text for later parsing.
func (w WebsterDictionary) Fetch(reqURL, inputWord string) (bool, string, string, error) {
	client := &http.Client{}
	resp, err := Fetch(reqURL, client)
	if err != nil {
		return false, "", "", err
	}
	defer resp.Body.Close()

	resURL := resp.Request.URL.String()
	resText, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", "", err
	}

	if resp.StatusCode == 200 {
		fmt.Printf("%s \"%s\" in %s at %s\n", settings.OP[5], inputWord, settings.WEBSTER, resURL)
		return true, resURL, string(resText), nil
	}

	if resp.StatusCode == 404 {
		fmt.Printf("%s \"%s\" in %s\n", settings.OP[6], inputWord, settings.WEBSTER)
		return false, resURL, string(resText), nil
	}

	return false, resURL, string(resText), nil
}

// FreshRun prints the result without cache.
func (w WebsterDictionary) FreshRun(con *sql.DB, reqURL, inputWord string) (string, *goquery.Document, error) {
	found, resURL, resText, err := w.Fetch(reqURL, inputWord)
	if err != nil {
		return "", nil, err
	}

	if found {
		soup, err := utils.MakeASoup(resText)
		if err != nil {
			return "", nil, err
		}
		responseWord := w.ParseResponseWord(soup)
		expected := soup.Find("div#left-content")
		html, err := expected.Html()
		if err != nil {
			return "", nil, err
		}
		soup.Find("body").SetHtml(html)

		html, err = soup.Html()
		if err != nil {
			return "", nil, err
		}
		Save(con, inputWord, responseWord, resURL, html)
		return resURL, soup, nil
	} else {
		fmt.Printf("%s the parsed result of %s\n", settings.OP[4], resURL)

		soup, err := utils.MakeASoup(resText)
		if err != nil {
			return "", nil, err
		}
		nodes := soup.Find("div.widget.spelling-suggestion")
		html, err := nodes.Html()
		if err != nil {
			return "", nil, err
		}
		soup.Find("body").SetHtml(html)
		return resURL, soup, nil
	}
}

// ParseResponseWord parses the response word from h1 tag.
func (w WebsterDictionary) ParseResponseWord(soup *goquery.Document) string {
	return soup.Find("h1.hword").Text()
}
