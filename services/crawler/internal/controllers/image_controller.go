package controllers

import (
	"context"
	"log"
	"time"

	"github.com/kawilkinson/search-engine/internal/crawlutil"
	"github.com/kawilkinson/search-engine/internal/pages"
	"github.com/kawilkinson/search-engine/internal/redisdb"
)

type ImageController struct {
	db *redisdb.RedisDatabase
}

func CreateImageController(db *redisdb.RedisDatabase) *ImageController {
	return &ImageController{
		db: db,
	}
}

func (imgCtrl *ImageController) SaveImages(ctx context.Context, images map[string][]*pages.Image) {
    pipeline := imgCtrl.db.Client.Pipeline()

    log.Println("saving images...")
    count := 0

    for normalizedURL, imageData := range images {
        for _, image := range imageData {
            imageKey := crawlutil.ImagePrefix + ":" + image.NormalizedSourceURL
            pipeline.HSet(ctx, imageKey, map[string]interface{}{
                "page_url": image.NormalizedPageURL,
                "alt":      image.Alt,
            })

            pipeline.Expire(ctx, imageKey, crawlutil.ImgCtrlExpirationTime*time.Hour)
            count++

            pageImagesKey := crawlutil.PageImagesPrefix + ":" + normalizedURL
            pipeline.SAdd(ctx, pageImagesKey, image.NormalizedSourceURL)
        }
    }

    _, err := pipeline.Exec(ctx)
    if err != nil {
        log.Printf("unable to save images to Redis: %v\n", err)
    } else {
        log.Printf("successfully wrote %d image entries to the Redis database\n", count)
    }
}
