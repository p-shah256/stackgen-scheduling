package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collections = []string{"users", "userAvailability", "events", "recommendations"}

func InitMongoDB(client *mongo.Client, dbName string) error {
	ctx := context.Background()
	db := client.Database(dbName)

	for _, coll := range collections {
		db.CreateCollection(ctx, coll)
	}

	userCollection := db.Collection("users")
	_, err := userCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "userId", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		log.Fatal("Failed to create unique index for user ID:", err)
	}

	userAvailCollection := db.Collection("userAvailability")
	_, err = userAvailCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "userId", Value: 1}, {Key: "date", Value: 1}},
		},
	})
	if err != nil {
		log.Fatal("Failed to create indices for user availability:", err)
	}

	eventCollection := db.Collection("events")
	_, err = eventCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		log.Fatal("Failed to create indices for events:", err)
	}

	recCollection := db.Collection("recommendations")
	_, err = recCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "eventId", Value: 1}, {Key: "date", Value: 1}},
		},
	})
	if err != nil {
		log.Fatal("Failed to create indices for recommendations:", err)
	}

	log.Println("MongoDB initialization complete with all indices set up.")
	return nil
}
