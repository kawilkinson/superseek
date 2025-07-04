package redisdb

import (
	"context"
	"fmt"

	"github.com/kawilkinson/search-engine/internal/crawlutil"
)

func (db *RedisDatabase) addURL(ctx context.Context, rawURL, normalizedURL string) error {
	urlToSearch := crawlutil.NormalizedURLPrefix + ":" + normalizedURL

	exists, err := db.Client.Exists(ctx, urlToSearch).Result()
	if err != nil {
		return fmt.Errorf("URL key not found in Redis database: %v", err)
	}

	if exists > 0 {
		return fmt.Errorf("URL key already exists in Redis database")
	}

	err = db.Client.HSet(ctx, urlToSearch, map[string]interface{}{
		"raw_url": rawURL,
		"visited": 0,
	}).Err()
	if err != nil {
		return fmt.Errorf("unable to store URL in Redis database: %v", err)
	}

	return nil
}
