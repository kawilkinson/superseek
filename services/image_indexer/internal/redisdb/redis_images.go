package redisdb

import (
	"context"
	"fmt"
	"log"

	"github.com/kawilkinson/superseek/services/image_indexer/internal/indexerutil"
	"github.com/kawilkinson/superseek/services/image_indexer/internal/models"
)

func (db *RedisClient) PopImage(ctx context.Context) string {
	if db.Client == nil {
		log.Println("no Redis client found for pop image")
		return ""
	}

	result, err := db.Client.BRPop(ctx, indexerutil.Timeout, indexerutil.ImageIndexerQueueKey).Result()
	if err != nil || len(result) != 2 {
		log.Printf("unable to get image from message queue in Redis database: %v", err)
		return ""
	}

	pageID := result[1]
	return pageID
}

func (db *RedisClient) PeekPage(ctx context.Context) *string {
	if db.Client == nil {
		log.Println("no Redis client found for peek page")
		return nil
	}

	result, err := db.Client.LRange(ctx, indexerutil.ImageIndexerQueueKey, -1, -1).Result()
	if err != nil || len(result) == 0 {
		log.Printf("unable to peek page from message queue in Redis database: %v", err)
		return nil
	}

	return &result[0]
}

func (db *RedisClient) GetPageImages(ctx context.Context, normalizedURL string) ([]string, error) {
	if db.Client == nil {
		return nil, fmt.Errorf("no Redis client found for get page images")
	}

	key := fmt.Sprintf("%s:%s", indexerutil.PageImagesPrefix, normalizedURL)

	set, err := db.Client.SMembers(ctx, key).Result()
	if err != nil || len(set) == 0 {
		return nil, fmt.Errorf("unable to get page images for key %s: %v", key, err)
	}

	return set, nil
}

func (db *RedisClient) DeletePageImages(ctx context.Context, normalizedURL string) {
	if db.Client == nil {
		log.Println("no Redis client found for delete page images")
		return
	}

	key := fmt.Sprintf("%s:%s", indexerutil.PageImagesPrefix, normalizedURL)
	deleted, err := db.Client.Del(ctx, key).Result()
	if err != nil || deleted <= 0 {
		log.Printf("unable to remove %s from Redis: %v", key, err)
		return
	}
}

func (db *RedisClient) PopImageData(ctx context.Context, imageURL string) *models.Image {
	if db.Client == nil {
		log.Println("no Redis client found for pop image data")
		return nil
	}
	key := fmt.Sprintf("%s:%s", indexerutil.ImagePrefix, imageURL)
	imageMap, err := db.Client.HGetAll(ctx, key).Result()
	if err != nil {
		log.Printf("unable to get image %s from Redis database: %v", imageURL, err)
		return nil
	}

	return models.FromHash(imageMap, imageURL)
}

func (db *RedisClient) DeleteImageData(ctx context.Context, imageURL string) {
	if db.Client == nil {
		log.Println("no Redis client found for delete image data")
		return
	}

	key := fmt.Sprintf("%s:%s", indexerutil.ImagePrefix, imageURL)
	deleted, err := db.Client.Del(ctx, key).Result()
	if err != nil || deleted <= 0 {
		log.Printf("unable to remove %s from Redis database: %v", key, err)
		return
	}
}
