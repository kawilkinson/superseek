package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/kawilkinson/services/indexer/internal/indexerutil"
	"github.com/kawilkinson/services/indexer/internal/models"
	"github.com/kawilkinson/services/indexer/internal/mongodb"
	"github.com/kawilkinson/services/indexer/internal/redisdb"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Operations struct {
	wordOps     []mongo.WriteModel
	metadataOps []mongo.WriteModel
	outlinkOps  []mongo.WriteModel
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

		pageID := redisClient.PopPage(ctx)
		if pageID == "" {
			log.Println("no page ID found")
			time.Sleep(1 * time.Second)
			continue
		}

		log.Printf("fetching %s...\n", pageID)
		page := redisClient.GetPageData(ctx, pageID)
		if page == nil {
			log.Printf("could not fetch %s. skipping...\n", pageID)
			continue
		}

		oldMetadata, err := mongoClient.GetMetadata(ctx, page.NormalizedURL)
		if err == nil && oldMetadata != nil && oldMetadata.LastCrawled.Equal(page.LastCrawled) {
			log.Printf("no updates to %s. skipping...\n", oldMetadata.ID)
			continue
		}

		htmlData, err := indexerutil.GetHTMLData(page.HTML)
		if err != nil {
			log.Printf("skipping %s, error with parsing HTML: %v\n", pageID, err)
			continue
		} else if htmlData["language"] != "English" {
			log.Printf("skipping %s, language is not in english...\n", pageID)
			continue
		} else if len(htmlData["text"].([]string)) == 0 {
			log.Printf("skipping %s, unable to process its text...\n", pageID)
			continue
		}

		log.Printf("counting words from %s...\n", pageID)
		wordsFrequency := countWords(htmlData["text"].([]string))
		keywords := mostCommonWords(wordsFrequency, indexerutil.MaxIndexWords)

		urlWords := indexerutil.SplitURL(page.NormalizedURL)
		for _, word := range urlWords {
			if pastScore, found := keywords[word]; found {
				keywords[word] = pastScore * 50
			} else {
				keywords[word] = 100
			}
		}

		for word, frequency := range keywords {
			op := mongodb.CreateWordsEntryOperation(word, page.NormalizedURL, frequency)
			ops.wordOps = append(ops.wordOps, op)
		}

		metadata, err := models.FromMap(htmlData)
		if err != nil {
			log.Printf("unable to convert metadata in html data for %s: %v", pageID, err)
		}
		metadataOp := mongodb.CreateMetadataEntryOperation(*page, *metadata, keywords)
		ops.metadataOps = append(ops.metadataOps, metadataOp)

		outlinks := redisClient.GetOutlinks(ctx, page.NormalizedURL)
		if outlinks != nil {
			outlinkOp := mongodb.CreateOutlinksEntryOperation(*outlinks)
			ops.outlinkOps = append(ops.outlinkOps, outlinkOp)
		}

		lowerWords := make([]string, 0, len(htmlData["text"].([]string)))
		for _, word := range htmlData["text"].([]string) {
			lowerWords = append(lowerWords, strings.ToLower(word))
		}
		mongoClient.AddWordsToDictionary(ctx, lowerWords)

		redisClient.DeletePageData(ctx, pageID)
		redisClient.DeleteOutlinks(ctx, page.NormalizedURL)
		redisClient.PushToImageIndexerQueue(ctx, page.NormalizedURL)

		ops.wordOps = flushIfNeeded(ctx, mongoClient, indexerutil.WordsCollection, ops.wordOps, indexerutil.WordsOpThreshold)
		ops.metadataOps = flushIfNeeded(ctx, mongoClient, indexerutil.MetadataCollection, ops.metadataOps, indexerutil.MetadataOpThreshold)
		ops.outlinkOps = flushIfNeeded(ctx, mongoClient, indexerutil.OutlinksCollection, ops.outlinkOps, indexerutil.OutlinksOpThreshold)
	}

	log.Println("final flush...")
	flushIfAny(ctx, mongoClient, indexerutil.WordsCollection, ops.wordOps)
	flushIfAny(ctx, mongoClient, indexerutil.MetadataCollection, ops.metadataOps)
	flushIfAny(ctx, mongoClient, indexerutil.OutlinksCollection, ops.outlinkOps)

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

func countWords(text []string) map[string]int {
	wordsFrequency := make(map[string]int)

	for _, word := range text {
		_, exists := wordsFrequency[word]
		if !exists {
			wordsFrequency[word] = 1
		} else {
			wordsFrequency[word] += 1
		}
	}

	return wordsFrequency
}

func flushIfNeeded(
	ctx context.Context,
	mc *mongodb.MongoClient,
	collection string,
	operations []mongo.WriteModel,
	threshold int,
) []mongo.WriteModel {
	if len(operations) >= threshold {
		log.Printf("flushing %s bulk ops (%d entries)...\n", collection, len(operations))
		_, err := mc.PerformBatchOperations(ctx, operations, collection)
		if err != nil {
			log.Printf("error flushing %s ops: %v\n", collection, err)
		}
		return nil
	}
	return operations
}

func flushIfAny(ctx context.Context, mc *mongodb.MongoClient, collection string, operations []mongo.WriteModel) {
	if len(operations) > 0 {
		log.Printf("final flush for %s (%d entries)\n", collection, len(operations))
		_, err := mc.PerformBatchOperations(ctx, operations, collection)
		if err != nil {
			log.Printf("error in final flush for %s: %v\n", collection, err)
		}
	}
}

func mostCommonWords(wordsFrequency map[string]int, n int) map[string]int {
	type kv struct {
		Key   string
		Value int
	}

	var freqList []kv
	for key, value := range wordsFrequency {
		freqList = append(freqList, kv{key, value})
	}

	// sort by descending
	sort.Slice(freqList, func(i, j int) bool {
		// if same frequency on two words, sort alphabetically
		if freqList[i].Value == freqList[j].Value {
			return freqList[i].Key < freqList[j].Key
		}

		return freqList[i].Value > freqList[j].Value
	})

	topWords := make(map[string]int)
	for i := 0; i < len(freqList) && i < n; i++ {
		topWords[freqList[i].Key] = freqList[i].Value
	}

	return topWords
}
