package controllers

import (
	"context"
	"log"

	"github.com/kawilkinson/search-engine/internal/crawlutil"
	"github.com/kawilkinson/search-engine/internal/redisdb"
	"github.com/kawilkinson/search-engine/internal/spider"
)

type LinksController struct {
	db *redisdb.RedisDatabase
}

func CreateLinksController(db *redisdb.RedisDatabase) *LinksController {
	return &LinksController{
		db: db,
	}
}

func (linksCtrl *LinksController) SaveLinks(ctx context.Context, cfg *spider.Config) {
	pipeline := linksCtrl.db.Client.Pipeline()

	log.Println("saving backlinks...")
	data := cfg.Backlinks
	count := len(data)
	for key, backlinks := range data {
		for _, link := range backlinks.GetLinks() {
			pipeline.SAdd(ctx, crawlutil.BacklinksPrefix+":"+key, link)
		}
	}

	log.Println("saving outlinks...")
	data = cfg.Outlinks
	count += len(data)
	for key, outlinks := range data {
		for _, link := range outlinks.GetLinks() {
			pipeline.SAdd(ctx, crawlutil.OutlinksPrefix+":"+key, link)
		}
	}

	_, err := pipeline.Exec(ctx)
	if err != nil {
		log.Printf("unable to save links to Redis: %v\n", err)
	} else {
		log.Printf("successfully wrote %d link entries to Redis database\n", count)
	}
}
