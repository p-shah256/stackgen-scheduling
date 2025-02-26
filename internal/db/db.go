package db

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	db     *mongo.Database
)

func Connect(ctx context.Context) error {
	uri := os.Getenv("MONGO_URI")
	clientOptions := options.Client().ApplyURI(uri)
	c, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	client = c
	db = client.Database(os.Getenv("MONGO_DB"))
	return nil
}

func GetCollection(name string) *mongo.Collection {
	return db.Collection(name)
}

func Close(ctx context.Context) error {
	return client.Disconnect(ctx)
}
