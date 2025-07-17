package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/kawilkinson/superseek/services/tf_idf/internal/models"
	"github.com/kawilkinson/superseek/services/tf_idf/internal/mongodb"
	"github.com/kawilkinson/superseek/services/tf_idf/internal/tfidfutils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	running        = true
	operationsLock sync.Mutex
	bulkOperations []mongo.WriteModel
)

func main() {
	mongoHost := loadEnvString("MONGO_HOST", "localhost")
	mongoPort := loadEnvInt("MONGO_PORT", 27017)
	mongoPassword := loadEnvString("MONGO_PASSWORD", "")
	mongoDatabase := loadEnvString("MONGO_DB", "test")
	mongoUsername := loadEnvString("MONGO_USERNAME", "")

	ctx := context.Background()

	mongoClient, err := mongodb.ConnectToMongo(ctx, mongoPort, mongoHost, mongoUsername, mongoPassword, mongoDatabase)
	if err != nil {
		log.Fatalf("unable to connect to Mongo database: %v", err)
	}

	totalDocs, err := mongoClient.GetDocumentCount(ctx)
	if err != nil {
		log.Fatalf("unable to get document count from Mongo database: %v", err)
	}

	uniqueWords, err := mongoClient.GetUniqueWords(ctx)
	if err != nil {
		log.Fatalf("unable to get unique words from Mongo database: %v", err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		handleExit(mongoClient)
	}()

	wordChan := make(chan models.WordItem, len(uniqueWords))
	for _, word := range uniqueWords {
		wordChan <- word
	}
	close(wordChan)

	var wg sync.WaitGroup
	for i := 0; i < tfidfutils.MaxWorkerThreads; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			processWords(ctx, id, wordChan, mongoClient, totalDocs)
		}(i + 1)
	}
	wg.Wait()

	operationsLock.Lock()
	if len(bulkOperations) > 0 {
		log.Println("performing final bulk operations...")
		mongoClient.UpdatePageTfidfBulk(ctx, bulkOperations)
	}
	operationsLock.Unlock()

	log.Println("TF-IDF processing complete. shutting down...")
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

func handleExit(mongo *mongodb.MongoClient) {
	log.Println("termination signal received - shutting down...")
	running = false

	operationsLock.Lock()
	defer operationsLock.Unlock()

	if len(bulkOperations) > 0 {
		log.Println("performing final bulk operations...")
		_, err := mongo.UpdatePageTfidfBulk(context.Background(), bulkOperations)
		if err != nil {
			log.Printf("final bulk operation error: %v\n", err)
		}
	}
	os.Exit(0)
}
