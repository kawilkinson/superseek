package spider

import (
	"fmt"
	"net/url"
)

func (cfg *Config) CrawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.Wg.Done()
	}()

	if cfg.isMaxPagesReached() {
		cfg.runOnce.Do(func() {
			fmt.Printf("max pages reached, exiting WebGopher early...\n")
		})
		return
	}

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("error trying to parse current URL: %s\n%v\n", rawCurrentURL, err)
		return
	}

	if currentURL.Hostname() != cfg.baseURL.Hostname() {
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

	parsedURLs, err := getURLsFromHTML(currentHTML, cfg.baseURL.String())
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
