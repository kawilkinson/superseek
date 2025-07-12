package main

import (
	"context"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/kawilkinson/superseek/services/pagerank/internal/mongodb"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SortedPageRanks struct {
	URL  string
	Rank float64
}

// entry into PageRank service using the PageRank algorithm
func main() {
	mongoHost := loadEnvString("MONGO_HOST", "localhost")
	mongoPort := loadEnvInt("MONGO_PORT", 27017)
	mongoPassword := loadEnvString("MONGO_PASSWORD", "")
	mongoUsername := loadEnvString("MONGO_USERNAME", "")
	mongoDatabase := loadEnvString("MONGO_DB", "test")

	ctx := context.Background()

	mongoClient, err := mongodb.ConnectToMongo(ctx, mongoPort, mongoDatabase, mongoHost, mongoUsername, mongoPassword)
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
	
	cursorBacklinks, err := backlinksColl.Find(ctx, bson.D{})
	if err != nil {
		log.Fatalf("unable to get backlinks from Mongo database: %v", err)
	}
	defer cursorBacklinks.Close(ctx)

	for cursorBacklinks.Next(ctx) {
		var doc struct {
			ID    string   `bson:"_id"`
			Links []string `bson:"links"`
		}

		if err := cursorBacklinks.Decode(&doc); err != nil {
			log.Fatalf("unable to decode backlink document: %v", err)
		}

		backlinks[doc.ID] = doc.Links
	}

	outlinksCount := make(map[string]int)
	cursorOutlinks, err := outlinksColl.Find(ctx, bson.D{})
	if err != nil {
		log.Fatalf("unable to get outlinks from Mongo Database")
	}
	defer cursorOutlinks.Close(ctx)

	for cursorOutlinks.Next(ctx) {
		var doc struct {
			ID    string   `bson:"_id"`
			Links []string `bson:"links"`
		}

		if err := cursorOutlinks.Decode(&doc); err != nil {
			log.Fatalf("unable to decode outlink document: %v", err)
		}

		outlinksCount[doc.ID] = len(doc.Links)
	}

	pageRank := make(map[string]float64)
	for url := range outlinksCount {
		pageRank[url] = 1.0 / float64(count)
	}

	log.Printf("number of URLs found: %d\n", count)

	iterations := 10
	damping := 0.85
	for i := 0; i < iterations; i++ {
		newPageRank := make(map[string]float64)

		for url := range pageRank {
			var newCumulativeRank float64

			if backlinksForURL, exists := backlinks[url]; exists {
				for _, backlink := range backlinksForURL {
					outlinkCount, outlinkExists := outlinksCount[backlink]
					backlinkRank, backlinkExists := pageRank[backlink]

					if backlinkExists && outlinkExists {
						newCumulativeRank += (backlinkRank / float64(outlinkCount))
					}
				}
			}
			newPageRank[url] = (1-damping)/float64(count) + (damping * newCumulativeRank)
		}
		pageRank = newPageRank
	}

	sortedPageRanks := []SortedPageRanks{}
	for url, rank := range pageRank {
		sortedPageRanks = append(sortedPageRanks, SortedPageRanks{
			URL:  url,
			Rank: rank,
		})
	}

	sort.Slice(sortedPageRanks, func(i, j int) bool {
		return sortedPageRanks[i].Rank > sortedPageRanks[j].Rank
	})

	log.Println("sorted page rank values:")

	for _, pageRank := range sortedPageRanks {
		log.Printf("page URL: %s, page rank: %f\n", pageRank.URL, pageRank.Rank)
	}

	var bulkOps []mongo.WriteModel
	for _, pageRank := range sortedPageRanks {
		bulkOps = append(bulkOps, mongo.NewUpdateOneModel().
			SetFilter(bson.D{
				{Key: "_id", Value: pageRank.URL}}).
			SetUpdate(bson.D{
				{Key: "$set", Value: bson.D{
					{Key: "rank", Value: pageRank.Rank}}}}).
			SetUpsert(true))
	}

	if len(bulkOps) > 0 {
		_, err := mongoClient.Database.Collection("pagerank").BulkWrite(ctx, bulkOps)
		if err != nil {
			log.Fatalf("unable to batch insert page rank values: %v", err)
		}
	}

	log.Println("page rank values are now saved to the Mongo database")
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
