package webpages

import "time"

type Page struct {
	NormalizedURL string
	HTML string
	ContentType string
	StatusCode int
	LastCrawled time.Time
}

func CreatePage(normalizedURL, html, contentType string, statusCode int) *Page {
	return &Page {
		NormalizedURL: normalizedURL,
		HTML: html,
		ContentType: contentType,
		StatusCode: statusCode,
		LastCrawled: time.Now(),
	}
}

