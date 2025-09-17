package db

import (
	"context"
	"time"

	env "GOLANG_SERVER/configs/env"
	schema "GOLANG_SERVER/pkg/schema"
)

// * store data to mongo db and use upper camel case for function name
func StoreGyroData(data schema.GyroData) (bool, error) {
	collection = client.Database(env.GetEnv("MONGO_DB")).Collection(env.GetEnv("MONGO_COLLECTION")) // Get collection data
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)                        // Create a context with timeout
	defer cancel()                                                                                  // Defer cancel the context

	// load Bangkok timezone
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return false, err
	}

	currentTime := time.Now().In(loc) // Get current time

	data.DateTime = currentTime.Format(time.RFC3339) // Set timestamp to current time
	data.TimeStamp = currentTime.UnixMilli()         // Set timestamp to current time
	_, err = collection.InsertOne(ctx, data)
	if err != nil {
		return false, err
	}
	return true, nil
}
