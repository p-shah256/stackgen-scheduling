package db

import (
	"context"
	"log/slog"
	"time"

	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var Client *mongo.Client
var Database *mongo.Database

func Connect(ctx context.Context) error {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://admin:password@localhost:27017/events_app?authSource=admin"
	}

	clientOptions := options.Client().ApplyURI(uri).SetConnectTimeout(10 * time.Second)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}

	Client = client
	Database = client.Database("events_app")

	slog.Info("Database sucessfully connected")
	return nil
}

func Disconnect(ctx context.Context) error {
	if Client != nil {
		return Client.Disconnect(ctx)
	}
	return nil
}
