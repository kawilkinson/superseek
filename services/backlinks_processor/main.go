package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/kawilkinson/superseek/services/backlinks_processor/internal/mongodb"
	"github.com/kawilkinson/superseek/services/backlinks_processor/internal/redisdb"
)

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

	ctx := context.Background()

	log.Println("initializing Redis...")

	redisClient, err := redisdb.ConnectToRedis(ctx, redisPort, redisDatabase, redisHost, redisPassword)
	if err != nil {
		log.Fatalf("unable to connect to Redis database: %v", err)
	}

	mongoClient, err := mongodb.ConnectToMongo(ctx, mongoPort, mongoHost, mongoUsername, mongoPassword, mongoDatabase)
	if err != nil {
		log.Fatalf("unable to connect to Mongo database: %v", err)
	}

	running := true
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signals
		running = false
	}()

	for running {
		log.Println("processing backlinks...")

		backlinksKeys, err := redisClient.GetAllBacklinksKeys(ctx)
		if err != nil {
			log.Printf("unable to get all backlinks keys from Redis database, retrying: %v\n", err)
			continue
		}
		if len(backlinksKeys) == 0 {
			log.Println("no backlinks to process - sleeping...")
			for i := 0; i < 10; i++ {
				if !running {
					log.Println("service stopped")
					os.Exit(1)
				}
				time.Sleep(1 * time.Second)
			}
			continue
		}

		backlinks, err := redisClient.GetAllBacklinks(ctx, backlinksKeys)
		if err != nil {
			log.Printf("unable to get all backlinks from the Redis database, retrying: %v\n", err)
			continue
		}

		log.Println("removing all backlinks from Redis...")
		numOfDelBacklinks, err := redisClient.RemoveAllBacklinks(ctx, backlinksKeys)
		if err != nil {
			log.Printf("unable to remove all backlinks from Redis database - retrying: %v\n", err)
			continue
		}

		if numOfDelBacklinks > 0 {
			log.Printf("%d backlinks removed from Redis database\n", numOfDelBacklinks)
		}

		_, err = mongoClient.SaveAllBacklinks(ctx, backlinks)
		if err != nil {
			log.Printf("WARNING: unable to save all backlinks removed from Redis to Mongo database: %v\n", err)
		}

		for i := 0; i < 10; i++ {
			if !running {
				log.Println("service stopped")
				os.Exit(1)
			}
			time.Sleep(1 * time.Second)
		}
		continue
	}

	log.Println("shutting down...")
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
