package mongodb

import (
	"context"
	"fmt"
	"log"

	"github.com/kawilkinson/superseek/services/tf_idf/internal/tfidfutils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (mc *MongoClient) GetDocumentCount(ctx context.Context) (int, error) {
	if mc.Client == nil {
		return 0, fmt.Errorf("unable to run get document count, no Mongo client found")
	}

	coll := mc.Database.Collection(tfidfutils.MetadataCollection)
	result, err := coll.EstimatedDocumentCount(ctx)
	if err != nil {
		return 0, fmt.Errorf("unable to get estimated document count from Mongo database: %v", err)
	}

	return int(result), nil
}

func (mc *MongoClient) GetUniqueWords(ctx context.Context) (*mongo.Cursor, error) {
	if mc.Client == nil {
		return nil, fmt.Errorf("unable to run get unique words, no Mongo client found")
	}

	coll := mc.Database.Collection(tfidfutils.WordsCollection)

	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.M{"_id": "$word"}}},
		{{Key: "$project", Value: bson.M{"word": "$_id", "_id": 0}}},
	}

	opts := options.Aggregate().SetAllowDiskUse(true)

	cursor, err := coll.Aggregate(ctx, pipeline, opts)
	if err != nil {
		return nil, fmt.Errorf("unable to aggregate collection for unique words: %v", err)
	}

	return cursor, nil
}

func (mc *MongoClient) GetWordDocumentCount(ctx context.Context, word string) (int, error) {
	if mc.Client == nil {
		return 0, fmt.Errorf("unable to run get word document count, no Mongo client found")
	}

	coll := mc.Database.Collection(tfidfutils.WordsCollection)

	count, err := coll.CountDocuments(ctx, bson.M{"word": word})
	if err != nil {
		return 0, fmt.Errorf("unable to count documents for %s: %v", word, err)
	}

	return int(count), nil
}

func (mc *MongoClient) GetWordDocuments(ctx context.Context, word string) (*mongo.Cursor, error) {
	if mc.Client == nil {
		return nil, fmt.Errorf("unable to run get word documents, no Mongo client found")
	}

	coll := mc.Database.Collection(tfidfutils.WordsCollection)

	cursor, err := coll.Find(ctx, bson.M{"word": word})
	if err != nil {
		return nil, fmt.Errorf("unable to find word %s in collection: %v", word, err)
	}

	return cursor, nil
}

func (mc *MongoClient) UpdatePageTfidfOp(ctx context.Context, word string, url string, idf, tfidf float64) (mongo.WriteModel, error) {
	if mc.Client == nil {
		return nil, fmt.Errorf("")
	}

	return mongo.NewUpdateOneModel().
		SetFilter(bson.M{"word": word, "url": url}).
		SetUpdate(bson.M{
			"$set": bson.M{
				"weight": tfidf,
				"idf":    idf,
			},
		}), nil
}

func (mc *MongoClient) UpdatePageTfidfBulk(ctx context.Context, ops []mongo.WriteModel) (*mongo.BulkWriteResult, error) {
	if mc.Client == nil {
		return nil, fmt.Errorf("unable to perform update page tfidf bulk, no Mongo client found")
	}
	if len(ops) == 0 {

		return nil, fmt.Errorf("no operations found for tfidf bulk update")
	}

	result, err := mc.PerformBatchOperations(ctx, ops, tfidfutils.WordsCollection)
	if err != nil {
		return nil, fmt.Errorf("unable to update page tfidf bulk: %v", err)
	}

	return result, nil
}

func (mc *MongoClient) PerformBatchOperations(
	ctx context.Context,
	operations []mongo.WriteModel,
	collectionName string) (*mongo.BulkWriteResult, error) {
	if len(operations) == 0 {
		log.Println("no operations found to perform on mongo database")
		return nil, nil
	}

	coll := mc.Database.Collection(collectionName)
	writeResults, err := coll.BulkWrite(ctx, operations, options.BulkWrite().SetOrdered(false))
	if err != nil {
		return nil, fmt.Errorf("unable to perform bulk write to mongo db: %v", err)
	}

	return writeResults, nil
}
