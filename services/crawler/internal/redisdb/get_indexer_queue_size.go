package redisdb

import (
	"context"
	"fmt"

	"github.com/kawilkinson/search-engine/internal/crawlutil"
)

func (db *RedisDatabase) GetIndexerQueueSize(ctx context.Context) (int64, error) {
	size, err := db.Client.LLen(ctx, crawlutil.IndexerQueueKey).Result()
	if err != nil {
		return -1, fmt.Errorf("unable to get %v's size: %v", crawlutil.IndexerQueueKey, err)
	}

	return size, nil
}
