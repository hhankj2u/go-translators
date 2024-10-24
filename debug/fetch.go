package main

import (
	"log"
	"translators/internal/cache"
	"translators/internal/dicts"
)

func main() {
	// test SearchCambridge
	con := cache.DB
	inputWord := "test"
	isFresh := false
	url, soup, err := dicts.SearchCambridge(con, inputWord, isFresh)
	if err != nil {
		panic(err)
	}
	html, err := soup.Html()
	if err != nil {
		panic(err)
	}
	log.Printf("url: %s, soup: %s", url, html)
}
