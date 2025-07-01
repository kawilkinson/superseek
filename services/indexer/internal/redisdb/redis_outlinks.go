package redisdb

import (
	"log"

	"github.com/kawilkinson/services/indexer/internal/indexerutil"
	"github.com/kawilkinson/services/indexer/internal/pages"
)

func (db *RedisDatabase) GetOutlinks(normalizedURL string) *pages.Outlinks {
	key := indexerutil.OutlinksPrefix + ":" + normalizedURL
	result, err := db.Client.SMembers(db.Context, key).Result()
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

	return &pages.Outlinks{
		ID:    normalizedURL,
		Links: linksSet,
	}
}

func (db *RedisDatabase) DeleteOutlinks(normalizedURL string) {
	key := indexerutil.OutlinksPrefix + ":" + normalizedURL
	result, err := db.Client.Del(db.Context, key).Result()
	if err != nil {
		log.Printf("unable to delete outlinks from key %s: %v", key, err)
		return
	}

	if result <= 0 {
		log.Printf("could not remove outlinks for %s from Redis", key)
		return
	}
}
