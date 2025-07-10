package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/kawilkinson/superseek/services/image_indexer/internal/indexerutil"
	"github.com/kawilkinson/superseek/services/image_indexer/internal/mongodb"
	"github.com/kawilkinson/superseek/services/image_indexer/internal/redisdb"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Operations struct {
	imageOps     []mongo.WriteModel
	wordImageOps []mongo.WriteModel
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
		pageID := redisClient.PopImage(ctx)
		if pageID == "" {
			log.Println("no page ID found")
			continue
		}

		log.Printf("getting keywords of %s...\n", pageID)
		mongoID := strings.TrimPrefix(pageID, "page_images:")
		keywords, err := mongoClient.GetKeywords(ctx, mongoID)
		if err != nil {
			log.Printf("unable to get keywords for %s: %v\n", mongoID, err)
			continue
		}

		pageImages, err := redisClient.GetPageImages(ctx, pageID)
		if err != nil {
			log.Printf("unable to get page images for %s: %v\n", pageID, err)
			continue
		}

		if len(pageImages) == 0 {
			log.Printf("no images found for %s, skipping...\n", pageID)
			continue
		}

		log.Printf("got %d images from Redis database...", len(pageImages))

		var mu sync.Mutex
		var validImageOps []mongo.WriteModel
		var wordImageOps []mongo.WriteModel
		var wg sync.WaitGroup

		for _, imageURL := range pageImages {
			wg.Add(1)
			go func(imageURL string) {
				defer wg.Done()

				if !indexerutil.IsValidImage(imageURL, indexerutil.ImgMinWidth, indexerutil.ImgMinHeight) {
					log.Printf("invalid image, deleting %s\n", imageURL)
					redisClient.DeleteImageData(ctx, imageURL)
					return
				}

				if strings.HasSuffix(imageURL, ".svg") || strings.Contains(imageURL, "icons") {
					redisClient.DeleteImageData(ctx, imageURL)
					return
				}

				imageData := redisClient.PopImageData(ctx, imageURL)
				if imageData == nil {
					log.Printf("unable to get image data for %s", imageURL)
					return
				}

				imageData.Filename = path.Base(imageURL)
				words := indexerutil.SplitName(imageData.Filename)

				var localOps []mongo.WriteModel
				for _, word := range words {
					score := 30
					if val, exists := keywords[word]; exists {
						score = val * 100
					}

					op := mongoClient.CreateWordImagesEntryOperation(word, imageURL, score)
					localOps = append(localOps, op)
				}

				saveImageOp := mongoClient.CreateImageOperation(imageData)

				mu.Lock()
				validImageOps = append(validImageOps, saveImageOp)
				wordImageOps = append(wordImageOps, localOps...)
				mu.Unlock()
			}(imageURL)
		}
		wg.Wait()

		for word, weight := range keywords {
			for _, imageURL := range pageImages {
				op := mongoClient.CreateWordImagesEntryOperation(word, imageURL, weight)
				ops.wordImageOps = append(ops.wordImageOps, op)
			}
		}

		ops.imageOps = append(ops.imageOps, validImageOps...)
		ops.wordImageOps = append(ops.wordImageOps, wordImageOps...)

		if len(ops.wordImageOps) >= indexerutil.WordImagesOpThreshold {
			log.Println("flushing word image ops...")
			_, err := mongoClient.CreateWordImagesBulk(ctx, ops.wordImageOps)
			if err != nil {
				log.Printf("error trying to create word images bulk: %v\n", err)
			}
			ops.wordImageOps = nil
		}

		if len(ops.imageOps) >= indexerutil.ImageOpThreshold {
			log.Println("flushing image ops...")
			_, err := mongoClient.CreateImagesBulk(ctx, ops.imageOps)
			if err != nil {
				log.Printf("error try to create images bulk: %v\n", err)
			}
			ops.imageOps = nil
		}

		redisClient.DeletePageImages(ctx, mongoID)
	}

	log.Println("performing final bulk operations before shutdown...")
	if len(ops.wordImageOps) > 0 {
		_, err := mongoClient.CreateWordImagesBulk(ctx, ops.wordImageOps)
		if err != nil {
			log.Printf("error trying to final create word images bulk: %v", err)
		}
	}
	if len(ops.imageOps) > 0 {
		_, err := mongoClient.CreateImagesBulk(ctx, ops.imageOps)
		if err != nil {
			log.Printf("error trying to final create images bulk: %v", err)
		}
	}

	log.Println("shutting down...")
	os.Exit(0)
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
