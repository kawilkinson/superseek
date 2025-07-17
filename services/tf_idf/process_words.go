package main

import (
	"context"
	"log"
	"math"

	"github.com/kawilkinson/superseek/services/tf_idf/internal/models"
	"github.com/kawilkinson/superseek/services/tf_idf/internal/mongodb"
	"github.com/kawilkinson/superseek/services/tf_idf/internal/tfidfutils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// main loop that applies TF-IDF logic to passed in words
func processWords(ctx context.Context, threadID int, wordChan <-chan models.WordItem, mongoClient *mongodb.MongoClient, totalDocs int) {
	for running {
		wordItem, ok := <-wordChan
		if !ok {
			break
		}

		word := wordItem.Word
		docCount, err := mongoClient.GetWordDocumentCount(ctx, word)
		if err != nil {
			log.Printf("unable to process words: %v\n", err)
			continue
		}

		entries, err := mongoClient.GetWordDocuments(ctx, word)
		if err != nil {
			log.Printf("unable to process words: %v\n", err)
			continue
		}

		idf := math.Log10(float64(totalDocs) / (1 + float64(docCount)))

		log.Printf("thread-%d: calculating TF-IDF for %q\n", threadID, word)

		localOps := make([]mongo.WriteModel, 0, len(entries))
		for _, entry := range entries {
			tfidf := entry.TF * idf

			op, err := mongoClient.UpdatePageTfidfOp(ctx, word, entry.URL, idf, tfidf)
			if err != nil {
				log.Printf("unable to process words: %v\n", err)
				continue
			}

			localOps = append(localOps, op)
		}

		operationsLock.Lock()
		bulkOperations = append(bulkOperations, localOps...)
		if len(bulkOperations) > tfidfutils.OperationsThreshold {
			log.Printf("Thread-%d: Performing bulk operations...", threadID)
			mongoClient.UpdatePageTfidfBulk(ctx, bulkOperations)
			bulkOperations = nil
		}
		operationsLock.Unlock()

		log.Printf("thread-%d: processed %d entries for word: %s\n", threadID, len(entries), word)
	}
}
