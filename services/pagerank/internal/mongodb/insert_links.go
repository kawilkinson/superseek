package mongodb

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func InsertBacklinks(ctx context.Context, backlinksColl *mongo.Collection, backlinks map[string][]string) {
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
}

func InsertOutlinks(ctx context.Context, outlinksColl *mongo.Collection, outlinksCount map[string]int) {
	cursorOutlinks, err := outlinksColl.Find(ctx, bson.D{})
	if err != nil {
		log.Fatalf("unable to get outlinks from Mongo Database: %v", err)
	}
	defer cursorOutlinks.Close(ctx)

	for cursorOutlinks.Next(ctx) {
		var doc struct {
			ID    string                 `bson:"_id"`
			Links map[string]interface{} `bson:"links"`
		}

		if err := cursorOutlinks.Decode(&doc); err != nil {
			log.Fatalf("unable to decode outlink document: %v", err)
		}

		outlinksCount[doc.ID] = len(doc.Links)
	}

	if err := cursorOutlinks.Err(); err != nil {
		log.Fatalf("cursor error: %v", err)
	}
}
