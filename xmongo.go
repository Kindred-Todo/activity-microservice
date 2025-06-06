package main

import (
	"context"
	"fmt"
	"slices"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	Client      *mongo.Client
	DB          *mongo.Database
	Collections map[string]*mongo.Collection
	Stream      *mongo.ChangeStream
}

func New(ctx context.Context, uri string, environment string) (*DB, error) {
	client, err := connectClient(ctx, uri)
	if err != nil {
		return nil, fmt.Errorf("failed to setup client: %w", err)
	}
	if err := validateEnvironment(ctx, client, environment); err != nil {
		return nil, err
	}
	db := client.Database(environment)
	collections, err := setupCollections(ctx, db)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	tasksStream, err := db.Collection("completed-tasks").Watch(ctx, 
		[]bson.D{
			{
				{"$match", bson.D{
					{"operationType", bson.D{
						{"$in", []string{"insert", "update", "replace"}},
					}},
					{"fullDocument.timestamp", bson.D{
						{"$gte", now.Add(-5 * time.Hour)},
					}},
				}},
			},
		},
		options.ChangeStream().
			SetFullDocument(options.UpdateLookup).
			SetStartAtOperationTime(&primitive.Timestamp{
				T: uint32(time.Now().Unix()),
				I: 0,
			}),
	)

	if err != nil {
		return nil, err
	}
	return &DB{
		Client:      client,
		DB:          db,
		Collections: collections,
		Stream:      tasksStream,
	}, nil
}

func connectClient(ctx context.Context, uri string) (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return client, nil
}

func validateEnvironment(ctx context.Context, client *mongo.Client, environment string) error {
	envList, err := client.ListDatabaseNames(ctx, bson.D{})
	if err != nil {
		return fmt.Errorf("failed to list databases: %w", err)
	}
	if !slices.Contains(envList, environment) {
		return fmt.Errorf("invalid database environment passed. choose from the following: %v", envList)
	}
	return nil
}

func setupCollections(ctx context.Context, db *mongo.Database) (map[string]*mongo.Collection, error) {
	collectionNames, err := db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to list collections: %w", err)
	}
	collections := make(map[string]*mongo.Collection)
	for _, name := range collectionNames {
		collections[name] = db.Collection(name)
	}
	return collections, nil
}
