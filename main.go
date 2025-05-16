package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Kindred-Todo/activity-microservice/config"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	fmt.Println("Starting activity microservice")
	
	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

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

	slog.Info("Connected to database", "uri", uri)

	wg := sync.WaitGroup{}
	wg.Add(1)

	// Start the change stream in a goroutine
	go func() {
		slog.Info("Starting Completed Tasks Stream")
		defer wg.Done()
		defer db.Stream.Close(ctx)
		for db.Stream.Next(ctx) {
			var data bson.M
			if err := db.Stream.Decode(&data); err != nil {
				slog.Error("Failed to decode task", "error", err)
			}

			task := CompletedTaskDocument{}
			if fullDoc, ok := data["fullDocument"].(bson.M); ok {
				bytes, err := bson.Marshal(fullDoc)
				if err != nil {
					slog.Error("Failed to marshal task", "error", err)
					continue
				}
				if err := bson.Unmarshal(bytes, &task); err != nil {
					slog.Error("Failed to unmarshal task", "error", err)
					continue
				}
			}

			slog.Info("Task", "task", task.Content, "timestamp", task.Timestamp)
			err := updateMonthlyActivity(ctx, task.UserID, task.Timestamp, db)
			if err != nil {
				slog.Error("Failed to update monthly activity", "error", err)
			}
			
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	slog.Info("Received shutdown signal, initiating graceful shutdown")
	
	// Cancel the context to stop the stream
	cancel()
	
	// Wait for the stream to close
	wg.Wait()

	// Disconnect from database
	slog.Info("Disconnecting from database")
	if err := db.Client.Disconnect(ctx); err != nil {
		slog.Error("Failed to disconnect from database", "error", err)
	}
}


