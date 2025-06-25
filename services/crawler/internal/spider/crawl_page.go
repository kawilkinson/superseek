package spider

import (
	"fmt"
)

func (cfg *Config) CrawlPage(rawCurrentURL string) {
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

	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("error trying to normalize the current URL: %s\n%v\n", rawCurrentURL, err)
		return
	}

	// Check if normalized URL already has been visited in our crawled pages to ensure no repeat visits
	isFirst := cfg.addPageVisit(normalizedURL)
	if !isFirst {
		return
	}

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
		go cfg.CrawlPage(URL)
	}
}
