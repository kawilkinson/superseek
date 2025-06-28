package redisdb

import "github.com/kawilkinson/search-engine/internal/crawlutil"

func (db *RedisDatabase) ExistsInQueue(rawURL string) (float64, bool) {
	normalizedURL, err := crawlutil.NormalizeURL(rawURL)
	if err != nil {
		return 0.0, false
	}

	result, err := db.Client.ZScore(db.Context, crawlutil.CrawlerQueueKey, normalizedURL).Result()
	if err != nil {
		return 0.0, false
	}

	return result, true
}
