package mongodb

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (mc *MongoClient) FlushIfNeeded(ctx context.Context, coll string, ops []mongo.WriteModel, threshold int) []mongo.WriteModel {
	if len(ops) >= threshold {
		log.Printf("flushing %s bulk ops (%d entries)...\n", coll, len(ops))
		_, err := mc.PerformBatchOperations(ctx, ops, coll)
		if err != nil {
			log.Printf("error flushing %s ops: %v\n", coll, err)
		}
		return nil
	}
	return ops
}

func (mc *MongoClient) FlushIfAny(ctx context.Context, coll string, ops []mongo.WriteModel) {
	if len(ops) > 0 {
		log.Printf("final flush for %s (%d entries)\n", coll, len(ops))
		_, err := mc.PerformBatchOperations(ctx, ops, coll)
		if err != nil {
			log.Printf("error in final flush for %s: %v\n", coll, err)
		}
	}
}
