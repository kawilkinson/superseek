package mongodb

import (
	"context"
	"fmt"
	"log"

	"github.com/kawilkinson/superseek/services/backlinks_processor/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (mc *MongoClient) SaveAllBacklinks(ctx context.Context, backlinks []models.Backlinks) (*mongo.BulkWriteResult, error) {
	if mc.Client == nil {
		return nil, fmt.Errorf("unable to save all backlinks, no Mongo client found")
	}

	var ops []mongo.WriteModel
	for _, backlink := range backlinks {
		for _, link := range backlink.Links {
			saveOp := mongo.NewUpdateOneModel().
				SetFilter(bson.M{"_id": backlink.ID}).
				SetUpdate(bson.M{"$addToSet": bson.M{"links": link}}).
				SetUpsert(true)

			ops = append(ops, saveOp)
		}
	}

	collection := mc.Database.Collection("backlinks")
	opts := options.BulkWrite().SetOrdered(false)
	result, err := collection.BulkWrite(ctx, ops, opts)
	if err != nil {
		return nil, fmt.Errorf("unable to bulk write to Mongo database: %v", err)
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
