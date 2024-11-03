package dicts

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hhankj2u/translators/pkg/settings"
	"github.com/hhankj2u/translators/pkg/utils"

	"github.com/PuerkitoBio/goquery"
)

type CambridgeDictionary struct{}

const (
	CAMBRIDGE_URL            = "https://dictionary.cambridge.org"
	CAMBRIDGE_DICT_BASE_URL  = CAMBRIDGE_URL + "/dictionary/english/"
	CAMBRIDGE_SPELLCHECK_URL = CAMBRIDGE_URL + "/spellcheck/english/?q="
)

// searchCambridge requests web resource and returns the result.
func (c CambridgeDictionary) Search(con *sql.DB, inputWord string, isFresh bool) (string, *goquery.Document, error) {
	reqURL := utils.GetRequestURL(CAMBRIDGE_DICT_BASE_URL, inputWord, settings.CAMBRIDGE)

	if !isFresh {
		cached, soup, err := CacheRun(con, inputWord, reqURL)
		if err != nil {
			return "", nil, err
		}
		if !cached {
			return c.FreshRun(con, reqURL, inputWord)
		}
		return reqURL, soup, nil
	} else {
		return c.FreshRun(con, reqURL, inputWord)
	}
}

// fetchCambridge gets response URL and response text for later parsing.
func (c CambridgeDictionary) Fetch(reqURL, inputWord string) (bool, string, string, error) {
	client := &http.Client{}
	resp, err := Fetch(reqURL, client)
	if err != nil {
		return false, "", "", err
	}
	defer resp.Body.Close()

	if resp.Request.URL.String() == CAMBRIDGE_DICT_BASE_URL {
		fmt.Printf("%s \"%s\" in %s\n", settings.OP[6], inputWord, settings.CAMBRIDGE)
		spellReqURL := utils.GetRequestURLSpellcheck(CAMBRIDGE_SPELLCHECK_URL, inputWord)

		spellResp, err := Fetch(spellReqURL, client)
		if err != nil {
			return false, "", "", err
		}
		defer spellResp.Body.Close()
		spellResText, err := io.ReadAll(spellResp.Body)
		if err != nil {
			return false, "", "", err
		}
		return false, spellResp.Request.URL.String(), string(spellResText), nil
	} else {
		resURL := utils.ParseResponseURL(resp.Request.URL.String())
		resText, err := io.ReadAll(resp.Body)
		if err != nil {
			return false, "", "", err
		}

		fmt.Printf("%s \"%s\" in %s at %s\n", settings.OP[5], inputWord, settings.CAMBRIDGE, resURL)
		return true, resURL, string(resText), nil
	}
}

// freshRun prints the result without cache.
func (c CambridgeDictionary) FreshRun(con *sql.DB, reqURL, inputWord string) (string, *goquery.Document, error) {
	found, resURL, resText, err := c.Fetch(reqURL, inputWord)
	if err != nil {
		return "", nil, err
	}

	if found {
		soup, err := utils.MakeASoup(resText)
		if err != nil {
			return "", nil, err
		}
		responseWord := c.ParseResponseWord(soup)
		expected := soup.Find("div.page")
		html, err := expected.Html()
		if err != nil {
			return "", nil, err
		}

		// hide #onetrust-consent-sdk
		html += `<style>#onetrust-consent-sdk {display: none;}</style>`

		soup.Find("body").SetHtml(html)

		html, err = soup.Html()
		if err != nil {
			return "", nil, err
		}
		Save(con, inputWord, responseWord, resURL, html)
		return resURL, soup, nil
	} else {
		spellResURL, spellResText := resURL, resText
		fmt.Printf("%s the parsed result of %s\n", settings.OP[4], spellResURL)

		soup, err := utils.MakeASoup(spellResText)
		if err != nil {
			return "", nil, err
		}
		nodes := soup.Find("div.hfl-s.lt2b.lmt-10.lmb-25.lp-s_r-20").Find("ul.hul-u")
		html, err := nodes.Html()
		if err != nil {
			return "", nil, err
		}
		soup.Find("body").SetHtml(html)
		return spellResURL, soup, nil
	}
}

// parseResponseWord parses the response word from HTML head title tag.
func (c CambridgeDictionary) ParseResponseWord(soup *goquery.Document) string {
	temp := strings.TrimSpace(soup.Find("title").Text())
	temp = strings.Split(temp, "-")[0]
	if strings.Contains(temp, "|") {
		return strings.ToLower(strings.TrimSpace(strings.Split(temp, "|")[0]))
	}
	return strings.ToLower(temp)
}
