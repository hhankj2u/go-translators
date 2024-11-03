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

type SohaDictionary struct{}

const (
	SOHA_BASE_URL      = "http://tratu.soha.vn"
	SOHA_DICT_BASE_URL = SOHA_BASE_URL + "/dict/en_vn/"
)

// Search requests web resource and returns the result.
func (s SohaDictionary) Search(con *sql.DB, inputWord string, isFresh bool) (string, *goquery.Document, error) {
	reqURL := utils.GetRequestURL(SOHA_DICT_BASE_URL, inputWord, settings.SOHA)

	if !isFresh {
		cached, soup, err := CacheRun(con, inputWord, reqURL)
		if err != nil {
			return "", nil, err
		}
		if !cached {
			return s.FreshRun(con, reqURL, inputWord)
		}
		return reqURL, soup, nil
	} else {
		return s.FreshRun(con, reqURL, inputWord)
	}
}

// Fetch gets response URL and response text for later parsing.
func (s SohaDictionary) Fetch(reqURL, inputWord string) (bool, string, string, error) {
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
		fmt.Printf("%s \"%s\" in %s at %s\n", settings.OP[5], inputWord, settings.SOHA, resURL)
		return true, resURL, string(resText), nil
	}

	if resp.StatusCode == 404 {
		fmt.Printf("%s \"%s\" in %s\n", settings.OP[6], inputWord, settings.SOHA)
		return false, resURL, string(resText), nil
	}

	return false, resURL, string(resText), nil
}

// FreshRun prints the result without cache.
func (s SohaDictionary) FreshRun(con *sql.DB, reqURL, inputWord string) (string, *goquery.Document, error) {
	found, resURL, resText, err := s.Fetch(reqURL, inputWord)
	if err != nil {
		return "", nil, err
	}

	if found {
		soup, err := utils.MakeASoup(resText)
		if err != nil {
			return "", nil, err
		}
		responseWord := s.ParseResponseWord(soup)
		expected := soup.Find("div#column-content")
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
		nodes := soup.Find("div#column-content")
		html, err := nodes.Html()
		if err != nil {
			return "", nil, err
		}
		soup.Find("body").SetHtml(html)
		return resURL, soup, nil
	}
}

// ParseResponseWord parses the response word from h1 tag.
func (s SohaDictionary) ParseResponseWord(soup *goquery.Document) string {
	return soup.Find("h1.firstHeading").Text()
}
