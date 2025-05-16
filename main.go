package main

import (
	"context"
	"log"
	"os"

	"github.com/Kindred-Todo/activity-microservice/config"
	"github.com/joho/godotenv"
)

func main() {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MONGODB_URI is not set")
	}
	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		fatal(ctx, "Failed to load .env", err)
	}

	config, err := config.Load()
	if err != nil {
		fatal(ctx, "Failed to load config", err)
	}

	db, err := New(context.TODO(), uri, config.Environment)
	if err != nil {
		panic(err)
	}

	// open up a change stream on the collection

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

}
