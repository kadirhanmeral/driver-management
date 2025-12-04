package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoClient struct {
	Client     *mongo.Client
	Database   *mongo.Database
	DriversCol *mongo.Collection
}

type Config struct {
	URI        string
	Database   string
	TimeoutSec int
}

func NewMongoClient(cfg Config) (*MongoClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.TimeoutSec)*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(cfg.URI)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("MongoDB connection failed: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("MongoDB ping failed: %w", err)
	}

	db := client.Database(cfg.Database)

	return &MongoClient{
		Client:     client,
		Database:   db,
		DriversCol: db.Collection("drivers"),
	}, nil
}

func (m *MongoClient) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return m.Client.Disconnect(ctx)
}
