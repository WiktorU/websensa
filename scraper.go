package main

import (
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

// Scraper represents a web scraper with an HTTP client.
type Scraper struct {
	client     *http.Client
	WordCount  map[string]int
	mu         sync.Mutex
	maxWorkers int
	wg         sync.WaitGroup
	workQueue  *URLQueue
}

// NewScraper initializes a new Scraper with the specified maximum number of workers.
func NewScraper(maxWorkers int) *Scraper {
	return &Scraper{
		client:     http.DefaultClient,
		WordCount:  make(map[string]int),
		maxWorkers: maxWorkers,
		workQueue:  NewURLQueue(),
	}
}

// ScrapeURL scrapes a given URL and counts word occurrences.
func (s *Scraper) ScrapeURL(url url.URL) {
	defer s.wg.Done()

	wordCount := make(map[string]int)

	resp, err := s.client.Get(url.String())
	if err != nil {
		log.Printf("Error fetching URL %s: %v", url.String(), err)
		return
	}
	defer resp.Body.Close()

	tokenizer := html.NewTokenizer(resp.Body)
	inStyleElement := false
	exit := false

	for !exit {
		tokenType := tokenizer.Next()

		switch tokenType {
		case html.ErrorToken:
			exit = true
		case html.TextToken:
			if inStyleElement {
				continue
			}

			text := tokenizer.Token().Data
			words := strings.Fields(text)

			for _, word := range words {
				word = strings.ToLower(word)
				wordCount[word]++
			}

		case html.StartTagToken, html.EndTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			if token.Data == "style" || token.Data == "script" || token.Data == "link" {
				if tokenType == html.StartTagToken {
					inStyleElement = true
				} else if tokenType == html.EndTagToken {
					inStyleElement = false
				}
			}

			if subsiteURL, ok := s.getURLStringFromToken(token, url); ok {
				parsedURL, err := url.Parse(subsiteURL)
				if err != nil {
					log.Fatal(err)
				}

				s.workQueue.Enqueue(*parsedURL)
			}
		}
	}

	// Merge word counts into the shared map
	s.mu.Lock()
	for word, count := range wordCount {
		s.WordCount[word] += count
	}
	s.mu.Unlock()
}

// getURLStringFromToken retrieves the url string from the HTML token
func (s *Scraper) getURLStringFromToken(token html.Token, mainURL url.URL) (string, bool) {
	if token.Data == "a" {
		for _, attr := range token.Attr {
			if attr.Key == "href" {
				subsiteURL := attr.Val
				parsedSubsiteURL, err := url.Parse(subsiteURL)
				if err == nil && parsedSubsiteURL.Host == mainURL.Host {
					return subsiteURL, true
				}
			}
		}
	}

	return "", false
}

// Scrape initiates the scraping process from a specified URL.
func (s *Scraper) Scrape(startURLs []string) {
	for _, startURL := range startURLs {
		parsedURL, err := url.Parse(startURL)
		if err != nil {
			log.Fatal(err)
		}

		s.workQueue.Enqueue(*parsedURL)
	}

	// Start multiple goroutines to process URLs from the queue
	for i := 0; i < s.maxWorkers; i++ {
		go func() {
			for {
				if urlToScrape, ok := s.workQueue.Dequeue(); ok {
					s.wg.Add(1)
					s.ScrapeURL(urlToScrape)
				}
			}
		}()
	}

	// Sleep for 2 seconds
	duration := 5 * time.Second
	time.Sleep(duration)

	s.wg.Wait()
}

// FindMostRecurringWords finds the most recurring words in the WordCount map.
func (s *Scraper) FindMostRecurringWords() []KeyValue {
	var keyValuePairs []KeyValue
	for key, value := range s.WordCount {
		keyValuePairs = append(keyValuePairs, KeyValue{key, value})
	}

	sort.Slice(keyValuePairs, func(i, j int) bool {
		return keyValuePairs[i].Value > keyValuePairs[j].Value
	})

	return keyValuePairs
}

// Define a struct to store key-value pairs.
type KeyValue struct {
	Key   string
	Value int
}
