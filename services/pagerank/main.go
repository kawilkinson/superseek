package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/kawilkinson/superseek/services/pagerank/internal/mongodb"
	"github.com/kawilkinson/superseek/services/pagerank/internal/pagerankutils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using OS environment variables")
	}

	mongoHost := loadEnvString("MONGO_HOST", "localhost")
	mongoPort := loadEnvInt("MONGO_PORT", 27017)
	mongoPassword := loadEnvString("MONGO_PASSWORD", "")
	mongoUsername := loadEnvString("MONGO_USERNAME", "")
	mongoDatabase := loadEnvString("MONGO_DB", "test")

	ctx := context.Background()

	mongoClient, err := mongodb.ConnectToMongo(ctx, mongoPort, mongoHost, mongoUsername, mongoPassword, mongoDatabase)
	if err != nil {
		log.Printf("unable to connect to Mongo database: %v", err)
		log.Fatal("exiting...")
	}

	outlinksColl := mongoClient.Database.Collection("outlinks")
	backlinksColl := mongoClient.Database.Collection("backlinks")

	count, err := outlinksColl.CountDocuments(ctx, bson.D{})
	if err != nil {
		log.Fatalf("unable to count documents for outlinks from Mongo database: %v", err)
	}

	backlinks := make(map[string][]string)
	mongodb.InsertBacklinks(ctx, backlinksColl, backlinks)

	outlinksCount := make(map[string]int)
	mongodb.InsertOutlinks(ctx, outlinksColl, outlinksCount)

	pageRank := make(map[string]float64)
	for url := range outlinksCount {
		pageRank[url] = 1.0 / float64(count)
	}

	log.Printf("number of URLs found: %d\n", count)

	sortedPageRanks := pagerankutils.PageRankSort(pageRank, backlinks, outlinksCount, count)

	var bulkOps []mongo.WriteModel
	err = mongoClient.CreatePageRankEntryOperation(ctx, bulkOps, sortedPageRanks)
	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Println("page rank values are now saved to the Mongo database")
}

// string env loading due to needing a string value for database connections
func loadEnvString(key string, fallback string) string {
	if strVal, exists := os.LookupEnv(key); exists {
		return strVal
	}

	log.Printf("unable to load environment variable %s, using string fallback %s\n", key, fallback)
	return fallback
}

// string env loading due to needing an int value for database connections
func loadEnvInt(key string, fallback int) int {
	if envVariable, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(envVariable); err == nil {
			return intVal
		}
	}

	log.Printf("unable to load environment variable %s, using int fallback %d\n", key, fallback)
	return fallback
}
