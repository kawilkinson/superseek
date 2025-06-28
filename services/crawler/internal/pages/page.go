package pages

import (
	"fmt"
	"time"

	"github.com/kawilkinson/search-engine/internal/crawlutil"
)

type Page struct {
	NormalizedURL string
	HTML          string
	ContentType   string
	StatusCode    int
	LastCrawled   time.Time
}

func CreatePage(normalizedURL, html, contentType string, statusCode int) *Page {
	return &Page{
		NormalizedURL: normalizedURL,
		HTML:          html,
		ContentType:   contentType,
		StatusCode:    statusCode,
		LastCrawled:   time.Now(),
	}
}

func HashPage(page *Page) (map[string]interface{}, error) {
	return map[string]interface{}{
		"normalized_url": page.NormalizedURL,
		"html":           page.HTML,
		"content_type":   page.ContentType,
		"status_code":    page.StatusCode,
		"last_crawled":   page.LastCrawled.Format(time.RFC1123),
	}, nil
}

func DehashPage(data map[string]string) (*Page, error) {
	lastCrawled, err := crawlutil.ParseTime(data["last_crawled"])
	if err != nil {
		return nil, fmt.Errorf("unable to parse 'last_crawled' in hash: %w", err)
	}

	statusCode, err := crawlutil.ParseInt(data["status_code"])
	if err != nil {
		return nil, fmt.Errorf("unable to parse 'status_code' in hash: %w", err)
	}

	return &Page{
		NormalizedURL: data["normalized_url"],
		HTML:          data["html"],
		ContentType:   data["content_type"],
		StatusCode:    statusCode,
		LastCrawled:   lastCrawled,
	}, nil
}
