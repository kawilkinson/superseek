package redisdb

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func ConnectToRedis(ctx context.Context, redisPort, redisDatabase int, redisHost, redisPassword string) (*RedisClient, error) {
	log.Println("Attempting connection to Redis database...")

	redisPortStr := strconv.Itoa(redisPort)

	client := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPortStr,
		Password: redisPassword,
		DB:       redisDatabase,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("unable to connect to Redis database of host %v: %v", redisHost, err)
	}

	log.Println("Connection to Redis database successful")

	return &RedisClient{
		Client: client,
	}, nil
}
