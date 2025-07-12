package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/kawilkinson/superseek/services/pagerank/internal/pagerankutils"
)

func (mc *MongoClient) CreatePageRankEntryOperation(ctx context.Context, bulkOps []mongo.WriteModel, sortedPageRanks []pagerankutils.SortedPageRanks) error {
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
		_, err := mc.Database.Collection("pagerank").BulkWrite(ctx, bulkOps)
		if err != nil {
			return fmt.Errorf("unable to batch insert page rank values: %v", err)
		}
	}

	return nil
}
