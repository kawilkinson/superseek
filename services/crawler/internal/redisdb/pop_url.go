package redisdb

import (
	"fmt"

	"github.com/kawilkinson/search-engine/internal/crawlutil"
)

func (db *RedisDatabase) PopURL() (string, float64, string, error) {
	result, err := db.Client.BZPopMin(db.Context, crawlutil.Timeout, crawlutil.CrawlerQueueKey).Result()
	if err != nil {
		return "", 0.0, "", fmt.Errorf("unable to pop URL from the crawler queue: %v", err)
	}

	normalizedURL := result.Z.Member.(string)
	rawURL := fmt.Sprintf("https://%v", normalizedURL)

	return rawURL, result.Z.Score, normalizedURL, nil
}
