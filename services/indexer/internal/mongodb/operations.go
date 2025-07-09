package mongodb

import (
	"context"
	"fmt"
	"log"

	"github.com/kawilkinson/services/indexer/internal/indexerutil"
	"github.com/kawilkinson/services/indexer/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func CreateWordsEntryOperation(word, url string, tf int) mongo.WriteModel {
	return mongo.NewUpdateOneModel().
		SetFilter(bson.M{"word": word, "url": url}).
		SetUpdate(bson.M{
			"$set": bson.M{
				"tf":     tf,
				"weight": 0,
			},
		}).
		SetUpsert(true)
}

func (m *MongoClient) CreateWordsBulk(ctx context.Context, ops []mongo.WriteModel) (*mongo.BulkWriteResult, error) {
	return m.PerformBatchOperations(ctx, ops, indexerutil.WordCollection)
}

func (m *MongoClient) GetMetadata(ctx context.Context, normalizedURL string) (*models.Metadata, error) {
	var meta models.Metadata
	coll := m.Database.Collection(indexerutil.MetadataCollection)

	err := coll.FindOne(ctx, bson.M{"_id": normalizedURL}).Decode(&meta)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("unable to find document in mongo database from %s: %v", normalizedURL, err)
	}

	return &meta, nil
}

func CreateMetadataEntryOperation(page models.Page, html models.Metadata, topWords map[string]int) mongo.WriteModel {
	doc := bson.M{
		"_id":          page.NormalizedURL,
		"title":        html.Title,
		"description":  html.Description,
		"summary_text": html.SummaryText,
		"last_crawled": page.LastCrawled,
		"keywords":     topWords,
	}

	return mongo.NewUpdateOneModel().
		SetFilter(bson.M{"_id": page.NormalizedURL}).
		SetUpdate(bson.M{"$set": doc}).
		SetUpsert(true)
}

func (m *MongoClient) CreateMetadataBulk(ctx context.Context, ops []mongo.WriteModel) (*mongo.BulkWriteResult, error) {
	return m.PerformBatchOperations(ctx, ops, indexerutil.MetadataCollection)
}

func CreateOutlinksEntryOperation(out models.Outlinks) mongo.WriteModel {
	return mongo.NewUpdateOneModel().
		SetFilter(bson.M{"_id": out.ID}).
		SetUpdate(bson.M{"$set": out}).
		SetUpsert(true)
}

func (m *MongoClient) CreateOutlinksBulk(ctx context.Context, ops []mongo.WriteModel) (*mongo.BulkWriteResult, error) {
	return m.PerformBatchOperations(ctx, ops, indexerutil.OutlinkCollection)
}

func CreateDictionaryEntryOperation(word string) mongo.WriteModel {
	return mongo.NewUpdateOneModel().
		SetFilter(bson.M{"_id": word}).
		SetUpdate(bson.M{"$set": bson.M{"_id": word}}).
		SetUpsert(true)
}

func (m *MongoClient) AddWordsToDictionary(ctx context.Context, words []string) (*mongo.BulkWriteResult, error) {
	ops := make([]mongo.WriteModel, 0, len(words))
	for _, word := range words {
		ops = append(ops, CreateDictionaryEntryOperation(word))
	}

	if len(ops) == 0 {
		return nil, nil
	}

	return m.PerformBatchOperations(ctx, ops, indexerutil.DictionaryCollection)
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
