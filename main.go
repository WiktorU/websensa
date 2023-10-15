package main

import (
	"fmt"
)

func main() {
	startURLs := []string{"https://www.onet.pl", "https://www.wp.pl", "https://www.pudelek.pl", "https://example.com"}

	// Limit the number of workers to 5
	scraper := NewScraper(5)

	scraper.Scrape(startURLs)

	mostRecurringWords := scraper.FindMostRecurringWords()

	fmt.Printf(`The most recurring word is 
	1.'%s' with a count of %d\n
	2.'%s' with a count of %d\n
	3.'%s' with a count of %d\n
	4.'%s' with a count of %d\n
	5.'%s' with a count of %d\n`,
		mostRecurringWords[0].Key, mostRecurringWords[0].Value,
		mostRecurringWords[1].Key, mostRecurringWords[1].Value,
		mostRecurringWords[2].Key, mostRecurringWords[2].Value,
		mostRecurringWords[3].Key, mostRecurringWords[3].Value,
		mostRecurringWords[4].Key, mostRecurringWords[4].Value)
}
