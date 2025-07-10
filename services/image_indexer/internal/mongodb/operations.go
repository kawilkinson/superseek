package mongodb

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/kawilkinson/superseek/services/image_indexer/internal/indexerutil"
	"github.com/kawilkinson/superseek/services/image_indexer/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (m *MongoClient) GetKeywords(ctx context.Context, mongoID string) (map[string]int, error) {
	log.Printf("getting keywords for %s\n", mongoID)
	collection := m.Database.Collection(indexerutil.MetadataCollection)

	filter := bson.M{"_id": mongoID}
	projection := bson.M{"keywords": 1}

	var result bson.M
	err := collection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&result)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("error finding document with %s: %v", mongoID, err)
	}

	if keywordsRaw, exists := result["keywords"]; exists {
		if keywordMap, ok := keywordsRaw.(map[string]interface{}); ok {
			keywords := make(map[string]int)
			for key, value := range keywordMap {
				if count, ok := value.(int32); ok {
					keywords[key] = int(count)
				} else if count, ok := value.(int64); ok {
					keywords[key] = int(count)
				}
			}
			return keywords, nil
		}
	}

	log.Printf("no keywords found for %s, reprocessing fields", mongoID)
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("error finding document with %s: %v", mongoID, err)
	}

	fields := []string{"summary_text", "description", "title"}
	totalWords := []string{}

	for _, field := range fields {
		if raw, exists := result[field]; exists {
			if str, ok := raw.(string); ok {
				words := strings.Fields(strings.ToLower(str))
				for _, word := range words {
					if _, isStopWord := indexerutil.StopWordsSet[word]; !isStopWord {
						totalWords = append(totalWords, word)
					}
				}
			}
		}
	}

	wordCount := make(map[string]int)
	for _, word := range totalWords {
		wordCount[word]++
	}

	return wordCount, nil
}

func (m *MongoClient) CreateWordImagesEntryOperation(word, url string, weight int) mongo.WriteModel {
	filter := bson.M{"word": word, "url": url}
	update := bson.M{"$set": bson.M{"weight": weight}}
	return mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true)
}

func (m *MongoClient) CreateWordImagesBulk(ctx context.Context, ops []mongo.WriteModel) (*mongo.BulkWriteResult, error) {
	if len(ops) == 0 {
		log.Println("no operations found to perform for word images bulk")
		return nil, nil
	}

	return m.PerformBatchOperations(ctx, ops, "word_images")
}

func (m *MongoClient) CreateImagesBulk(ctx context.Context, ops []mongo.WriteModel) (*mongo.BulkWriteResult, error) {
	if len(ops) == 0 {
		log.Println("no operations found to perform for create images bulk")
		return nil, nil
	}

	return m.PerformBatchOperations(ctx, ops, "images")
}

func (m *MongoClient) CreateImageOperation(image *models.Image) mongo.WriteModel {
	filter := bson.M{"_id": image.ID}
	update := bson.M{"$set": image.ToMap()}

	return mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true)
}

func (m *MongoClient) PerformBatchOperations(
	ctx context.Context,
	operations []mongo.WriteModel,
	collectionName string) (*mongo.BulkWriteResult, error) {
	if len(operations) == 0 {
		log.Println("no operations found to perform on mongo database")
		return nil, nil
	}

	coll := m.Database.Collection(collectionName)
	writeResults, err := coll.BulkWrite(ctx, operations, options.BulkWrite().SetOrdered(false))
	if err != nil {
		return nil, fmt.Errorf("unable to perform bulk write to mongo db: %v", err)
	}

	return writeResults, nil
}
