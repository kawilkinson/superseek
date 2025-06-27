package redis_db

import "github.com/kawilkinson/search-engine/internal/crawler_utilities"

func (db *RedisDatabase) ExistsInQueue(rawURL string) (float64, bool) {
	normalizedURL, err := crawler_utilities.NormalizeURL(rawURL)
	if err != nil {
		return 0.0, false
	}

	result, err := db.Client.ZScore(db.Context, crawler_utilities.CrawlerQueueKey, normalizedURL).Result()
	if err != nil {
		return 0.0, false
	}

	return result, true
}
