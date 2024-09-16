package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var visited = make(map[string]bool)

func fetch(urlStr string) {
	// Check if the URL has been visited
	if visited[urlStr] {
		return
	}
	visited[urlStr] = true // Mark URL as visited

	// Fetch the URL
	res, err := http.Get(urlStr)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("Error: Status code %d while fetching %s", res.StatusCode, urlStr)
	}

	// Parse the page
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find all links and crawl them
	doc.Find("a").Each(func(index int, item *goquery.Selection) {
		link, exists := item.Attr("href")
		if exists {
			absoluteURL := resolveURL(urlStr, link)
			if strings.HasPrefix(absoluteURL, "http") && !strings.Contains(absoluteURL, "#") {
				fmt.Println("URL found:", absoluteURL)
				fetch(absoluteURL) // Recursively fetch the URL
			}
		}
	})
}

func resolveURL(base, href string) string {
	u, err := url.Parse(href)
	if err != nil {
		return ""
	}
	baseURL, err := url.Parse(base)
	if err != nil {
		return ""
	}
	return baseURL.ResolveReference(u).String()
}

func main() {
	startURL := "https://www.kodoka.fr/index.php"
	fmt.Println("Fetching URLs recursively")
	fetch(startURL)
	fmt.Println("Done")
}
