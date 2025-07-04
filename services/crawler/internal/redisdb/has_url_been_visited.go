package redisdb

import (
	"context"
	"fmt"
	"strconv"

	"github.com/kawilkinson/search-engine/internal/crawlutil"
)

func (db *RedisDatabase) HasURLBeenVisited(ctx context.Context, normalizedURL string) (bool, error) {
	urlToSearch := crawlutil.NormalizedURLPrefix + ":" + normalizedURL
	result, err := db.Client.HGet(ctx, urlToSearch, "visited").Result()
	if err != nil {
		return false, fmt.Errorf("unable to get %v from Redis Database: %v", urlToSearch, err)
	}

	visited, err := strconv.Atoi(result)
	if err != nil {
		return false, fmt.Errorf("unable to parse the 'visited' value: %v", err)
	}

	if visited == 0 {
		return false, nil
	}

	return true, nil
}

func (db *RedisDatabase) VisitPage(ctx context.Context, normalizedURL string) error {
	urlToSearch := crawlutil.NormalizedURLPrefix + ":" + normalizedURL

	_, err := db.Client.HSet(ctx, urlToSearch, "visited", 1).Result()
	if err != nil {
		return fmt.Errorf("unable to update visit %v from Redis: %v", urlToSearch, err)
	}

	return nil
}
