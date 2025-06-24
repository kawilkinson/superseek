package webpages

import "time"

type Page struct {
	NormalizedURL string
	HTML string
	ContentType string
	StatusCode int
	LastCrawled time.Time
}

