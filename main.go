package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"log"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func fetchURL() {
	baseURL := "http://158.178.197.230:8081/index.html"

	// Request the HTML page.
	res, err := http.Get(baseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	log.Print()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Parser l'URL de base
	base, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}

	var urls []string

	// trouver les liens
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		link, ok := s.Attr("href")
		if !ok {
			return
		}

		parsedLink, err := url.Parse(link)
		if err != nil {
			log.Println("Erreur parsing URL:", err)
			return
		}

		fullURL := base.ResolveReference(parsedLink)

		urls = append(urls, fullURL.String())
	})

	fmt.Println(urls)

}

func main() {
	log.Println("Fetching URL")
	page := rod.New().MustConnect().MustPage("https://www.instagram.com/")
	page.MustWaitStable().MustScreenshot("a.png")
	log.Println("URL fetched")
}
