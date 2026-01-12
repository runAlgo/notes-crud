package db

import (
	"context"
	"fmt"
	"time"

	"github.com/runAlgo/notes-api/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connects to MongoDB 
// Verifies the connection 
// Returns: Mongo client, Mongo database, error if something fails
func Connect(cfg config.Config) (*mongo.Client, *mongo.Database, error) {
	// Prevents the app from freezing forever if MongoDB is unreachable
	// Gives MongoDB 10 seconds to respond
	// defer cancel() cleans up resources after function finishes
	// If MongoDB doesn't respond in 10 seconds -> operation fails.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create MongoDB client options
	// Reads MongoDB connection string from config
	// These options tell MongoDB where and how to connect
	clientOpts := options.Client().ApplyURI(cfg.MongoURI)

	//Connect to MongoDB
	// Tries to establish a connection using:
	// * the context(timeout)
	// * the connection URI
	// If connection fails -> return error
	// If successful -> continue
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, nil, fmt.Errorf("mongo connect failed")
	}

	//mongo.Connect() does NOT guarantee MongoDB is reachable
	// Ping() confirms MongoDB is actually running and responding
	// If MongoDB is down -> this fails early
	if err := client.Ping(ctx, nil); err != nil {
		return nil, nil, fmt.Errorf("mongo ping failed")
	}

	// Select the database
	// Picks the database name from config
	// MongoDB creates the database only when data is inserted
	database := client.Database(cfg.MongoDB)


	// client -> MongoDB client (used globally)
	// database -> Database instance
	// nil -> No error
	return client, database, nil
}


// Gracefully closes the MongoDB connection when your app stops.
// create context with timeout
// Gives MongoDB 5 seconds to disconnect cleanly
func Disconnect(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Disconnect client
	// Closes all open connections
	// Frees resources
	// Prevents memory leaks
	return client.Disconnect(ctx)
}
