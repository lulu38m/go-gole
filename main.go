package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"log"
	"net/url"
)

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
			fmt.Println(fullURL.String())
		}
	}

	log.Println("URL fetched")
}

func main() {
	fetchURL()
}
