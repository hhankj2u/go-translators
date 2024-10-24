package utils

import (
	"fmt"
	"net/url"
	"strings"

	"translators/internal/settings"

	"github.com/PuerkitoBio/goquery"
)

var DICTS = settings.DICTS

func MakeASoup(text string, baseURL string) (*goquery.Document, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(text))
	updateRelativeURLs(doc, baseURL)
	return doc, err
}

func ReplaceAll(input string) string {
	replacements := []struct {
		old string
		new string
	}{
		{"\n            (", "("},
		{"(\n    \t                ", "("},
		{"\n", " "},
		{"      \t                ", " "},
		{"\\'", "'"},
		{"  ", " "},
		{"\xa0 ", ""},
		{"[ ", "["},
		{" ]", "]"},
		{"A1", ""},
		{"A2", ""},
		{"B1", ""},
		{"B2", ""},
		{"C1", ""},
		{"C2", ""},
	}

	for _, r := range replacements {
		input = strings.ReplaceAll(input, r.old, r.new)
	}

	return strings.TrimSpace(input)
}

func ParseResponseURL(inputURL string) string {
	parts := strings.Split(inputURL, "?")
	return parts[0]
}

func GetRequestURL(baseURL, inputWord, dict string) string {
	if dict == DICTS[0] {
		queryWord := strings.ReplaceAll(inputWord, " ", "-")
		queryWord = strings.ReplaceAll(queryWord, "/", "-")
		return baseURL + queryWord
	}
	return baseURL + url.PathEscape(inputWord)
}

func GetRequestURLSpellcheck(baseURL, inputWord string) string {
	queryWord := strings.ReplaceAll(inputWord, " ", "+")
	queryWord = strings.ReplaceAll(queryWord, "/", "+")
	return baseURL + queryWord
}

func updateRelativeURLs(doc *goquery.Document, baseURL string) error {
	// Parse the base URL
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return fmt.Errorf("invalid base URL: %v", err)
	}

	// Update <a href> elements
	doc.Find("a[href]").Each(func(i int, sel *goquery.Selection) {
		href, exists := sel.Attr("href")
		if exists {
			newURL := resolveURL(parsedBaseURL, href)
			sel.SetAttr("href", newURL)
		}
	})

	// Update <img src> elements
	doc.Find("img[src]").Each(func(i int, sel *goquery.Selection) {
		src, exists := sel.Attr("src")
		if exists {
			newURL := resolveURL(parsedBaseURL, src)
			sel.SetAttr("src", newURL)
		}
	})

	// Update <link href> elements (CSS files, etc.)
	doc.Find("link[href]").Each(func(i int, sel *goquery.Selection) {
		href, exists := sel.Attr("href")
		if exists {
			newURL := resolveURL(parsedBaseURL, href)
			sel.SetAttr("href", newURL)
		}
	})

	// Update <script src> elements (JavaScript files, etc.)
	doc.Find("script[src]").Each(func(i int, sel *goquery.Selection) {
		src, exists := sel.Attr("src")
		if exists {
			newURL := resolveURL(parsedBaseURL, src)
			sel.SetAttr("src", newURL)
		}
	})

	return nil
}

// resolveURL converts a relative URL to an absolute URL based on the base URL
func resolveURL(base *url.URL, relativeURL string) string {
	// Parse the relative URL
	parsedURL, err := url.Parse(relativeURL)
	if err != nil || parsedURL.IsAbs() {
		// If it's already an absolute URL or invalid, return it as is
		return relativeURL
	}

	// Resolve relative URL against the base URL
	return base.ResolveReference(parsedURL).String()
}

// Extract the scheme and host (e.g., "https://dictionary.cambridge.org") from a URL
func ExtractSchemeAndHost(inputURL string) string {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return ""
	}
	return parsedURL.Scheme + "://" + parsedURL.Host
}