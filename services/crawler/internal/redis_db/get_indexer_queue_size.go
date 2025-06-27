package redis_db

import (
	"fmt"

	"github.com/kawilkinson/search-engine/internal/crawler_utilities"
)

func (db *RedisDatabase) GetIndexerQueueSize() (int64, error) {
	size, err := db.Client.LLen(db.Context, crawler_utilities.IndexerQueueKey).Result()
	if err != nil {
		return -1, fmt.Errorf("unable to get %v's size: %v", crawler_utilities.IndexerQueueKey, err)
	}

	return size, nil
}
