package redis_db

import (
	"fmt"

	"github.com/kawilkinson/search-engine/internal/spider"
)

func (db *RedisDatabase) GetIndexerQueueSize() (int64, error) {
	size, err := db.Client.LLen(db.Context, spider.IndexerQueueKey).Result()
	if err != nil {
		return -1, fmt.Errorf("unable to get %v's size: %v", spider.IndexerQueueKey, err)
	}

	return size, nil
}
