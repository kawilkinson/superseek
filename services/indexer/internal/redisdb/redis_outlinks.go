package redisdb

import (
	"context"
	"log"

	"github.com/kawilkinson/services/indexer/internal/indexerutil"
	"github.com/kawilkinson/services/indexer/internal/models"
)

func (db *RedisClient) GetOutlinks(ctx context.Context, normalizedURL string) *models.Outlinks {
	if db.Client == nil {
		log.Println("no Redis client found for get outlinks")
		return nil
	}
	key := indexerutil.OutlinksPrefix + ":" + normalizedURL
	result, err := db.Client.SMembers(ctx, key).Result()
	if err != nil {
		log.Printf("unable to get outlinks from key %s: %v", key, err)
		return nil
	}

	if len(result) == 0 {
		log.Printf("unable to get outlinks, no outlinks found for %s", key)
		return nil
	}

	linksSet := make(map[string]struct{}, len(result))
	for _, link := range result {
		linksSet[link] = struct{}{}
	}

	return &models.Outlinks{
		ID:    normalizedURL,
		Links: linksSet,
	}
}

func (db *RedisClient) DeleteOutlinks(ctx context.Context, normalizedURL string) {
	if db.Client == nil {
		log.Println("no Redis client found for delete outlinks")
		return
	}
	key := indexerutil.OutlinksPrefix + ":" + normalizedURL
	result, err := db.Client.Del(ctx, key).Result()
	if err != nil {
		log.Printf("unable to delete outlinks from key %s: %v", key, err)
		return
	}

	if result <= 0 {
		log.Printf("could not remove outlinks for %s from Redis", key)
		return
	}
}
