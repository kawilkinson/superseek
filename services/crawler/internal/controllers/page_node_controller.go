package controllers

import (
	"context"
	"log"

	"github.com/kawilkinson/search-engine/internal/crawlutil"
	"github.com/kawilkinson/search-engine/internal/pages"
	"github.com/kawilkinson/search-engine/internal/redisdb"
)

type LinksController struct {
	db *redisdb.RedisDatabase
}

func CreateLinksController(db *redisdb.RedisDatabase) *LinksController {
	return &LinksController{
		db: db,
	}
}



func (linksCtrl *LinksController) SaveLinks(ctx context.Context, backlinks, outlinks map[string]*pages.PageNode) {
	pipeline := linksCtrl.db.Client.Pipeline()

	log.Println("saving backlinks...")
	count := len(backlinks)
	for key, backlinks := range backlinks {
		for _, link := range backlinks.GetLinks() {
			pipeline.SAdd(ctx, crawlutil.BacklinksPrefix+":"+key, link)
		}
	}

	log.Println("saving outlinks...")
	count += len(outlinks)
	for key, outlinks := range outlinks {
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
