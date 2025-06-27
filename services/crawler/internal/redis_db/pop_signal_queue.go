package redis_db

import (
	"fmt"

	"github.com/kawilkinson/search-engine/internal/crawler_utilities"
)

func (db *RedisDatabase) PopSignalQueue() (string, error) {
	result, err := db.Client.BRPop(db.Context, 0, crawler_utilities.SignalQueueKey).Result()
	if err != nil {
		return "", fmt.Errorf("unable to pop from the signal queue: %v", err)
	}

	return result[1], nil
}
