package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Update function for month+array design
func updateMonthlyActivity(ctx context.Context, userId primitive.ObjectID, completedAt time.Time, db *DB) error {
	year := completedAt.Year()
	month := int(completedAt.Month())
	day := completedAt.Day()

	collection := db.Collections["activity"]

	// array filter is that the day has to match
	arrayFilters := options.ArrayFilters{
			Filters: []interface{}{
					bson.M{"elem.day": day},
			},
	}

	// array filters
	opts := options.Update().SetUpsert(true)
	optsFilter := options.Update().SetUpsert(true).SetArrayFilters(arrayFilters)
				
  // upsert a doucment for the current month 
	result, err := collection.UpdateOne(
		ctx, bson.M{"user": userId, "year": year, "month": month},
		 bson.M{
			"$inc": bson.M{"totalCount": 1},
			"$set": bson.M{"lastUpdated": time.Now().UTC()},
			"$setOnInsert": bson.M{"days": []bson.M{{"day": day, "count": 1, "level": 1}}}},
		 opts)
	if err != nil {
		return err
	}

	// now try to edit the actual day
	result, err = collection.UpdateOne(
		ctx, bson.M{"user": userId, "year": year, "month": month},
		bson.M{
			"$inc": bson.M{"days.$[elem].count": 1},
			"$setOnInsert": bson.M{"days.$[elem].level": 1},
		},
		optsFilter,
	)
	if err != nil {
		return err
	}

	// If no document was modified, it means the day doesn't exist, so add it
	if result.ModifiedCount == 0 {
		_, err = collection.UpdateOne(
			ctx, bson.M{"user": userId, "year": year, "month": month},
			bson.M{"$push": bson.M{"days": bson.M{"day": day, "count": 1, "level": 1}}},
			opts,
		)
		if err != nil {
			return err
		}
	}

	fmt.Println(result)
		

	return nil 	
}