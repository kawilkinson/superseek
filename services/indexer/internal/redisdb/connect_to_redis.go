package redisdb

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RedisDatabase struct {
	Client  *redis.Client
	Context context.Context
}

func (db *RedisDatabase) ConnectToRedis(redisHost, redisPort, redisPassword, redisDB string) error {
	log.Println("Attempting connection to Redis database...")

	dbIndex, err := strconv.Atoi(redisDB)
	if err != nil {
		return fmt.Errorf("unable to parse RedisDB value: %v", err)
	}

	db.Client = redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassword,
		DB:       dbIndex,
	})

	db.Context = context.Background()

	_, err = db.Client.Ping(db.Context).Result()
	if err != nil {
		return fmt.Errorf("unable to connect to Redis database of host %v: %v", redisHost, err)
	}

	log.Println("Connection to Redis database successful")
	return nil
}
