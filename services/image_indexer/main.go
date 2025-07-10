package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/kawilkinson/superseek/services/image_indexer/internal/mongodb"
	"github.com/kawilkinson/superseek/services/image_indexer/internal/redisdb"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Operations struct {
	imageOps []mongo.WriteModel
}

func main() {
	redisHost := loadEnvString("REDIS_HOST", "localhost")
	redisPort := loadEnvInt("REDIS_PORT", 6379)
	redisPassword := loadEnvString("REDIS_PASSWORD", "")
	redisDatabase := loadEnvInt("REDIS_DB", 0)

	mongoHost := loadEnvString("MONGO_HOST", "localhost")
	mongoPort := loadEnvInt("MONGO_PORT", 27017)
	mongoPassword := loadEnvString("MONGO_PASSWORD", "")
	mongoDatabase := loadEnvString("MONGO_DB", "test")
	mongoUsername := loadEnvString("MONGO_USERNAME", "")

	ops := Operations{}

	ctx := context.Background()

	log.Println("initializing Redis connection...")

	redisClient, err := redisdb.ConnectToRedis(ctx, redisPort, redisDatabase, redisHost, redisPassword)
	if err != nil {
		log.Printf("unable to connect to Redis database: %v", err)
		log.Fatal("exiting...")
	}

	log.Println("initializing Mongo connection...")

	mongoClient, err := mongodb.ConnectToMongo(ctx, mongoPort, mongoDatabase, mongoHost, mongoUsername, mongoPassword)
	if err != nil {
		log.Printf("unable to connect to Mongo database: %v", err)
		log.Fatal("exiting...")
	}

	running := true
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signals
		running = false
	}()

	// main loop for indexing
	for running {
		queueSize := redisClient.GetQueueSize(ctx)
		if queueSize == 0 {
			redisClient.SignalCrawler(ctx)
			log.Println("RESUME_CRAWL signal sent to Redis database")
		}

		log.Println("waiting for message queue...")

		// continue the loop with all the operations to the redis and mongo database
	}
}

func loadEnvString(key string, fallback string) string {
	if envVariable, exists := os.LookupEnv(key); exists {
		return envVariable
	}

	log.Println("unable to load environment variable, using string fallback")
	return fallback
}

func loadEnvInt(key string, fallback int) int {
	if envVariable, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(envVariable); err == nil {
			return intVal
		}
	}

	log.Println("unable to load environment variable, using integer fallback")
	return fallback
}
