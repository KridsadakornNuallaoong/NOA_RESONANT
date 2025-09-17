package db

import (
	"context"
	"time"

	env "GOLANG_SERVER/configs/env"

	"go.mongodb.org/mongo-driver/bson"
)

func GetDeviceAddress() ([]string, error) {
	collection = client.Database(env.GetEnv("MONGO_DB")).Collection(env.GetEnv("MONGO_DEVICECOLLECTION"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Create a context with timeout
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var deviceAddresses []string
	for cursor.Next(ctx) {
		var result struct {
			DeviceAddress string `bson:"deviceaddress"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		deviceAddresses = append(deviceAddresses, result.DeviceAddress)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return deviceAddresses, nil
}
