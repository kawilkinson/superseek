package redis_db

import (
	"fmt"
	"log"

	"github.com/kawilkinson/search-engine/internal/spider"
)

func (db *RedisDatabase) PushURLToQueue(rawURL string, score float64) error {
	strippedURL, err := spider.StripURL(rawURL)
	if err != nil {
		return fmt.Errorf("unable to strip URL: %v", err)
	}

	normalizedURL, err := spider.NormalizeURL(strippedURL)
	if err != nil {
		return fmt.Errorf("unable to normalize URL: %v", err)
	}

	err = db.addURL(rawURL, normalizedURL)
	if err != nil {
		return fmt.Errorf("unable to add URL: %v, it is already in the queue: %v", rawURL, err)
	}

	log.Printf("URL %v (%v) to Redis queue\n", rawURL, normalizedURL)

	return nil
}
