package db

import (
	"context"
	"time"

	env "GOLANG_SERVER/configs/env"
	schema "GOLANG_SERVER/pkg/schema"

	"go.mongodb.org/mongo-driver/bson"
)

// get data from collection data in mongoDB by device address
func GetDataByDeviceAddress(deviceAddress string) ([]schema.GyroData, error) {
	collection = client.Database(env.GetEnv("MONGO_DB")).Collection(env.GetEnv("MONGO_COLLECTION")) // Get collection data
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)                        // Create a context with timeout
	defer cancel()                                                                                  // Defer cancel the context
	cursor, err := collection.Find(ctx, bson.M{"deviceaddress": deviceAddress})                     // Find data by device address
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var gyroData []schema.GyroData
	if err = cursor.All(ctx, &gyroData); err != nil {
		return nil, err
	}
	return gyroData, nil
}
