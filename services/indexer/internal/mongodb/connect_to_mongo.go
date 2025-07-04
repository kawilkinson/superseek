package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

type MongoClient struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoClient(ctx context.Context, cfg MongoConfig) (*MongoClient, error) {
	uri := "mongodb://"
	if cfg.Username != "" {
		uri += cfg.Username + ":" + cfg.Password + "@"
	}
	uri += cfg.Host + ":" + cfg.Port + "/" + cfg.Database + "?authSource=admin"

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
		Database: client.Database(cfg.Database),
	}, nil
}
