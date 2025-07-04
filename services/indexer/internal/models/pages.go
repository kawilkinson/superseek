package models

import (
	"log"
	"strconv"
	"time"
)

type Page struct {
	NormalizedURL string
	HTML          string
	ContentType   string
	StatusCode    int
	LastCrawled   time.Time
}

func FromHash(pageData map[string]string) *Page {
	if pageData == nil {
		log.Println("unable to get page data for FromHash, no data found")
		return nil
	}

	lastCrawled, err := time.Parse(time.RFC1123, pageData["last_crawled"])
	if err != nil {
		log.Printf("unable to parse 'last_crawled' key of page data: %v", err)
		return nil
	}

	statusCode, err := strconv.Atoi(pageData["status_code"])
	if err != nil {
		log.Printf("unable to convert to int from 'status_code' key of page data: %v", err)
		return nil
	}

	return &Page{
		NormalizedURL: pageData["normalized_url"],
		HTML:          pageData["html"],
		ContentType:   pageData["content_type"],
		StatusCode:    statusCode,
		LastCrawled:   lastCrawled,
	}
}
