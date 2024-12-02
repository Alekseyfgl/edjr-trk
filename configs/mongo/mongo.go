package mongo

import (
	"context"
	"edjr-trk/configs/env"
	"edjr-trk/pkg/log"
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
