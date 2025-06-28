package redisdb

import (
	"fmt"

	"github.com/kawilkinson/search-engine/internal/crawlutil"
)

func (db *RedisDatabase) PopSignalQueue() (string, error) {
	result, err := db.Client.BRPop(db.Context, 0, crawlutil.SignalQueueKey).Result()
	if err != nil {
		return "", fmt.Errorf("unable to pop from the signal queue: %v", err)
	}

	return result[1], nil
}
