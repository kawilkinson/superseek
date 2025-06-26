package redis_db

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/kawilkinson/search-engine/internal/spider"
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

func (db *RedisDatabase) addURL(rawURL, normalizedURL string) error {
	urlToSearch := spider.NormalizedURLPrefix + ":" + normalizedURL

	exists, err := db.Client.Exists(db.Context, urlToSearch).Result()
	if err != nil {
		return fmt.Errorf("URL key not found in Redis database: %v", err)
	}

	if exists > 0 {
		return fmt.Errorf("URL key already exists in Redis database")
	}

	err = db.Client.HSet(db.Context, urlToSearch, map[string]interface{}{
		"raw_url": rawURL,
		"visited": 0,
	}).Err()
	if err != nil {
		return fmt.Errorf("unable to store URL in Redis database: %v", err)
	}

	return nil
}

func (db *RedisDatabase) PushURLToQueue(rawURL string, score float64) error {
	strippedURL, err := spider.StripURL(rawURL)
	if err != nil {
		return fmt.Errorf("unable to strip URL: %v", err)
	}

	normalizedURL, err := spider.NormalizeURL(strippedURL)
	if err != nil {
		return fmt.Errorf("unable to normalize URL: %v", err)
	}

	err = db.addURL(rawURL, normalizedURL)
	if err != nil {
		return fmt.Errorf("unable to add URL: %v, it is already in the queue: %v", rawURL, err)
	}

	log.Printf("URL %v (%v) to Redis queue\n", rawURL, normalizedURL)

	return nil
}
