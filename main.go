package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var visited = make(map[string]bool)
var finalURLs = make([]string, 0)
var mu sync.Mutex // Pour protéger les accès concurrents à visited et finalURLs

const maxURLs = 1000
const maxWorkers = 100
const maxDepth = 5

var queue = make(chan Task, maxURLs) // File d'attente d'URLs avec profondeur

type Task struct {
	URL   string
	Depth int
}

// Worker function
func worker(id int, jobs <-chan Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range jobs {

		//fmt.Printf("Worker %d: Processing URL: %s (Depth %d) %d\n", id, task.URL, task.Depth, len(queue))
		fetch(task.URL, task.Depth)
		time.Sleep(3 * time.Second) // Attente pour éviter l'erreur 429
	}
}

// Fonction pour récupérer et traiter une URL
func fetch(urlStr string, depth int) {
	mu.Lock()
	if visited[urlStr] || len(finalURLs) >= maxURLs {
		mu.Unlock()
		return
	}
	visited[urlStr] = true
	mu.Unlock()

	res, err := http.Get(urlStr)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Printf("Error: Status code %d while fetching %s\n", res.StatusCode, urlStr)
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println(err)
		return
	}

	doc.Find("a").Each(func(index int, item *goquery.Selection) {
		if depth >= maxDepth {
			return
		}

		link, exists := item.Attr("href")
		if exists {
			absoluteURL := resolveURL(urlStr, link)
			if strings.HasPrefix(absoluteURL, "http") && !strings.Contains(absoluteURL, "#") {
				mu.Lock()
				if !visited[absoluteURL] {
					finalURLs = append(finalURLs, absoluteURL)
					//fmt.Println("URL found:", absoluteURL)
					queue <- Task{URL: absoluteURL, Depth: depth + 1}
				}
				mu.Unlock()
			}
		}
	})
}

// Fonction pour résoudre les URLs relatives en URLs absolues
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

// Fonction principale
func main() {
	startURL := "https://www.kodoka.fr/"

	var wg sync.WaitGroup
	for i := 1; i <= maxWorkers; i++ {
		wg.Add(1)
		go worker(i, queue, &wg)
	}

	queue <- Task{URL: startURL, Depth: 0}

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for range ticker.C {
			fmt.Println("Visited URLs:", len(finalURLs))
		}
	}()

	wg.Wait()
	close(queue)

	fmt.Println("Crawling finished. Found", len(finalURLs), "URLs.")
	fmt.Println("Final URLs:")
}
