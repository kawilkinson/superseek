package redisdb

import (
	"context"
	"fmt"

	"github.com/kawilkinson/superseek/services/backlinks_processor/internal/models"
	"github.com/redis/go-redis/v9"
)

func (db *RedisClient) GetAllBacklinksKeys(ctx context.Context) ([]string, error) {
	if db.Client == nil {
		return nil, fmt.Errorf("unable to perform get all backlinks keys operation, no Redis client found")
	}

	result, err := db.Client.Keys(ctx, "backlinks:*").Result()
	if err != nil {
		return nil, fmt.Errorf("unable to get backlinks keys from Redis database: %v", err)
	}

	return result, nil
}

func (db *RedisClient) GetAllBacklinks(ctx context.Context, backlinkKeys []string) ([]models.Backlinks, error) {
	if db.Client == nil {
		return nil, fmt.Errorf("unable to perform get all backlinks operation, no Redis client found")
	}

	var backlinksURLs []string
	pipeline := db.Client.Pipeline()
	cmds := make([]*redis.StringSliceCmd, 0, len(backlinkKeys))

	for _, backlinksID := range backlinkKeys {
		url := backlinksID[10:]
		backlinksURLs = append(backlinksURLs, url)
		cmd := pipeline.SMembers(ctx, backlinksID)
		cmds = append(cmds, cmd)
	}

	_, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to execute pipeline for get all backlinks in Redis database: %v", err)
	}

	var results []models.Backlinks
	for i, cmd := range cmds {
		members, err := cmd.Result()
		if err != nil && err != redis.Nil {
			return nil, fmt.Errorf("unable to get SMembers for key %s: %v", backlinkKeys[i], err)
		}

		linkSet := make(map[string]struct{}, len(members))
		for _, link := range members {
			linkSet[link] = struct{}{}
		}

		results = append(results, models.Backlinks{
			ID:    backlinksURLs[i],
			Links: linkSet,
		})
	}

	return results, nil
}

func (db *RedisClient) RemoveAllBacklinks(ctx context.Context, keys []string) (int, error) {
	if db.Client == nil {
		return 0, fmt.Errorf("unable to perform remove all backlinks operation, no Redis client found")
	}

	pipeline := db.Client.Pipeline()
	cmds := make([]*redis.IntCmd, 0, len(keys))

	for _, key := range keys {
		cmd := pipeline.Del(ctx, key)
		cmds = append(cmds, cmd)
	}

	_, err := pipeline.Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("unable to execute pipeline for backlinks deletion in Redis database: %v", err)
	}

	deleted := 0
	for _, cmd := range cmds {
		count, err := cmd.Result()
		if err != nil && err != redis.Nil {
			return 0, fmt.Errorf("unable to delete key in Redis database: %v", err)
		}

		deleted += int(count)
	}

	if deleted == 0 {
		return 0, nil
	}

	return deleted, nil
}
