package redis_db

import (
	"fmt"

	"github.com/kawilkinson/search-engine/internal/crawler_utilities"
)

func (db *RedisDatabase) addURL(rawURL, normalizedURL string) error {
	urlToSearch := crawler_utilities.NormalizedURLPrefix + ":" + normalizedURL

	exists, err := db.Client.Exists(db.Context, urlToSearch).Result()
	if err != nil {
		return fmt.Errorf("URL key not found in Redis database: %v", err)
	}

	if exists > 0 {
		return fmt.Errorf("URL key already exists in Redis database")
	}

	err = db.Client.HSet(db.Context, urlToSearch, map[string]interface{}{
		"raw_url": rawURL,
		"visited": 0,
	}).Err()
	if err != nil {
		return fmt.Errorf("unable to store URL in Redis database: %v", err)
	}

	return nil
}
