package controllers

import (
	"log"

	"github.com/kawilkinson/search-engine/internal/crawlutil"
	"github.com/kawilkinson/search-engine/internal/pages"
	"github.com/kawilkinson/search-engine/internal/redisdb"
	"github.com/kawilkinson/search-engine/internal/spider"
	"github.com/redis/go-redis/v9"
)

type PageController struct {
	db *redisdb.RedisDatabase
}

func CreatePageController(db *redisdb.RedisDatabase) *PageController {
	return &PageController{
		db: db,
	}
}

func (pgCtrl *PageController) GetAllPages() map[string]*pages.Page {
	log.Println("grabbing pages data from Redis...")
	redisPages := make(map[string]*pages.Page)

	keys, err := pgCtrl.db.Client.Keys(pgCtrl.db.Context, crawlutil.PagePrefix+":*").Result()
	if err != nil {
		log.Printf("unable to fetch page data from Redis: %v\n", err)
		return nil
	}

	pipeline := pgCtrl.db.Client.Pipeline()
	commands := make([]*redis.MapStringStringCmd, len(keys))

	for i, key := range keys {
		commands[i] = pipeline.HGetAll(pgCtrl.db.Context, key)
	}

	_, err = pipeline.Exec(pgCtrl.db.Context)
	if err != nil {
		log.Printf("unable to fetch page data from Redis pipeline: %v\n", err)
		return nil
	}

	for _, command := range commands {
		data, err := command.Result()
		if err != nil {
			log.Printf("unable to fetch pipeline result of page data from Redis: %v\n", err)
			return nil
		}

		page, err := pages.DehashPage(data)
		if err != nil {
			log.Printf("unable to dehash page data from Redis: %v\n", err)
			return nil
		}

		redisPages[page.NormalizedURL] = page
	}

	return redisPages
}

func (pgCtrl *PageController) SavePages(cfg *spider.Config) {
	data := cfg.Pages
	log.Printf("writing %d entries to the Redis database...\n", len(data))

	pipeline := pgCtrl.db.Client.Pipeline()

	for _, page := range data {
		pageHash, err := pages.HashPage(page)
		if err != nil {
			log.Printf("unable to hash page for Redis database %s: %v\n", page.NormalizedURL, err)
			continue
		}

		pageKey := crawlutil.PagePrefix + ":" + page.NormalizedURL
		pipeline.HSet(pgCtrl.db.Context, pageKey, pageHash)

		pgCtrl.db.Client.LPush(pgCtrl.db.Context, crawlutil.IndexerQueueKey, pageKey)
	}

	_, err := pipeline.Exec(pgCtrl.db.Context)
	if err != nil {
		log.Printf("unable to execute page pipeline: %v\n", err)
	} else {
		log.Printf("successfully wrote %d page entires to the Redis database\n", len(data))
	}
}
