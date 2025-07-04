package redisdb

import (
	"context"
	"fmt"

	"github.com/kawilkinson/search-engine/internal/crawlutil"
)

func (db *RedisDatabase) PopURL(ctx context.Context) (string, float64, string, error) {
	result, err := db.Client.BZPopMin(ctx, crawlutil.Timeout, crawlutil.CrawlerQueueKey).Result()
	if err != nil {
		return "", 0.0, "", fmt.Errorf("unable to pop URL from the crawler queue: %v", err)
	}

	normalizedURL := result.Z.Member.(string)
	rawURL := fmt.Sprintf("https://%v", normalizedURL)

	return rawURL, result.Z.Score, normalizedURL, nil
}
