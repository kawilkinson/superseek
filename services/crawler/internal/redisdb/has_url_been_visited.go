package redisdb

import (
	"fmt"
	"strconv"

	"github.com/kawilkinson/search-engine/internal/crawlutil"
)

func (db *RedisDatabase) HasURLBeenVisited(normalizedURL string) (bool, error) {
	urlToSearch := crawlutil.NormalizedURLPrefix + ":" + normalizedURL
	result, err := db.Client.HGet(db.Context, urlToSearch, "visited").Result()
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

func (db *RedisDatabase) VisitPage(normalizedURL string) error {
	urlToSearch := crawlutil.NormalizedURLPrefix + ":" + normalizedURL

	_, err := db.Client.HSet(db.Context, urlToSearch, "visited", 1).Result()
	if err != nil {
		return fmt.Errorf("unable to update visit %v from Redis: %v", urlToSearch, err)
	}

	return nil
}
