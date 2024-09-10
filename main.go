package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"log"
	"net/url"
)

func fetchURLsFromPage(pageURL string) ([]string, error) {
	page := rod.New().MustConnect().MustPage(pageURL)
	links := page.MustElements("a")

	// Parser l'URL de base
	base, err := url.Parse(pageURL)
	if err != nil {
		return nil, err
	}

	var urls []string
	for _, link := range links {
		href := link.MustAttribute("href")
		if href != nil {
			// Résoudre l'URL relative
			parsedLink, err := url.Parse(*href)
			if err != nil {
				log.Println("Erreur parsing URL:", err)
				continue
			}
			fullURL := base.ResolveReference(parsedLink)
			urls = append(urls, fullURL.String())
		}
	}
	return urls, nil
}

func fetchRecursiveURLs(startURL string) []string {
	visited := make(map[string]bool)
	var result []string

	//// Fonction récursive pour parcourir les pages
	var recursiveFetch func(pageURL string)
	recursiveFetch = func(pageURL string) {
		if visited[pageURL] {
			return
		}
		visited[pageURL] = true

		urls, err := fetchURLsFromPage(pageURL)
		if err != nil {
			log.Println("Erreur fetching URLs from page:", err)
			return
		}

		result = append(result, urls...)

		for _, url := range urls {
			recursiveFetch(url)
		}
	}
	recursiveFetch(startURL)
	return result
}

func main() {
	baseURL := "http://158.178.197.230:8081/index.html"
	log.Println("Fetching URLs recursively")

	urls := fetchRecursiveURLs(baseURL)
	fmt.Println("All URLs found:", urls)
	log.Println("Done")
}
