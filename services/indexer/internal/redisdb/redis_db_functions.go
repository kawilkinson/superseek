package redisdb

import (
	"log"

	"github.com/kawilkinson/services/indexer/internal/indexerutil"
)

func (db *RedisDatabase) GetQueueSize() int64 {
	size, err := db.Client.LLen(db.Context, indexerutil.IndexerQueueKey).Result()
	if err != nil {
		log.Printf("unable to get %v's size: %v", indexerutil.IndexerQueueKey, err)
		return -1
	}

	return size
}

func (db *RedisDatabase) SignalCrawler() {
	db.Client.LPush(db.Context, indexerutil.SignalQueueKey, indexerutil.ResumeCrawl)
	log.Println("signaled crawler")
}

func (db *RedisDatabase) PushToImageIndexerQueue(normalizedURL string) {
	db.Client.LPush(db.Context, indexerutil.IndexerQueueKey, normalizedURL)
	log.Printf("pushed %s to image indexer queue", normalizedURL)
}
