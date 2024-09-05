package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func fetchURL() {
	// Request the HTML page.
	res, err := http.Get("http://158.178.197.230:8081")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		link, ok := s.Attr("href")
		if !ok {
			return
		}
		fmt.Println(link)
	})
}

func main() {
	fetchURL()
}
