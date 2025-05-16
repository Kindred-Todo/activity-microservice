package main

import (
	"context"

	"github.com/Kindred-Todo/activity-microservice/config"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		fatal(ctx, "Failed to load .env", err)
	}

	config, err := config.Load()
	if err != nil {
		fatal(ctx, "Failed to load config", err)
	}

	uri := config.Atlas.URI()

	db, err := New(ctx, uri, config.Environment)
	if err != nil {
		fatal(ctx, "Failed to connect to database", err)
	}

	// open up a change stream on the collection

	defer func() {
		if err := db.Client.Disconnect(ctx); err != nil {
			fatal(ctx, "Failed to disconnect from database", err)
		}
	}()

}
