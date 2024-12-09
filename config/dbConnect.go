package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

// ConnectDB establishes a connection to the MongoDB database.
func ConnectDB() {
	// MongoDB URI (update with your URI as needed)
	mongoURI := "mongodb+srv://madil9227583:R8d6BjEpgN6t6o5w@cluster0.1jlm2.mongodb.net"

	// Set client options
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Create a new MongoDB client
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal("Error creating MongoDB client:", err)
	}

	// Set a timeout context for the connection
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Error pinging MongoDB:", err)
	}

	// Get a reference to the database
	DB = client.Database("codeCollab")
	log.Println("Connected to MongoDB successfully")
}
