package spider

import (
	"fmt"
	"log"

	"github.com/kawilkinson/search-engine/internal/redis_db"
)

func (cfg *Config) CrawlPage(rawCurrentURL string, db *redis_db.RedisDatabase) {
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.Wg.Done()
	}()

	if cfg.isMaxPagesReached() {
		cfg.runOnce.Do(func() {
			fmt.Printf("max pages reached...\n")
		})
		return
	}

	log.Println("Waiting for the message queue...")
	rawCurrentURL, depth, normalizedURL, err := db.PopURL()
	if err != nil {
		fmt.Printf("No more URLs in the message queue: %v\n", err)
	}

	fmt.Printf("Popped URL: %v | Depth Level: %v | Normalized URL: %v\n", rawCurrentURL, depth, normalizedURL)

	// normalizedURL, err := NormalizeURL(rawCurrentURL)
	// if err != nil {
	// 	fmt.Printf("error trying to normalize the current URL: %s\n%v\n", rawCurrentURL, err)
	// 	return
	// }

	visited, err := db.HasURLBeenVisited(normalizedURL)
	if err != nil {
		log.Printf("Error: %v - skipping this URL...\n", err)
		return
	}

	if visited {
		log.Printf("Skipping %v - already visited", normalizedURL)
	}

	// // Check if normalized URL already has been visited in our crawled pages to ensure no repeat visits
	// isFirst := cfg.addPageVisit(normalizedURL)
	// if !isFirst {
	// 	return
	// }

	fmt.Printf("getting HTML of %s...\n", rawCurrentURL)
	currentHTML, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("error trying to get HTML of current URL: %s\n%v\n", rawCurrentURL, err)
		return
	}

	parsedURLs, _, err := getURLsFromHTML(currentHTML, rawCurrentURL)
	if err != nil {
		fmt.Printf("error trying to parse URLs from HTML of %s\n%v\n", rawCurrentURL, err)
		return
	}

	for _, URL := range parsedURLs {
		cfg.Wg.Add(1)
		fmt.Printf("crawling to next URL: %s...\n", URL)
		go cfg.CrawlPage(URL, db)
	}
}
