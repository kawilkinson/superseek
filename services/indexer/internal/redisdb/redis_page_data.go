package redisdb

import (
	"log"

	"github.com/kawilkinson/services/indexer/internal/indexerutil"
	"github.com/kawilkinson/services/indexer/internal/pages"
)

func (db *RedisDatabase) PopPage() string {
	poppedPage, err := db.Client.BRPop(db.Context, indexerutil.Timeout, indexerutil.IndexerQueueKey).Result()
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

func (db *RedisDatabase) PeekPage() string {
	peekedPage, err := db.Client.LRange(db.Context, indexerutil.IndexerQueueKey, -1, -1).Result()
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

func (db *RedisDatabase) GetPageData(key string) *pages.Page {
	pageHashed, err := db.Client.HGetAll(db.Context, key).Result()
	if err != nil {
		log.Printf("unable to get page data from %s: %v", key, err)
		return nil
	}

	if len(pageHashed) == 0 {
		log.Printf("page with key %s not found in Redis", key)
		return nil
	}

	log.Printf("page with key %s successfully fetched", key)

	page := pages.FromHash(pageHashed)

	return page
}

func (db *RedisDatabase) DeletePageData(key string) {
	result, err := db.Client.Del(db.Context, key).Result()
	if err != nil {
		log.Printf("unable to delete key %s from Redis: %v", key, err)
		return
	}

	if result <= 0 {
		log.Printf("could not remove page data for %s from Redis", key)
		return
	}
}
