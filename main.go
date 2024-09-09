package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"log"
	"net/url"
)

// http://158.178.197.230:8081/index.html
func fetchURL() {
	log.Println("Fetching URL")

	baseURL := "https://kodoka.fr/"
	page := rod.New().MustConnect().MustPage(baseURL)
	links := page.MustElements("a")

	// Parser l'URL de base
	base, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
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
	fmt.Println(urls)
	log.Println("URL fetched")
}

func main() {
	fetchURL()
}
