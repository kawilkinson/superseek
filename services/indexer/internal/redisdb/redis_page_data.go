package redisdb

import (
	"context"
	"log"

	"github.com/kawilkinson/services/indexer/internal/indexerutil"
	"github.com/kawilkinson/services/indexer/internal/models"
)

func (db *RedisClient) PopPage(ctx context.Context) string {
	if db.Client == nil {
		log.Println("no Redis client found for pop page")
		return ""
	}
	
	poppedPage, err := db.Client.BRPop(ctx, indexerutil.Timeout, indexerutil.IndexerQueueKey).Result()
	if err != nil {
		log.Printf("unable to fetch page from message queue: %v", err)
		return ""
	}

	if len(poppedPage) < 2 {
		log.Printf("unexpected pop result from message queue: %v", poppedPage)
		return ""
	}

	pageID := poppedPage[1]

	return pageID
}

func (db *RedisClient) PeekPage(ctx context.Context) string {
	if db.Client == nil {
		log.Println("no Redis client found for peek page")
		return ""
	}

	peekedPage, err := db.Client.LRange(ctx, indexerutil.IndexerQueueKey, -1, -1).Result()
	if err != nil {
		log.Printf("unable to peek page from message queue: %v", err)
		return ""
	}

	if len(peekedPage) == 0 {
		log.Printf("unable to peek page from message queue, nothing return")
		return ""
	}

	pageID := peekedPage[0]
	log.Printf("peeked page from message queue: %s", pageID)

	return pageID
}

func (db *RedisClient) GetPageData(ctx context.Context, key string) *models.Page {
	if db.Client == nil {
		log.Println("no Redis client found for get page data")
		return nil
	}

	pageHashed, err := db.Client.HGetAll(ctx, key).Result()
	if err != nil {
		log.Printf("unable to get page data from %s: %v", key, err)
		return nil
	}

	if len(pageHashed) == 0 {
		log.Printf("page with key %s not found in Redis", key)
		return nil
	}

	log.Printf("page with key %s successfully fetched", key)

	page := models.FromHash(pageHashed)

	return page
}

func (db *RedisClient) DeletePageData(ctx context.Context, key string) {
	if db.Client == nil {
		log.Println("no Redis client found for delete page data")
		return
	}

	result, err := db.Client.Del(ctx, key).Result()
	if err != nil {
		log.Printf("unable to delete key %s from Redis: %v", key, err)
		return
	}

	if result <= 0 {
		log.Printf("could not remove page data for %s from Redis", key)
		return
	}
}
