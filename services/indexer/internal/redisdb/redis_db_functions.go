package redisdb

import (
	"context"
	"log"

	"github.com/kawilkinson/services/indexer/internal/indexerutil"
)

func (db *RedisClient) GetQueueSize(ctx context.Context) int64 {
	size, err := db.Client.LLen(ctx, indexerutil.IndexerQueueKey).Result()
	if err != nil {
		log.Printf("unable to get %v's size: %v", indexerutil.IndexerQueueKey, err)
		return -1
	}

	return size
}

func (db *RedisClient) SignalCrawler(ctx context.Context) {
	db.Client.LPush(ctx, indexerutil.SignalQueueKey, indexerutil.ResumeCrawl)
	log.Println("signaled crawler")
}

func (db *RedisClient) PushToImageIndexerQueue(ctx context.Context, normalizedURL string) {
	db.Client.LPush(ctx, indexerutil.ImageIndexerQueueKey, normalizedURL)
	log.Printf("pushed %s to image indexer queue", normalizedURL)
}
