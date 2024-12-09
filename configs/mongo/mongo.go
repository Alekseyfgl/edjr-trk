package mongo

import (
	"context"
	"edjr-trk/configs/env"
	"edjr-trk/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"sync"
	"time"
)

var (
	mongoClient *mongo.Client
	once        sync.Once
)

// InitMongoSingleton initializes the MongoDB client as a singleton.
func InitMongoSingleton() {
	once.Do(func() {
		mongoURI := env.GetEnv("MONGO_URI", "")
		if mongoURI == "" {
			log.Fatal("MONGO_URI is not set in environment variables")
		}

		// Create a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Connect to MongoDB
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
		if err != nil {
			log.Fatal("Failed to connect to MongoDB", zap.Error(err))
		}

		// Ping MongoDB to verify connection
		if err := client.Ping(ctx, nil); err != nil {
			log.Fatal("Failed to ping MongoDB", zap.Error(err))
		}

		log.Info("Connected to MongoDB successfully!")
		mongoClient = client

		// Ensure unique index on email field
		ensureEmailUniqueIndex(ctx)
	})
}

// GetClient returns the initialized MongoDB client.
func GetClient() *mongo.Client {
	if mongoClient == nil {
		log.Fatal("MongoDB client is not initialized. Call InitMongoSingleton() first.")
	}
	return mongoClient
}

// CloseMongoClient disconnects the MongoDB client.
func CloseMongoClient() {
	if mongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Error("Failed to disconnect MongoDB", zap.Error(err))
		} else {
			log.Info("MongoDB client disconnected successfully.")
		}
	}
}

// ensureEmailUniqueIndex creates a unique index on the "email" field in the "users" collection.
func ensureEmailUniqueIndex(ctx context.Context) {
	usersCollection := GetClient().Database(env.GetEnv("MONGO_DB_NAME", "default_db")).Collection("users")

	// Define the index model
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "email", Value: 1}}, // Create index on "email" field
		Options: options.Index().
			SetUnique(true).               // Make the index unique
			SetName("unique_email_index"), // Optional: name for the index
	}

	// Create the index
	if _, err := usersCollection.Indexes().CreateOne(ctx, indexModel); err != nil {
		log.Fatal("Failed to create unique index on email field", zap.Error(err))
	} else {
		log.Info("Unique index on email field created successfully.")
	}
}
