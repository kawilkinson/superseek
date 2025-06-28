package redisdb

import (
	"fmt"

	"github.com/kawilkinson/search-engine/internal/crawlutil"
)

func (db *RedisDatabase) GetIndexerQueueSize() (int64, error) {
	size, err := db.Client.LLen(db.Context, crawlutil.IndexerQueueKey).Result()
	if err != nil {
		return -1, fmt.Errorf("unable to get %v's size: %v", crawlutil.IndexerQueueKey, err)
	}

	return size, nil
}
