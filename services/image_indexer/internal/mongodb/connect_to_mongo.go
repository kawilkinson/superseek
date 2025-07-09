package mongodb

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoClient struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func ConnectToMongo(ctx context.Context, mongoPort int, mongoHost, mongoUsername, mongoPassword, mongoDatabase string) (*MongoClient, error) {
	uri := "mongodb://"
	if mongoUsername != "" {
		uri += mongoUsername + ":" + mongoPassword + "@"
	}

	mongoPortStr := strconv.Itoa(mongoPort)

	uri += mongoHost + ":" + mongoPortStr + "/" + mongoDatabase + "?authSource=admin"

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongo database: %w", err)
	}

	ctxPing, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := client.Ping(ctxPing, nil); err != nil {
		return nil, fmt.Errorf("failed to ping mongo database: %w", err)
	}

	log.Println("successfully connected to mongo database")

	return &MongoClient{
		Client:   client,
		Database: client.Database(mongoDatabase),
	}, nil
}