package spider

import (
	"fmt"
	"log"

	"github.com/kawilkinson/search-engine/internal/crawlutil"
	"github.com/kawilkinson/search-engine/internal/pages"
	"github.com/kawilkinson/search-engine/internal/redisdb"
)

func (cfg *Config) CrawlPage(db *redisdb.RedisDatabase) {
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

	fmt.Printf("Getting HTML of %s...\n", rawCurrentURL)
	currentHTML, statusCode, contentType, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("error trying to get HTML of URL: %s\n%v\n", rawCurrentURL, err)
		return
	}

	parsedURLs, images, err := getURLsFromHTML(currentHTML, rawCurrentURL)
	if err != nil {
		fmt.Printf("error trying to parse URLs from HTML of %s\n%v\n", rawCurrentURL, err)
		return
	}

	cfg.AddImages(normalizedURL, images)
	cfg.UpdateLinks(normalizedURL, parsedURLs)

	page := pages.CreatePage(normalizedURL, currentHTML, contentType, statusCode)

	err = cfg.addPage(page)
	if err != nil {
		fmt.Printf("\terror adding page: %v\n", err)
		return
	}

	err = db.VisitPage(normalizedURL)
	if err != nil {
		fmt.Printf("\terror visiting page: %v\n", err)
	}

	log.Printf("Adding links from %v (%v)...\n", normalizedURL, rawCurrentURL)
	for _, URL := range parsedURLs {
		if !crawlutil.IsValidURL(URL) {
			continue
		}

		score, exists := db.ExistsInQueue(URL)
		if exists {
			score -= 0.001
		} else {
			score = depth + 1
		}

		// calculate the score here based on minimum values

		err = db.PushURLToQueue(URL, score)
		if err != nil {
			fmt.Printf("error trying to push URL to queue to update score")
			continue
		}

		cfg.Wg.Add(1)
		fmt.Printf("crawling to next URL: %s...\n", URL)
		go cfg.CrawlPage(db)
	}
}
