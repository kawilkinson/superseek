package redisdb

import (
	"context"
	"fmt"
	"log"

	"github.com/kawilkinson/search-engine/internal/crawlutil"
)

func (db *RedisDatabase) PushURLToQueue(ctx context.Context, rawURL string, score float64) error {
	strippedURL, err := crawlutil.StripURL(rawURL)
	if err != nil {
		return fmt.Errorf("unable to strip URL: %v", err)
	}

	normalizedURL, err := crawlutil.NormalizeURL(strippedURL)
	if err != nil {
		return fmt.Errorf("unable to normalize URL: %v", err)
	}

	err = db.addURL(ctx, rawURL, normalizedURL)
	if err != nil {
		return fmt.Errorf("unable to add URL: %v, it is already in the queue: %v", rawURL, err)
	}

	log.Printf("URL %v (%v) to Redis queue\n", rawURL, normalizedURL)

	return nil
}
