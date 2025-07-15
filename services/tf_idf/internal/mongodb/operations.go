package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

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
