package redis_db

import (
	"fmt"
	"strconv"

	"github.com/kawilkinson/search-engine/internal/spider"
)

func (db *RedisDatabase) HasURLBeenVisited(normalizedURL string) (bool, error) {
	urlToSearch := spider.NormalizedURLPrefix + ":" + normalizedURL
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
