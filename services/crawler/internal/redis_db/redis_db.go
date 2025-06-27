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

func (db *RedisDatabase) ExistsInQueue(rawURL string) (float64, bool) {
	normalizedURL, err := spider.NormalizeURL(rawURL)
	if err != nil {
		return 0.0, false
	}

	result, err := db.Client.ZScore(db.Context, spider.CrawlerQueueKey, normalizedURL).Result()
	if err != nil {
		return 0.0, false
	}

	return result, true
}

func (db *RedisDatabase) HasURLBeenVisited(normalizedURL string) (bool, error) {
	urlToSearch := spider.NormalizedURLPrefix + ":" + normalizedURL
	result, err := db.Client.HGet(db.Context, urlToSearch, "visited").Result()
	if err != nil {
		return false, fmt.Errorf("unable to get %v from Redis Database: %v", urlToSearch, err)
	}

	visited, err := strconv.Atoi(result)
	if err != nil {
		return false, fmt.Errorf("unable to parse the 'visited' value: %v", err)
	}

	if visited == 0 {
		return false, nil
	}

	return true, nil
}

func (db *RedisDatabase) PopURL() (string, float64, string, error) {
	result, err := db.Client.BZPopMin(db.Context, spider.Timeout, spider.CrawlerQueueKey).Result()
	if err != nil {
		return "", 0.0, "", fmt.Errorf("unable to pop URL from the crawler queue: %v", err)
	}

	normalizedURL := result.Z.Member.(string)
	rawURL := fmt.Sprintf("https://%v", normalizedURL)

	return rawURL, result.Z.Score, normalizedURL, nil
}

func (db *RedisDatabase) PopSignalQueue() (string, error) {
	result, err := db.Client.BRPop(db.Context, 0, spider.SignalQueueKey).Result()
	if err != nil {
		return "", fmt.Errorf("unable to pop from the signal queue: %v", err)
	}

	return result[1], nil
}

func (db *RedisDatabase) GetIndexerQueueSize() (int64, error) {
	size, err := db.Client.LLen(db.Context, spider.IndexerQueueKey).Result()
	if err != nil {
		return -1, fmt.Errorf("unable to get %v's size: %v", spider.IndexerQueueKey, err)
	}

	return size, nil
}
