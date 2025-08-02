package redisdb

import (
	"context"
	"fmt"
	"log"

	"github.com/kawilkinson/search-engine/internal/crawlutil"
	"github.com/redis/go-redis/v9"
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

	// this adds a normalized_url:rawURL key value pair permanently to the Redis database so duplicates aren't crawled
	err = db.addURL(ctx, rawURL, normalizedURL)
	if err != nil {
		return fmt.Errorf("unable to add URL: %v, it is already in the queue: %v", rawURL, err)
	}

	err = db.Client.ZAdd(ctx, crawlutil.CrawlerQueueKey, redis.Z{
		Score:  score,
		Member: normalizedURL,
	}).Err()
	if err != nil {
		return fmt.Errorf("unable to add URL to crawler queue: %w", err)
	}

	log.Printf("Pushed URL %v (%v) to Redis queue\n", rawURL, normalizedURL)

	return nil
}
