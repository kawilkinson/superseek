package mongodb

import (
	"context"
	"fmt"
	"log"

	"github.com/kawilkinson/superseek/services/image_indexer/internal/indexerutil"
	"github.com/kawilkinson/superseek/services/image_indexer/internal/models"
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
