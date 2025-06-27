package redis_db

import "github.com/kawilkinson/search-engine/internal/spider"

func (db *RedisDatabase) ExistsInQueue(rawURL string) (float64, bool) {
	normalizedURL, err := spider.NormalizeURL(rawURL)
	if err != nil {
		return 0.0, false
	}

	result, err := db.Client.ZScore(db.Context, spider.CrawlerQueueKey, normalizedURL).Result()
	if err != nil {
		return 0.0, false
	}

	return result, true
}
