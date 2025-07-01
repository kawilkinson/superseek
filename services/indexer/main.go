package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RedisDatabase struct {
	Client  *redis.Client
	Context context.Context
}

// entry into indexer service, a core service that parses the data from the crawler
func main() {
	redisHost := loadEnv("REDIS_HOST", "localhost")
	redisPort := loadEnv("REDIS_PORT", "6379")
	redisPassword := loadEnv("REDIS_PASSWORD", "")
	redisDB := loadEnv("REDIS_DB", "0")

}

func loadEnv(key, fallback string) string {
	if envVariable, exists := os.LookupEnv(key); exists {
		return envVariable
	}

	return fallback
}
